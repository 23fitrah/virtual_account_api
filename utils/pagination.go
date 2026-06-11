package utils

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type PaginationParams struct {
	Page    int               `json:"page"`
	PerPage int               `json:"perPage"`
	SortBy  string            `json:"sortBy"`
	SortDir string            `json:"sortDir"`
	Search  string            `json:"search"`
	Filter  map[string]string `json:"filter"`
}

type PaginationResult struct {
	CurrentPage int         `json:"currentPage"`
	Data        interface{} `json:"data"`
	From        int         `json:"from"`
	LastPage    int         `json:"lastPage"`
	PerPage     int         `json:"perPage"`
	To          int         `json:"to"`
	Total       int64       `json:"total"`
}

func GetPaginationParams(ctx *gin.Context) PaginationParams {
	page := cast.ToInt(ctx.DefaultQuery("page", "1"))
	perPage := cast.ToInt(ctx.DefaultQuery("perPage", "10"))

	sortBy := ctx.Query("sortBy")
	if sortBy == "" {
		sortBy = "ROWID"
	}

	sortDir := ctx.Query("sortDir")
	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "desc"
	}

	search := ctx.DefaultQuery("search", "")
	filterMap := make(map[string]string)

	// Reserved keywords for pagination
	reservedKeys := map[string]bool{
		"page":     true,
		"per_page": true,
		"sort_by":  true,
		"sort_dir": true,
		"search":   true,
	}

	for key, values := range ctx.Request.URL.Query() {
		if len(values) == 0 {
			continue
		}

		if strings.HasPrefix(key, "filters[") && strings.HasSuffix(key, "]") {
			filterKey := key[len("filters[") : len(key)-1]
			filterMap[filterKey] = values[0]
			continue
		}

		if !reservedKeys[key] {
			filterMap[key] = values[0]
		}
	}

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 30
	}

	return PaginationParams{
		Page:    page,
		PerPage: perPage,
		SortBy:  sortBy,
		SortDir: sortDir,
		Search:  search,
		Filter:  filterMap,
	}
}

func PaginateQuery[T any, R any](
	db *gorm.DB,
	params PaginationParams,
	mapFunc func([]T) []R,
	searchableFields ...string,
) PaginationResult {
	var total int64
	var rawData []T

	wrapper := db.Session(&gorm.Session{}).Table("(?) as sub", db)

	if params.Search != "" && len(searchableFields) > 0 {
		searchQuery := "%" + params.Search + "%"

		var conditions []string
		var args []interface{}

		for _, field := range searchableFields {
			conditions = append(conditions, field+" LIKE ?")
			args = append(args, searchQuery)
		}

		whereClause := "(" + strings.Join(conditions, " OR ") + ")"
		wrapper = wrapper.Where(whereClause, args...)
	}

	countQuery := wrapper.Session(&gorm.Session{})
	if err := countQuery.Count(&total).Error; err != nil {
		LogError("[FAILED] Failed count data", err)
	}

	if params.SortBy == "" {
		params.SortBy = "(SELECT NULL)"
	}

	err := wrapper.
		Order(fmt.Sprintf("%s %s", params.SortBy, params.SortDir)).
		Offset((params.Page - 1) * params.PerPage).
		Limit(params.PerPage).
		Scan(&rawData).Error

	if err != nil {
		LogError("[FAILED] Failed paginate data", err)
	}
	mapped := mapFunc(rawData)
	if mapped == nil {
		mapped = make([]R, 0)
	}

	lastPage := int((total + int64(params.PerPage) - 1) / int64(params.PerPage))
	from := (params.Page-1)*params.PerPage + 1
	to := params.Page * params.PerPage
	if to > int(total) {
		to = int(total)
	}
	if total == 0 {
		from = 0
		to = 0
	}

	return PaginationResult{
		CurrentPage: params.Page,
		Data:        mapped,
		From:        from,
		LastPage:    lastPage,
		PerPage:     params.PerPage,
		To:          to,
		Total:       total,
	}
}

func ScanPaginateQuery[T any](
	query *gorm.DB,
	params PaginationParams,
	searchableFields ...string,
) PaginationResult {
	var total int64
	var datas []T

	if params.Search != "" && len(searchableFields) > 0 {
		searchQuery := "%" + params.Search + "%"
		var conditions []string
		var args []interface{}
		for _, field := range searchableFields {
			conditions = append(conditions, fmt.Sprintf(`%s LIKE ?`, field))
			args = append(args, searchQuery)
		}
		whereClause := strings.Join(conditions, " OR ")
		query = query.Where(whereClause, args...)
	}

	countQuery := query.Session(&gorm.Session{})

	err := countQuery.Count(&total).Error
	if err != nil {
		log.Println("Error counting data:", err)
	}

	if params.SortBy == "" {
		params.SortBy = "ROW_ID"
	}

	offset := (params.Page - 1) * params.PerPage

	err = query.
		Order(fmt.Sprintf("%s %s", params.SortBy, params.SortDir)).
		Offset(offset).
		Limit(params.PerPage).
		Scan(&datas).Error

	if err != nil {
		log.Println("Error scanning data:", err)
		datas = make([]T, 0)
	}

	lastPage := 0
	if total > 0 && params.PerPage > 0 {
		lastPage = int((total + int64(params.PerPage) - 1) / int64(params.PerPage))
	}
	from := 0
	if total > 0 {
		from = (params.Page-1)*params.PerPage + 1
	}
	to := from + len(datas) - 1
	if total == 0 {
		to = 0
	}

	return PaginationResult{
		CurrentPage: params.Page,
		Data:        datas,
		From:        from,
		LastPage:    lastPage,
		PerPage:     params.PerPage,
		To:          to,
		Total:       total,
	}
}

func PaginateQueryWithTotal[T any, R any](
	db *gorm.DB,
	params PaginationParams,
	total int64,
	mapFunc func([]T) []R,
	searchableFields ...string,
) PaginationResult {
	var rawData []T

	wrapper := db.Session(&gorm.Session{}).Table("(?) as sub", db)

	if params.Search != "" && len(searchableFields) > 0 {
		searchQuery := "%" + params.Search + "%"

		var conditions []string
		var args []interface{}

		for _, field := range searchableFields {
			conditions = append(conditions, field+" LIKE ?")
			args = append(args, searchQuery)
		}

		whereClause := "(" + strings.Join(conditions, " OR ") + ")"
		wrapper = wrapper.Where(whereClause, args...)
	}

	if params.SortBy == "" {
		params.SortBy = "(SELECT NULL)"
	}

	err := wrapper.
		Order(fmt.Sprintf("%s %s", params.SortBy, params.SortDir)).
		Offset((params.Page - 1) * params.PerPage).
		Limit(params.PerPage).
		Scan(&rawData).Error

	if err != nil {
		LogError("[FAILED] Failed paginate data", err)
	}
	mapped := mapFunc(rawData)
	if mapped == nil {
		mapped = make([]R, 0)
	}

	lastPage := int((total + int64(params.PerPage) - 1) / int64(params.PerPage))
	from := (params.Page-1)*params.PerPage + 1
	to := params.Page * params.PerPage
	if to > int(total) {
		to = int(total)
	}
	if total == 0 {
		from = 0
		to = 0
	}

	return PaginationResult{
		CurrentPage: params.Page,
		Data:        mapped,
		From:        from,
		LastPage:    lastPage,
		PerPage:     params.PerPage,
		To:          to,
		Total:       total,
	}
}
