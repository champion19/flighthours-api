package middleware

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidJSONFormat           = errors.New("invalid JSON format")
	ErrUnmarshalBody               = errors.New("failed to process request body")
	ErrSchemaValidation            = errors.New("validation failed")
	ErrInternalServer              = errors.New("internal server error")
	ErrModuleRootNotFound          = errors.New("could not find module root")
	ErrSchemaFileNotFound          = errors.New("schema file not found")
	ErrSchemaFileRead              = errors.New("failed to read schema file")
	ErrSchemaCompilation           = errors.New("failed to compile JSON schema")
	ErrSchemaEmpty                 = errors.New("JSON schema is empty or null")
	ErrValidatorInitFailed         = errors.New("validator initialization failed")
	ErrValidationUserFailed        = errors.New("user validation failed")
	ErrValidationUserNotFound      = errors.New("user not found")
	ErrValidationUserAlreadyExists = errors.New("user already exists")
)

type SchemaError struct {
	Code    string
	Message string
	err     error
}

func (e *SchemaError) Error() string {
	return e.Message
}

func (e *SchemaError) Unwrap() error {
	return e.err
}

func NewSchemaError(code, message string) *SchemaError {
	return &SchemaError{
		Code:    code,
		Message: message,
		err:     errors.New(message),
	}
}

// NewFieldSchemaError creates a schema error with field-specific information
func NewFieldSchemaError(baseError *SchemaError, fieldName string) *SchemaError {
	message := fmt.Sprintf("%s: %s", baseError.Message, fieldName)
	return &SchemaError{
		Code:    baseError.Code,
		Message: message,
		err:     errors.New(message),
	}
}

// NewMultipleFieldSchemaError creates a schema error for multiple fields
func NewMultipleFieldSchemaError(fieldNames []string) *SchemaError {
	message := "Errores de validación en los campos: " + fieldNames[0]
	for i := 1; i < len(fieldNames); i++ {
		message += ", " + fieldNames[i]
	}
	return &SchemaError{
		Code:    "MOD_V_VAL_ERR_00011",
		Message: message,
		err:     errors.New(message),
	}
}
