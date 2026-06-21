package utils

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"reflect"
	"strings"
	"virtual_account_api/constants"
	"virtual_account_api/resources"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ValidationUtils struct {
	db *gorm.DB
}

func NewValidationUtils(db *gorm.DB) *ValidationUtils {
	v := &ValidationUtils{db: db}

	_ = Validate.RegisterValidation("exists", func(fl validator.FieldLevel) bool {
		return v.validateExists(fl)
	})
	_ = Validate.RegisterValidation("not_exists", func(fl validator.FieldLevel) bool {
		return v.validateNotExists(fl)
	})
	_ = Validate.RegisterValidation("nefield", func(fl validator.FieldLevel) bool { return neField(fl) })

	return v
}

var Validate = validator.New()

func ValidateAndBind[T any](c *gin.Context) (*T, *resources.GeneralResponse[any], int) {
	var req T

	// Debug: Print raw form values
	if c.Request.Method == "POST" || c.Request.Method == "PUT" {
		c.Request.ParseMultipartForm(32 << 20) // 32MB
	}

	if err := c.ShouldBind(&req); err != nil {
		fmt.Println("Error binding request:", err)
		return nil, &resources.GeneralResponse[any]{
			BaseResponse: resources.BaseResponse{
				Status:       constants.StatusErrorValidation,
				ResponseCode: constants.CodeErrorSendMidTier,
				Message:      constants.StatusIncompleteData,
			},
		}, http.StatusBadRequest
	}

	injectFilesFromRequest(c, &req)

	if err := Validate.Struct(req); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			return nil, &resources.GeneralResponse[any]{
				BaseResponse: resources.BaseResponse{
					ResponseCode: constants.CodeErrorSendMidTier,
					Message:      constants.StatusErrorCustom + validationMsg(e),
				},
			}, http.StatusBadRequest
		}
	}

	if err := validateAllFileFields(&req); err != nil {
		return nil, &resources.GeneralResponse[any]{
			BaseResponse: resources.BaseResponse{
				ResponseCode: constants.CodeErrorSendMidTier,
				Message:      constants.StatusErrorCustom + err.Error(),
			},
		}, http.StatusBadRequest
	}

	return &req, nil, http.StatusOK
}

func injectFilesFromRequest[T any](c *gin.Context, req *T) {
	rv := reflect.ValueOf(req).Elem()
	rt := reflect.TypeOf(req).Elem()

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		formTag := field.Tag.Get("form")

		if field.Type == reflect.TypeOf((*multipart.FileHeader)(nil)) && formTag != "" {
			if file, err := c.FormFile(formTag); err == nil && file != nil {
				rv.Field(i).Set(reflect.ValueOf(file))
			}
		}
	}
}

func (u *ValidationUtils) validateExists(fl validator.FieldLevel) bool {
	param := fl.Param()
	parts := strings.Split(param, "#")
	if len(parts) != 2 {
		return false
	}

	table, column := parts[0], parts[1]
	value := fl.Field().String()

	var exists bool
	err := u.db.
		Raw(fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s = ?)", table, column), value).
		Scan(&exists).Error

	return err == nil && exists
}

func (u *ValidationUtils) validateNotExists(fl validator.FieldLevel) bool {
	param := fl.Param()
	parts := strings.Split(param, "#")
	if len(parts) != 2 {
		return false
	}

	table, column := parts[0], parts[1]
	value := fl.Field().String()

	var exists bool
	err := u.db.
		Raw(fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s = ?)", table, column), value).
		Scan(&exists).Error

	return err == nil && !exists
}

func validateAllFileFields(req interface{}) error {
	rv := reflect.ValueOf(req).Elem()
	rt := reflect.TypeOf(req).Elem()

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		value := rv.Field(i)

		if !isMultipartFileField(field) {
			continue
		}

		formName := getFormName(field)
		validateTag := field.Tag.Get("Validate")
		extTag := field.Tag.Get("file_ext")
		maxTag := field.Tag.Get("file_max")

		file, _ := value.Interface().(*multipart.FileHeader)

		if err := validateRequired(file, validateTag, formName); err != nil {
			return err
		}
		if err := validateFileExtension(file, extTag, formName); err != nil {
			return err
		}
		if err := validateFileMax(file, maxTag, formName); err != nil {
			return err
		}
	}
	return nil
}

