package domain

import "net/http"

const (
	CodeDuplicateUser             = "ERR_DUPLICATE_USER"
	CodeUserCannotSave            = "ERR_USER_CANNOT_SAVE"
	CodePersonNotFound            = "ERR_PERSON_NOT_FOUND"
	CodeInvalidTransaction        = "ERR_INVALID_TRANSACTION"
	CodeGettingUserByEmail        = "ERR_GETTING_USER_BY_EMAIL"
	CodeNotFoundUserByEmail       = "ERR_NOT_FOUND_USER_BY_EMAIL"
	CodeNotFoundUserById          = "ERR_NOT_FOUND_USER_BY_ID"
	CodeUserCannotFound           = "ERR_USER_CANNOT_FOUND"
	CodeUserCannotGet             = "ERR_USER_CANNOT_GET"
	CodeEmailNotVerified          = "ERR_EMAIL_NOT_VERIFIED"
	CodeVerificationTokenNotFound = "ERR_VERIFICATION_TOKEN_NOT_FOUND"
	CodeTokenExpired              = "ERR_TOKEN_EXPIRED"
	CodeTokenAlreadyUsed          = "ERR_TOKEN_ALREADY_USED"
	CodeRegistrationFailed        = "ERR_REGISTRATION_FAILED"
	CodeRoleRequired              = "ERR_ROLE_REQUIRED"
	CodeDBQueryFailed             = "ERR_DB_QUERY_FAILED"
	CodeInvalidJSONFormat         = "ERR_INVALID_JSON_FORMAT"
	CodeInvalidRequest            = "ERR_INVALID_REQUEST"
	CodeInvalidID                 = "ERR_INVALID_ID"
	CodeRoleAssignmentFailed      = "ERR_ROLE_ASSIGNMENT_FAILED"
	CodeRoleRemovalFailed         = "ERR_ROLE_REMOVAL_FAILED"
	CodeRoleCheckFailed           = "ERR_ROLE_CHECK_FAILED"
	CodeGetUserRolesFailed        = "ERR_GET_USER_ROLES_FAILED"
)

// MessageManager global - se inicializa en el bootstrap
var Messages *MessageManager

var (
	ErrDuplicateUser             *AppError
	ErrUserCannotSave            *AppError
	ErrPersonNotFound            *AppError
	ErrInvalidTransaction        *AppError
	ErrGettingUserByEmail        *AppError
	ErrNotFoundUserByEmail       *AppError
	ErrNotFoundUserById          *AppError
	ErrUserCannotFound           *AppError
	ErrUserCannotGet             *AppError
	ErrorEmailNotVerified        *AppError
	ErrVerificationTokenNotFound *AppError
	ErrTokenExpired              *AppError
	ErrTokenAlreadyUsed          *AppError
	ErrRegistrationFailed        *AppError
	ErrRoleRequired              *AppError
	ErrDBQueryFailed             *AppError
	ErrInvalidJSONFormat         *AppError
	ErrInvalidRequest            *AppError
	ErrInvalidID                 *AppError
	ErrRoleAssignmentFailed      *AppError
	ErrRoleRemovalFailed         *AppError
	ErrRoleCheckFailed           *AppError
	ErrGetUserRolesFailed        *AppError
)

func InitErrors() {

	ErrDuplicateUser = NewAppError(CodeDuplicateUser, http.StatusConflict)
	ErrUserCannotSave = NewAppError(CodeUserCannotSave, http.StatusInternalServerError)
	ErrPersonNotFound = NewAppError(CodePersonNotFound, http.StatusNotFound)
	ErrInvalidTransaction = NewAppError(CodeInvalidTransaction, http.StatusBadRequest)
	ErrGettingUserByEmail = NewAppError(CodeGettingUserByEmail, http.StatusInternalServerError)
	ErrNotFoundUserByEmail = NewAppError(CodeNotFoundUserByEmail, http.StatusNotFound)
	ErrNotFoundUserById = NewAppError(CodeNotFoundUserById, http.StatusNotFound)
	ErrUserCannotFound = NewAppError(CodeUserCannotFound, http.StatusNotFound)
	ErrUserCannotGet = NewAppError(CodeUserCannotGet, http.StatusInternalServerError)
	ErrorEmailNotVerified = NewAppError(CodeEmailNotVerified, http.StatusForbidden)
	ErrVerificationTokenNotFound = NewAppError(CodeVerificationTokenNotFound, http.StatusNotFound)
	ErrTokenExpired = NewAppError(CodeTokenExpired, http.StatusUnauthorized)
	ErrTokenAlreadyUsed = NewAppError(CodeTokenAlreadyUsed, http.StatusBadRequest)
	ErrRegistrationFailed = NewAppError(CodeRegistrationFailed, http.StatusInternalServerError)
	ErrRoleRequired = NewAppError(CodeRoleRequired, http.StatusBadRequest)
	ErrDBQueryFailed = NewAppError(CodeDBQueryFailed, http.StatusInternalServerError)
	ErrInvalidJSONFormat = NewAppError(CodeInvalidJSONFormat, http.StatusBadRequest)
	ErrInvalidRequest = NewAppError(CodeInvalidRequest, http.StatusBadRequest)
	ErrInvalidID = NewAppError(CodeInvalidID, http.StatusBadRequest)
	ErrRoleAssignmentFailed = NewAppError(CodeRoleAssignmentFailed, http.StatusInternalServerError)
	ErrRoleRemovalFailed = NewAppError(CodeRoleRemovalFailed, http.StatusInternalServerError)
	ErrRoleCheckFailed = NewAppError(CodeRoleCheckFailed, http.StatusInternalServerError)
	ErrGetUserRolesFailed = NewAppError(CodeGetUserRolesFailed, http.StatusInternalServerError)
}
