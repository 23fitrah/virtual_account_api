package routes

import (
	"context"
	"net/http"
	"os"
	"time"
	"virtual_account_api/constants"
	"virtual_account_api/internal/injector"
	"virtual_account_api/internal/middleware"
	"virtual_account_api/resources"
	"virtual_account_api/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(ct *injector.AppContainer) *gin.Engine {
	// Gunakan gin.New() agar tidak ada middleware bawaan (Logger & Recovery)
	r := gin.New()

	// Pasang Recovery middleware agar server tidak crash jika ada panic
	r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.NoRoute(func(c *gin.Context) {
		utils.Responds(c, resources.BaseResponse{
			ResponseCode: constants.CodeEndpointNotFound,
			Message:      constants.StatusEndpointNotFound,
		}, http.StatusNotFound)
	})

	r.Use(middleware.ESMiddleware())

	r.GET("/", func(c *gin.Context) {
		utils.Responds(c, resources.BaseResponse{
			ResponseCode: constants.CodeTransactionSuccess,
			Message:      constants.StatusGetSuccess,
		}, http.StatusOK)
		return
	})

	api := r.Group("/api")

	api.GET("/health", func(c *gin.Context) {
		healthStatus := gin.H{
			"status":   "ok",
			"services": gin.H{},
		}

		sqlDB, err := ct.DB.DB()
		if err != nil {
			healthStatus["status"] = "error"
			healthStatus["services"].(gin.H)["database"] = "connection error: " + err.Error()
		} else {
			if err := sqlDB.Ping(); err != nil {
				healthStatus["status"] = "error"
				healthStatus["services"].(gin.H)["database"] = "unreachable: " + err.Error()
			} else {
				healthStatus["services"].(gin.H)["database"] = "connected"
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if _, err := ct.RedisClient.Ping(ctx).Result(); err != nil {
			healthStatus["status"] = "error"
			healthStatus["services"].(gin.H)["redis"] = "unreachable: " + err.Error()
		} else {
			healthStatus["services"].(gin.H)["redis"] = "connected"
		}

		logDir := "logs"
		if _, err := os.Stat(logDir); os.IsNotExist(err) {
			healthStatus["status"] = "error"
			healthStatus["services"].(gin.H)["logging"] = "logs directory missing"
		} else {
			testFile := logDir + "/health_check.tmp"
			if f, err := os.Create(testFile); err != nil {
				healthStatus["status"] = "error"
				healthStatus["services"].(gin.H)["logging"] = "not writable: " + err.Error()
			} else {
				f.Close()
				os.Remove(testFile)
				healthStatus["services"].(gin.H)["logging"] = "writable (ready for filebeat)"
			}
		}

		httpStatus := http.StatusOK
		if healthStatus["status"] != "ok" {
			httpStatus = http.StatusServiceUnavailable
		}

		c.JSON(httpStatus, healthStatus)
	})

	// V1 Routes
	RegisterV1Router(api, ct)

	//api.Use(middleware.MustBeAuthenticated(ct))
	//{
	//	RegisterWebRouter(api, ct)
	//}

	return r
}