func validateRequired(file *multipart.FileHeader, tag, fieldName string) error {
	if strings.Contains(tag, "required") && (file == nil || file.Size == 0) {
		return fmt.Errorf("%s is required and must be a valid file", fieldName)
	}
	return nil
}

func validateFileExtension(file *multipart.FileHeader, extTag string, fieldName string) error {
	if file == nil || file.Size == 0 || extTag == "" {
		return nil
	}

	allowed := strings.Split(extTag, "|")
	dotIdx := strings.LastIndex(file.Filename, ".")
	if dotIdx < 0 {
		return fmt.Errorf("%s must have a valid file extension", fieldName)
	}

	ext := strings.ToLower(file.Filename[dotIdx+1:])
	for _, allow := range allowed {
		if ext == strings.ToLower(strings.TrimSpace(allow)) {
			return nil
		}
	}
	return fmt.Errorf("%s must be one of the following file types: %v", fieldName, allowed)
}

func validateFileMax(file *multipart.FileHeader, maxTag string, fieldName string) error {
	if file == nil || file.Size == 0 || maxTag == "" {
		return nil
	}

	var maxKB int64
	_, err := fmt.Sscanf(maxTag, "%d", &maxKB)
	if err != nil || maxKB <= 0 {
		return fmt.Errorf("%s has invalid max file size tag", fieldName)
	}

	if file.Size > maxKB*1024 {
		return fmt.Errorf("%s must not be larger than %d KB", fieldName, maxKB)
	}
	return nil
}

func isMultipartFileField(field reflect.StructField) bool {
	return field.Type == reflect.TypeOf((*multipart.FileHeader)(nil))
}

func getFormName(field reflect.StructField) string {
	formName := field.Tag.Get("form")
	if formName == "" {
		formName = field.Name
	}
	return formName
}

func neField(fl validator.FieldLevel) bool {
	fieldName := fl.Param()
	otherField := fl.Parent().FieldByName(fieldName)

	if !otherField.IsValid() {
		return false
	}

	return fl.Field().Interface() != otherField.Interface()
}

func validationMsg(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("Field '%s' is required", e.Field())
	case "email":
		return fmt.Sprintf("Field '%s' must be a valid email address", e.Field())
	case "min":
		if e.Kind().String() == "slice" || e.Kind().String() == "array" {
			return fmt.Sprintf("Field '%s' must have at least %s items", e.Field(), e.Param())
		}
		return fmt.Sprintf("Field '%s' must be at least %s", e.Field(), e.Param())
	case "max":
		if e.Kind().String() == "slice" || e.Kind().String() == "array" {
			return fmt.Sprintf("Field '%s' must have no more than %s items", e.Field(), e.Param())
		}
		return fmt.Sprintf("Field '%s' must not be more than %s", e.Field(), e.Param())
	case "len":
		if e.Kind().String() == "slice" || e.Kind().String() == "array" {
			return fmt.Sprintf("Field '%s' must contain exactly %s items", e.Field(), e.Param())
		}
		return fmt.Sprintf("Field '%s' must be exactly %s characters", e.Field(), e.Param())
	case "uuid":
		return fmt.Sprintf("Field '%s' must be a valid UUID", e.Field())
	case "datetime":
		return fmt.Sprintf("Field '%s' must match datetime format '%s'", e.Field(), e.Param())
	case "exists":
		return fmt.Sprintf("Field '%s' must refer to an existing record in the database", e.Field())
	case "not_exists":
		return fmt.Sprintf("Field '%s' must not refer to an existing record in the database", e.Field())
	case "oneof":
		return fmt.Sprintf("Field '%s' must be one of the following: %s", e.Field(), e.Param())
	case "numeric":
		return fmt.Sprintf("Field '%s' must be a numeric value", e.Field())
	case "dive":
		return fmt.Sprintf("Invalid data in field '%s'", e.Field())
	case "nefield":
		return fmt.Sprintf("Field '%s' must not be equal to field '%s'", e.Field(), e.Param())
	default:
		return "Invalid value: " + e.Error()
	}
}
