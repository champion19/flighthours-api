package domain

import "errors"

var (
	ErrDuplicateUser             = errors.New("Err_DUPLICATE_USER")
	ErrUserCannotSave            = errors.New("Err_USER_CANNOT_SAVE")
	ErrPersonNotFound            = errors.New("Err_PERSON_NOT_FOUND")
	ErrUserNotFound              = errors.New("Err_USER_NOT_FOUND")
	ErrEmailAlreadyVerified      = errors.New("Err_EMAIL_ALREADY_VERIFIED")
	ErrInvalidTransaction        = errors.New("Err_INVALID_TRANSACTION")
	ErrGettingUserByEmail        = errors.New("Err_GETTING_USER_BY_EMAIL")
	ErrNotFoundUserByEmail       = errors.New("Err_NOT_FOUND_USER_BY_EMAIL")
	ErrNotFoundUserById          = errors.New("Err_NOT_FOUND_USER_BY_ID")
	ErrUserCannotFound           = errors.New("Err_USER_CANNOT_FOUND")
	ErrUserCannotGet             = errors.New("Err_USER_CANNOT_GET")
	ErrorEmailNotVerified        = errors.New("Err_EMAIL_NOT_VERIFIED")
	ErrVerificationTokenNotFound = errors.New("Err_VERIFICATION_TOKEN_NOT_FOUND")
	ErrTokenExpired              = errors.New("Err_TOKEN_EXPIRED")
	ErrTokenAlreadyUsed          = errors.New("Err_TOKEN_ALREADY_USED")
	ErrInvalidToken              = errors.New("Err_INVALID_TOKEN")
	ErrRegistrationFailed        = errors.New("Err_REGISTRATION_FAILED")
	ErrRoleRequired              = errors.New("Err_ROLE_REQUIRED")
	ErrDBQueryFailed             = errors.New("Err_DB_QUERY_FAILED")
	ErrUserCannotDelete          = errors.New("Err_USER_CANNOT_DELETE")
	ErrPasswordMismatch          = errors.New("Err_PASSWORD_MISMATCH")
	ErrPasswordUpdateFailed      = errors.New("Err_PASSWORD_UPDATE_FAILED")
	ErrUserCannotUpdate          = errors.New("Err_USER_CANNOT_UPDATE")
	ErrKeycloakUpdateFailed      = errors.New("Err_KEYCLOAK_UPDATE_FAILED")
	ErrRoleUpdateFailed          = errors.New("Err_ROLE_UPDATE_FAILED")
	// Data validation errors (from DB constraints)
	ErrInvalidForeignKey = errors.New("Err_INVALID_FOREIGN_KEY") // 1452 - FK constraint fails
	ErrDataTooLong       = errors.New("Err_DATA_TOO_LONG")       // 1406 - Data too long for column
	ErrInvalidData       = errors.New("Err_INVALID_DATA")        // Generic invalid data
)

// Infrastructure Errors
var (
	ErrKeycloakInconsistentState  = errors.New("ERR_KC_INCONSISTENT_STATE")
	ErrKeycloakUserCreationFailed = errors.New("ERR_KC_USER_CREATION_FAILED")
	ErrKeycloakCleanupFailed      = errors.New("ERR_KC_CLEANUP_FAILED")
	// Dependency availability errors
	ErrKeycloakUnavailable = errors.New("ERR_KC_UNAVAILABLE")
	ErrDatabaseUnavailable = errors.New("ERR_DB_UNAVAILABLE")
	// Specific user existence errors
	ErrKeycloakUserExists     = errors.New("ERR_KC_USER_EXISTS")
	ErrDatabaseUserExists     = errors.New("ERR_DB_USER_EXISTS")
	ErrIncompleteRegistration = errors.New("ERR_INCOMPLETE_REGISTRATION")
)

