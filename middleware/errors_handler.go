package middleware

import (
	"net/http"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/gin-gonic/gin"
)

var mapError = map[error]ErrorResponse{
	domain.ErrDuplicateUser: {
		Code:    "MODE_U_USU_ERR_00001",
		Message: "User already exists",
		Status:  http.StatusConflict,
	},
	domain.ErrUserCannotSave: {
		Code:    "MODE_U_USU_ERR_00002",
		Message: "User cannot be saved",
		Status:  http.StatusInternalServerError,
	},
	domain.ErrPersonNotFound: {
		Code:    "MODE_U_USU_ERR_00003",
		Message: "Person not found",
		Status:  http.StatusNotFound,
	},
	domain.ErrGettingUserByEmail: {
		Code:    "MODE_U_USU_ERR_00004",
		Message: "Error getting user by email",
		Status:  http.StatusInternalServerError,
	},
	domain.ErrNotFoundUserByEmail: {
		Code:    "MODE_U_USU_ERR_00005",
		Message: "User not found by email",
		Status:  http.StatusNotFound,
	},
	domain.ErrNotFoundUserById: {
		Code:    "MODE_U_USU_ERR_00006",
		Message: "User not found by id",
		Status:  http.StatusNotFound,
	},
	domain.ErrUserCannotFound: {
		Code:    "MODE_U_USU_ERR_00007",
		Message: "User cannot be found",
		Status:  http.StatusNotFound,
	},
	domain.ErrUserCannotGet: {
		Code:    "MODE_U_USU_ERR_00008",
		Message: "User cannot be retrieved",
		Status:  http.StatusInternalServerError,
	},
	domain.ErrorEmailNotVerified: {
		Code:    "MODE_U_USU_ERR_00009",
		Message: "Email not verified",
		Status:  http.StatusUnauthorized,
	},
	domain.ErrVerificationTokenNotFound: {
		Code:    "MODE_U_USU_ERR_00010",
		Message: "Verification token not found",
		Status:  http.StatusNotFound,
	},
	domain.ErrTokenExpired: {
		Code:    "MODE_U_USU_ERR_00011",
		Message: "Token expired",
		Status:  http.StatusUnauthorized,
	},
	domain.ErrTokenAlreadyUsed: {
		Code:    "MODE_U_USU_ERR_00012",
		Message: "Token already used",
		Status:  http.StatusForbidden,
	},
	domain.ErrRegistrationFailed: {
		Code:    "MODE_U_USU_ERR_00013",
		Message: "Registration failed",
		Status:  http.StatusBadRequest,
	},
	domain.ErrRoleRequired: {
		Code:    "MODE_U_USU_ERR_00014",
		Message: "Role required",
		Status:  http.StatusBadRequest,
	},
	domain.ErrInvalidJSONFormat: {
		Code:    "MODE_U_USU_ERR_00015",
		Message: "Invalid JSON format",
		Status:  http.StatusBadRequest,
	},
	domain.ErrInvalidRequest: {
		Code:    "MODE_U_USU_ERR_00016",
		Message: "Invalid request parameters",
		Status:  http.StatusBadRequest,
	},
	domain.ErrRoleAssignmentFailed: {
		Code:    "MODE_U_USU_ERR_00017",
		Message: "Error assigning role",
		Status:  http.StatusInternalServerError,
	},
	domain.ErrRoleRemovalFailed: {
		Code:    "MODE_U_USU_ERR_00018",
		Message: "Error removing role",
		Status:  http.StatusInternalServerError,
	},
	domain.ErrRoleCheckFailed: {
		Code:    "MODE_U_USU_ERR_00019",
		Message: "Error checking role",
		Status:  http.StatusInternalServerError,
	},
	domain.ErrGetUserRolesFailed: {
		Code:    "MODE_U_USU_ERR_00020",
		Message: "Error retrieving user roles",
		Status:  http.StatusInternalServerError,
	},
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			if response, ok := mapError[err]; ok {
				c.JSON(response.Status, response)
				return
			}

			c.JSON(http.StatusInternalServerError, map[string]any{
				"success": false,
				"message": err.Error(),
			})

		}
	}

}
