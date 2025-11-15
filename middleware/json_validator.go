package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	json_schema "github.com/champion19/flighthours-api/platform/schema"
	"github.com/gin-gonic/gin"
	"github.com/kaptinlin/jsonschema"
)

type Builder struct {
	Validators *json_schema.Validators
	isLogin    bool
}

func NewMiddlewareValidator(validators *json_schema.Validators) *Builder {

	return &Builder{
		Validators: validators,
	}
}

func (b *Builder) WithValidateRegister() gin.HandlerFunc {
	b.isLogin = false
	return b.jsonValidator(b.Validators.RegisterValidator)
}

func (b *Builder) jsonValidator(schema *jsonschema.Schema) gin.HandlerFunc {
	return func(c *gin.Context) {

		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Error(json_schema.NewFieldSchemaError(json_schema.ErrBodyReadFailed, err.Error()))
			c.Abort()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		var data map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			c.Error(json_schema.ErrBadRequest)
			c.Abort()
			return
		}
		result := schema.Validate(data)
		if !result.IsValid() {
			var fieldNames []string


			for _, validationError := range result.Errors {
				if validationError.Params != nil {
					if properties, exists := validationError.Params["properties"]; exists {
						propertiesStr := fmt.Sprintf("%v", properties)
						fields := strings.Split(propertiesStr, ",")
						for _, field := range fields {
							field = strings.TrimSpace(field)
							field = strings.Trim(field, "'\"")
							if field != "" {
								fieldNames = append(fieldNames, field)
							}
						}
					} else if property, exists := validationError.Params["property"]; exists {
						propertyName := fmt.Sprintf("%v", property)
						propertyName = strings.Trim(propertyName, "'\"")
						if propertyName != "" {
							fieldNames = append(fieldNames, propertyName)
						}
					}
				}
			}

			var schemaError *json_schema.SchemaError

			if len(fieldNames) == 1 {
				var firstError *jsonschema.EvaluationError
				for _, err := range result.Errors {
					firstError = err
					break
				}

				fieldName := fieldNames[0]
				switch firstError.Code {
				case "property_mismatch":
					schemaError = json_schema.NewFieldSchemaError(json_schema.ErrFieldPropertyMismatch, fieldName)
				case "required":
					schemaError = json_schema.NewFieldSchemaError(json_schema.ErrFieldRequired, fieldName)
				case "type":
					schemaError = json_schema.NewFieldSchemaError(json_schema.ErrFieldTypeInvalid, fieldName)
				default:
					schemaError = json_schema.NewFieldSchemaError(json_schema.ErrValidationFailed, fieldName)
				}
			} else if len(fieldNames) > 1 {
				schemaError = json_schema.NewMultipleFieldSchemaError(fieldNames)
			} else {
				schemaError = json_schema.ErrValidationFailed
			}

			c.Error(schemaError)
			c.Abort()
			return
		}
		c.Next()
	}
}