// Request Validation errors
var (
	ErrInvalidJSONFormat = errors.New("Err_INVALID_JSON_FORMAT")
	ErrInvalidRequest    = errors.New("Err_INVALID_REQUEST")
	ErrInvalidID         = errors.New("Err_INVALID_ID")
	ErrInternalServer    = errors.New("Err_INTERNAL_SERVER")

	// Schema validation errors
	ErrSchemaBadRequest       = errors.New("ERR_SCHEMA_BAD_REQUEST")
	ErrSchemaInvalidRequest   = errors.New("ERR_SCHEMA_INVALID_REQUEST")
	ErrSchemaReadFailed       = errors.New("ERR_SCHEMA_READ_FAILED")
	ErrSchemaEmpty            = errors.New("ERR_SCHEMA_EMPTY")
	ErrSchemaCompileFailed    = errors.New("ERR_SCHEMA_COMPILE_FAILED")
	ErrSchemaValidationFailed = errors.New("ERR_SCHEMA_VALIDATION_FAILED")
	ErrSchemaBodyReadFailed   = errors.New("ERR_SCHEMA_BODY_READ_FAILED")
	ErrSchemaFieldFormat      = errors.New("ERR_SCHEMA_FIELD_FORMAT")
	ErrSchemaFieldRequired    = errors.New("ERR_SCHEMA_FIELD_REQUIRED")
	ErrSchemaFieldType        = errors.New("ERR_SCHEMA_FIELD_TYPE")
	ErrSchemaMultipleFields   = errors.New("ERR_SCHEMA_MULTIPLE_FIELDS")
)

// Authorization Errors
var (
	ErrRoleAssignmentFailed = errors.New("Err_ROLE_ASSIGNMENT_FAILED")
	ErrRoleRemovalFailed    = errors.New("Err_ROLE_REMOVAL_FAILED")
	ErrRoleCheckFailed      = errors.New("Err_ROLE_CHECK_FAILED")
	ErrGetUserRolesFailed   = errors.New("Err_GET_USER_ROLES_FAILED")
)

// Message Management Errors (MOD_M_*)
var (
	ErrMessageNotFound         = errors.New("ERR_MESSAGE_NOT_FOUND")
	ErrMessageCodeRequired     = errors.New("ERR_MESSAGE_CODE_REQUIRED")
	ErrMessageTypeRequired     = errors.New("ERR_MESSAGE_TYPE_REQUIRED")
	ErrMessageTitleRequired    = errors.New("ERR_MESSAGE_TITLE_REQUIRED")
	ErrMessageContentRequired  = errors.New("ERR_MESSAGE_CONTENT_REQUIRED")
	ErrMessageModuleRequired   = errors.New("ERR_MESSAGE_MODULE_REQUIRED")
	ErrMessageCategoryRequired = errors.New("ERR_MESSAGE_CATEGORY_REQUIRED")
	ErrMessageCodeDuplicate    = errors.New("ERR_MESSAGE_CODE_DUPLICATE")
	ErrMessageCannotSave       = errors.New("ERR_MESSAGE_CANNOT_SAVE")
	ErrMessageCannotUpdate     = errors.New("ERR_MESSAGE_CANNOT_UPDATE")
	ErrMessageCannotDelete     = errors.New("ERR_MESSAGE_CANNOT_DELETE")
	ErrMessageInvalidType      = errors.New("ERR_MESSAGE_INVALID_TYPE")
	ErrMessageListFailed       = errors.New("ERR_MESSAGE_LIST_FAILED")
	ErrMessageNotRegistered    = errors.New("ERR_MESSAGE_NOT_REGISTERED")
	ErrMessageInactive         = errors.New("ERR_MESSAGE_INACTIVE")
)

// ============================================
// MESSAGE CODES - Constants for use in code
// ============================================

