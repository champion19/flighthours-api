package domain

import "errors"

var (
	ErrDuplicateUser             = errors.New("user already exists")
	ErrUserCannotSave            = errors.New("user cannot be saved")
	ErrPersonNotFound            = errors.New("person not found")
	
	ErrGettingUserByEmail        = errors.New("error getting user by email")
	ErrNotFoundUserByEmail       = errors.New("user not found by email")
	ErrNotFoundUserById          = errors.New("user not found by id")
	ErrUserCannotFound           = errors.New("user cannot be found")
	ErrUserCannotGet             = errors.New("user cannot be retrieved")
	ErrorEmailNotVerified        = errors.New("email not verified")
	ErrVerificationTokenNotFound = errors.New("verification token not found")
	ErrTokenExpired              = errors.New("token expired")
	ErrTokenAlreadyUsed          = errors.New("token already used")
	ErrRegistrationFailed        = errors.New("registration failed")
	ErrRoleRequired              = errors.New("role required")
)

var (
	ErrInvalidJSONFormat = errors.New("invalid json format")
	ErrInvalidRequest    = errors.New("invalid request parameters")
)

var (
	ErrRoleAssignmentFailed = errors.New("error assigning role")
	ErrRoleRemovalFailed    = errors.New("error removing role")
	ErrRoleCheckFailed      = errors.New("error checking role")
	ErrGetUserRolesFailed   = errors.New("error retrieving user roles")
)
