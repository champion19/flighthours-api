package domain

import "errors"

type DomainError struct {
	Code    string
	Message string
	err     error
}

func (e *DomainError) Error() string {
	return e.Message
}

func (e *DomainError) Unwrap() error {
	return e.err
}

func NewDomainError(code, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		err:     errors.New(message),
	}
}

var (
	ErrDuplicateUser             = NewDomainError("MOD_U_USU_ERR_00001", "User already exists")
	ErrUserCannotSave            = NewDomainError("MOD_U_USU_ERR_00002", "User cannot be saved")
	ErrPersonNotFound            = NewDomainError("MOD_U_USU_ERR_00003", "Person not found")
	ErrGettingUserByEmail        = NewDomainError("MOD_U_USU_ERR_00004", "Error getting user by email")
	ErrNotFoundUserByEmail       = NewDomainError("MOD_U_USU_ERR_00005", "User not found by email")
	ErrNotFoundUserById          = NewDomainError("MOD_U_USU_ERR_00006", "User not found by id")
	ErrUserCannotFound           = NewDomainError("MOD_U_USU_ERR_00007", "User cannot be found")
	ErrUserCannotGet             = NewDomainError("MOD_U_USU_ERR_00008", "User cannot be retrieved")
	ErrorEmailNotVerified        = NewDomainError("MOD_U_USU_ERR_00009", "Email not verified")
	ErrVerificationTokenNotFound = NewDomainError("MOD_U_USU_ERR_000010", "Verification token not found")
	ErrTokenExpired              = NewDomainError("MOD_U_USU_ERR_00011", "Token expired")
	ErrTokenAlreadyUsed          = NewDomainError("MOD_U_USU_ERR_00012", "Token already used")
	ErrRegistrationFailed        = NewDomainError("MOD_U_USU_ERR_00013", "Registration failed")
	ErrRoleRequired              = NewDomainError("MOD_U_USU_ERR_00014", "Role required")
)

var (
	ErrInvalidJSONFormat = NewDomainError("MOD_V_VAL_ERR_00001", "Invalid JSON format")
	ErrInvalidRequest    = NewDomainError("MOD_V_VAL_ERR_00002", "Invalid request parameters")
)

var (
	ErrRoleAssignmentFailed = NewDomainError("MOD_A_AUT_ERR_00001", "Error assigning role")
	ErrRoleRemovalFailed    = NewDomainError("MOD_A_AUT_ERR_00002", "Error removing role")
	ErrRoleCheckFailed      = NewDomainError("MOD_A_AUT_ERR_00003", "Error checking role")
	ErrGetUserRolesFailed   = NewDomainError("MOD_A_AUT_ERR_00004", "Error retrieving user roles")
)