// User Module (MOD_U_*)
const (
	MsgUserDuplicate        = "MOD_U_DUP_ERR_00001"
	MsgUserCannotSave       = "MOD_U_SAVE_ERR_00002"
	MsgUserNotFound         = "MOD_U_GET_ERR_00003"
	MsgUserEmailError       = "MOD_U_EMAIL_ERR_00004"
	MsgUserEmailNotFound    = "MOD_U_EMAIL_NF_ERR_00005"
	MsgUserEmailNotVerified = "MOD_U_EMAIL_NV_ERR_00006"
	MsgUserTokenNotFound    = "MOD_U_TOKEN_NF_ERR_00007"
	MsgUserTokenExpired     = "MOD_U_TOKEN_EXP_ERR_00008"
	MsgUserTokenUsed        = "MOD_U_TOKEN_USED_ERR_00009"
	MsgUserRegError         = "MOD_U_REG_ERR_00010"
	MsgUserRoleRequired     = "MOD_U_ROLE_REQ_ERR_00011"
	MsgUserCannotDelete     = "MOD_U_DEL_ERR_00012"

	MsgUserRegistered    = "MOD_U_REG_EXI_00001"
	MsgUserUpdated       = "MOD_U_UPD_EXI_00002"
	MsgUserDeleted       = "MOD_U_DEL_EXI_00003"
	MsgUserEmailVerified = "MOD_U_VER_EXI_00004"
	MsgUserFound         = "MOD_U_GET_EXI_00005"

	// Update errors
	MsgUserUpdateError         = "MOD_U_UPD_ERR_00013"
	MsgUserKeycloakUpdateError = "MOD_U_KC_UPD_ERR_00014"
	MsgUserRoleUpdateError     = "MOD_U_ROLE_UPD_ERR_00015"
	// Data validation errors (from DB constraints) - These are 400/422, not 500
	MsgInvalidForeignKey = "MOD_V_FK_ERR_00014"   // Invalid reference (e.g., airline doesn't exist)
	MsgDataTooLong       = "MOD_V_LEN_ERR_00015"  // Data exceeds column length
	MsgInvalidData       = "MOD_V_DATA_ERR_00016" // Generic invalid data
)

// Person Module (MOD_P_*)
const (
	MsgPersonNotFound     = "MOD_P_NOT_FOUND_ERR_00001"
	MsgPersonInvalidTx    = "MOD_P_TRANS_ERR_00002"
	MsgPersonRegistered   = "MOD_P_REG_EXI_00001"
	MsgPersonUpdated      = "MOD_P_UPD_EXI_00002"
	MsgPersonCannotDelete = "MOD_P_DEL_ERR_00003"
)

// Validation Module (MOD_V_*)
const (
	MsgValBadFormat     = "MOD_V_VAL_ERR_00001"
	MsgValInvalidReq    = "MOD_V_VAL_ERR_00002"
	MsgValSchemaRead    = "MOD_V_VAL_ERR_00003"
	MsgValSchemaEmpty   = "MOD_V_VAL_ERR_00004"
	MsgValSchemaCompile = "MOD_V_VAL_ERR_00005"
	MsgValFailed        = "MOD_V_VAL_ERR_00006"
	MsgValBodyRead      = "MOD_V_VAL_ERR_00007"
	MsgValFieldFormat   = "MOD_V_VAL_ERR_00008"
	MsgValFieldRequired = "MOD_V_VAL_ERR_00009"
	MsgValFieldType     = "MOD_V_VAL_ERR_00010"
	MsgValMultiple      = "MOD_V_VAL_ERR_00011"
	MsgValJSONInvalid   = "MOD_V_JSON_ERR_00012"
	MsgValIDInvalid     = "MOD_V_ID_ERR_00013"
)

// Authorization Module (MOD_A_*)
const (
	MsgRoleAssignError = "MOD_A_ROLE_ASSIGN_ERR_00001"
	MsgRoleRemoveError = "MOD_A_ROLE_REMOVE_ERR_00002"
	MsgRoleCheckError  = "MOD_A_ROLE_CHECK_ERR_00003"
	MsgRoleGetError    = "MOD_A_ROLE_GET_ERR_00004"
	MsgRoleAssigned    = "MOD_A_ROLE_ASSIGN_EXI_00001"
	MsgRoleRemoved     = "MOD_A_ROLE_REMOVE_EXI_00002"
)

// General Module (GEN_*)
const (
	MsgServerError   = "GEN_SRV_ERR_00001"
	MsgUnauthorized  = "GEN_AUTH_ERR_00002"
	MsgForbidden     = "GEN_FORBIDDEN_ERR_00003"
	MsgOpSuccess     = "GEN_OPE_EXI_00001"
	MsgInfoProcess   = "GEN_INFO_00001"
	MsgWarningAction = "GEN_WARN_00001"
)

// Message Module (MOD_M_*)
const (
	MsgMessageNotFound       = "MOD_M_NOT_FOUND_ERR_00001"
	MsgMessageCodeRequired   = "MOD_M_CODE_REQ_ERR_00002"
	MsgMessageTypeRequired   = "MOD_M_TYPE_REQ_ERR_00003"
	MsgMessageTitleRequired  = "MOD_M_TITLE_REQ_ERR_00004"
	MsgMessageContentReq     = "MOD_M_CONTENT_REQ_ERR_00005"
	MsgMessageModuleRequired = "MOD_M_MODULE_REQ_ERR_00006"
	MsgMessageCategoryReq    = "MOD_M_CATEGORY_REQ_ERR_00007"
	MsgMessageCodeDuplicate  = "MOD_M_CODE_DUP_ERR_00008"
	MsgMessageSaveError      = "MOD_M_SAVE_ERR_00009"
	MsgMessageUpdateError    = "MOD_M_UPDATE_ERR_00010"
	MsgMessageDeleteError    = "MOD_M_DELETE_ERR_00011"
	MsgMessageInvalidType    = "MOD_M_TYPE_INV_ERR_00012"
	MsgMessageListError      = "MOD_M_LIST_ERR_00013"
	MsgMessageCannotDelete   = "MOD_M_DEL_ERR_00014"
	MsgMessageNotRegistered  = "MOD_M_NOT_REG_ERR_00015"
	MsgMessageInactive       = "MOD_M_INACTIVE_ERR_00016"

	MsgMessageCreated = "MOD_M_CREATE_EXI_00001"
	MsgMessageUpdated = "MOD_M_UPDATE_EXI_00002"
	MsgMessageDeleted = "MOD_M_DELETE_EXI_00003"
	MsgMessageListed  = "MOD_M_LIST_EXI_00004"
)

// Infrastructure Module (MOD_INFRA_*)
const (
	MsgKeycloakInconsistentState = "MOD_INFRA_KC_INCONSISTENT_ERR_00001"
	MsgKeycloakCreateError       = "MOD_INFRA_KC_CREATE_ERR_00002"
	MsgKeycloakCleanupError      = "MOD_INFRA_KC_CLEANUP_ERR_00003"
	// Dependency availability messages
	MsgKeycloakUnavailable = "MOD_INFRA_KC_UNAVAIL_ERR_00004"
	MsgDatabaseUnavailable = "MOD_INFRA_DB_UNAVAIL_ERR_00005"
	MsgDependencyFailure   = "MOD_INFRA_DEP_FAIL_ERR_00006"
	// Specific user existence messages
	MsgKeycloakUserExists     = "MOD_INFRA_KC_USER_EXISTS_ERR_00007"
	MsgDatabaseUserExists     = "MOD_INFRA_DB_USER_EXISTS_ERR_00008"
	MsgIncompleteRegistration = "MOD_INFRA_INCOMPLETE_REG_ERR_00009"
)

// Keycloak Module (MOD_KC_*) - Email Verification and Password Reset
const (
	// Email Verification
	MsgKCEmailVerified        = "MOD_KC_EMAIL_VERIFIED_EXI_00001"
	MsgKCInvalidToken         = "MOD_KC_INVALID_TOKEN_ERR_00001"
	MsgKCEmailVerifyError     = "MOD_KC_EMAIL_VERIFY_ERROR_ERR_00001"
	MsgKCUserNotFound         = "MOD_KC_USER_NOT_FOUND_ERR_00001"
	MsgKCEmailAlreadyVerified = "MOD_KC_EMAIL_ALREADY_VERIFIED_WARN_00001"
	// Verification Email Sending
	MsgKCVerifEmailSent  = "MOD_KC_VERIF_EMAIL_SENT_EXI_00001"
	MsgKCVerifEmailError = "MOD_KC_VERIF_EMAIL_ERROR_ERR_00001"
	// Password Reset
	MsgKCPwdResetSent  = "MOD_KC_PWD_RESET_SENT_EXI_00001"
	MsgKCPwdResetError = "MOD_KC_PWD_RESET_ERROR_ERR_00001"
	// Password Update
	MsgKCPwdUpdated            = "MOD_KC_PWD_UPDATED_EXI_00001"
	MsgKCPwdUpdateError        = "MOD_KC_PWD_UPDATE_ERROR_ERR_00001"
	MsgKCPwdMismatch           = "MOD_KC_PWD_MISMATCH_ERR_00001"
	MsgKCPwdUpdateTokenInvalid = "MOD_KC_PWD_UPDATE_TOKEN_INVALID_ERR_00001"
	// Login
	MsgKCLoginEmailNotVerified = "MOD_KC_LOGIN_EMAIL_NOT_VERIFIED_ERR_00001"
	MsgKCLoginSuccess          = "MOD_KC_LOGIN_SUCCESS_EXI_00001"
)
