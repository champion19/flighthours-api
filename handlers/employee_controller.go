package handlers

import (
	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/middleware"
	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/gin-gonic/gin"
)

// RegisterEmployee godoc
// @Summary      Registrar nueva cuenta de empleado
// @Description  Crea una nueva cuenta de empleado en el sistema con sincronización a Keycloak. Incluye validación de datos, verificación de duplicados y creación de usuario en el sistema de autenticación.
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        account  body      EmployeeRequest  true  "Datos del empleado a registrar"
// @Success      201      {object}  middleware.APIResponse{data=RegisterEmployeeResponse}  "Cuenta creada exitosamente"
// @Failure      400      {object}  middleware.APIResponse  "Error de validación - Datos inválidos o incompletos"
// @Failure      409      {object}  middleware.APIResponse  "Conflicto - Email o número de identidad ya registrado"
// @Failure      500      {object}  middleware.APIResponse  "Error interno del servidor"
// @Router       /register [post]
func (h handler) RegisterEmployee() func(c *gin.Context) {

	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogRegRequestReceived,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		var employeeRequest EmployeeRequest
		if err := c.ShouldBindJSON(&employeeRequest); err != nil {
			log.Error(logger.LogRegJSONParseError, "error", err, "client_ip", c.ClientIP())
			c.Error(domain.ErrInvalidJSONFormat)
			return
		}

		log.Info(logger.LogRegProcessing,
			"email", employeeRequest.Email,
			"role", employeeRequest.Role)

		result, err := h.Interactor.RegisterEmployee(c, employeeRequest.ToDomain())
		if err != nil {
			log.Error(logger.LogRegProcessError,
				"email", employeeRequest.Email,
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(err)
			return
		}

		// Ofuscar el ID antes de exponerlo en la API
		encodedID, err := h.EncodeID(result.Employee.ID)
		if err != nil {
			h.HandleIDEncodingError(c, result.Employee.ID, err)
			return
		}

		//TODO;TENERLO EN CUENTA, ESTO ES DE COOKIES HTTTPONLY
		c.SetCookie(
			"employee_id",        // name
			result.Employee.ID,   // value
			3600,                 // expira en 1 hora
			"/",                  // path
			c.Request.Host,       // domain
			c.Request.TLS != nil, // secure
			true,                 // httpOnly
		)

		// Construir respuesta con HATEOAS
		baseURL := GetBaseURL(c)
		links := BuildAccountLinks(baseURL, encodedID)
		SetLocationHeader(c, baseURL, "accounts", encodedID)

		response := RegisterEmployeeResponse{
			Links: links,
		}

		log.Success("register employee success",
			result.Employee.ToLogger(),
			"encoded_id", encodedID,
			"client_ip", c.ClientIP())

		// Record Prometheus metric for employee registration
		middleware.RecordEmployeeRegistration()
		h.Response.SuccessWithData(c, domain.MsgUserRegistered, response)
	}
}

// ResendVerificationEmail godoc
// @Summary      Reenviar email de verificación
// @Description  Reenvía el email de verificación a un usuario registrado que no ha verificado su cuenta
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        request  body      ResendVerificationEmailRequest  true  "Email del usuario"
// @Success      200      {object}  middleware.APIResponse{data=ResendVerificationEmailResponse}  "Email reenviado exitosamente"
// @Failure      400      {object}  middleware.APIResponse  "Error de validación - Email inválido"
// @Failure      404      {object}  middleware.APIResponse  "Usuario no encontrado"
// @Failure      500      {object}  middleware.APIResponse  "Error interno del servidor"
// @Router       /auth/resend-verification [post]
func (h handler) ResendVerificationEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ResendVerificationEmailRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}

		err := h.Interactor.ResendVerificationEmail(c, req.Email)
		if err != nil {
			// Manejar diferentes tipos de errores
			switch err {
			case domain.ErrUserNotFound:
				h.Response.Error(c, "MOD_KC_USER_NOT_FOUND_ERR_00001")
			case domain.ErrEmailAlreadyVerified:
				h.Response.Warning(c, "MOD_KC_EMAIL_ALREADY_VERIFIED_WARN_00001")
			default:
				h.Response.Error(c, "MOD_KC_VERIF_EMAIL_ERROR_ERR_00001")
			}
			return
		}
		h.Response.Success(c, "MOD_KC_VERIF_EMAIL_RESENT_EXI_00001", req.Email)
	}
}

// RequestPasswordReset godoc
// @Summary      Solicitar recuperación de contraseña
// @Description  Envía un email con instrucciones para recuperar la contraseña de un usuario
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        request  body      PasswordResetRequest  true  "Email del usuario"
// @Success      200      {object}  middleware.APIResponse{data=PasswordResetResponse}  "Email de recuperación enviado"
// @Failure      400      {object}  middleware.APIResponse  "Error de validación - Email inválido"
// @Failure      404      {object}  middleware.APIResponse  "Usuario no encontrado"
// @Failure      500      {object}  middleware.APIResponse  "Error interno del servidor"
// @Router       /auth/password-reset [post]
func (h handler) RequestPasswordReset() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req PasswordResetRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}

		// Este método SIEMPRE retorna nil por seguridad (no revela si el email existe)
		// El logging interno sí registra el resultado real
		_ = h.Interactor.RequestPasswordReset(c, req.Email)

		// Siempre responder con éxito genérico
		h.Response.Success(c, "MOD_KC_PWD_RESET_SENT_EXI_00001")
	}
}

// @Summary Login de usuario
// @Description Autentica un usuario y retorna tokens de acceso
// @Tags Autenticación
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Credenciales de login"
// @Success 200 {object} middleware.APIResponse{data=LoginResponse} "Login exitoso"
// @Failure 400 {object} middleware.APIResponse "Credenciales inválidas"
// @Failure 401 {object} middleware.APIResponse "Email no verificado o credenciales incorrectas"
// @Failure 500 {object} middleware.APIResponse "Error interno del servidor"
// @Router /auth/login [post]
func (h handler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogRegJSONParseError, "error", err)
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}

		log.Info(logger.LogKeycloakUserLogin, "email", req.Email, "client_ip", c.ClientIP())

		token, err := h.Interactor.Login(c, req.Email, req.Password)
		if err != nil {
			log.Error(logger.LogKeycloakUserLoginError, "email", req.Email, "error", err, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgUnauthorized)
			return
		}

		response := LoginResponse{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			ExpiresIn:    token.ExpiresIn,
			TokenType:    token.TokenType,
		}

		log.Success(logger.LogKeycloakUserLoginOK, "email", req.Email, "client_ip", c.ClientIP())
		middleware.RecordEmployeeRegistration()
		h.Response.SuccessWithData(c, "MOD_AUTH_LOGIN_SUCCESS_EXI_00001", response)
	}
}

// @Summary Verificar email de usuario (Proxy)
// @Description Verifica el email de un usuario usando un token JWT. Este endpoint actúa como proxy para no exponer Keycloak directamente.
// @Tags Autenticación
// @Accept json
// @Produce json
// @Param request body VerifyEmailRequest true "Token de verificación del email"
// @Success 200 {object} middleware.APIResponse{data=VerifyEmailResponse} "Email verificado exitosamente"
// @Failure 400 {object} middleware.APIResponse "Token inválido o expirado"
// @Failure 404 {object} middleware.APIResponse "Usuario no encontrado"
// @Failure 409 {object} middleware.APIResponse "Email ya estaba verificado"
// @Failure 500 {object} middleware.APIResponse "Error interno del servidor"
// @Router /auth/verify-email [post]
func (h handler) VerifyEmailByToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		var req verifyEmailRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogRegJSONParseError, "error", err)
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}

		log.Info(logger.LogKeycloakEmailVerify, "client_ip", c.ClientIP())

		// Pasar el token al Interactor - la extracción del email se hace en la capa de negocio
		email, err := h.Interactor.VerifyEmailByToken(c, req.Token)
		if err != nil {
			switch err {
			case domain.ErrInvalidToken:
				h.Response.Error(c, domain.MsgKCInvalidToken)
			case domain.ErrUserNotFound:
				h.Response.Error(c, domain.MsgKCUserNotFound)
			case domain.ErrEmailAlreadyVerified:
				h.Response.Warning(c, domain.MsgKCEmailAlreadyVerified)
			default:
				h.Response.Error(c, domain.MsgKCEmailVerifyError)
			}
			return
		}

		response := verifyEmailResponse{
			Verified: true,
			Email:    email,
		}

		log.Success(logger.LogKeycloakEmailVerifyOK, "email", email, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgKCEmailVerified, response)
	}
}

// UpdatePassword godoc
// @Summary      Actualizar contraseña con token
// @Description  Actualiza la contraseña del usuario usando el token recibido por email
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        request  body      UpdatePasswordRequest  true  "Token y nueva contraseña"
// @Success      200      {object}  middleware.APIResponse{data=UpdatePasswordResponse}  "Contraseña actualizada exitosamente"
// @Failure      400      {object}  middleware.APIResponse  "Error de validación - Contraseñas no coinciden o token inválido"
// @Failure      401      {object}  middleware.APIResponse  "Token inválido o expirado"
// @Failure      500      {object}  middleware.APIResponse  "Error interno del servidor"
// @Router       /auth/update-password [post]
func (h handler) UpdatePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		var req UpdatePasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogRegJSONParseError, "error", err)
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}

		// Validate passwords match
		if req.NewPassword != req.ConfirmPassword {
			log.Warn(logger.LogKeycloakPasswordMismatch, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgKCPwdMismatch)
			return
		}

		log.Info(logger.LogKeycloakPasswordUpdate, "client_ip", c.ClientIP())

		email, err := h.Interactor.UpdatePassword(c, req.Token, req.NewPassword, req.ConfirmPassword)
		if err != nil {
			switch err {
			case domain.ErrInvalidToken:
				h.Response.Error(c, domain.MsgKCPwdUpdateTokenInvalid)
			case domain.ErrPasswordMismatch:
				h.Response.Error(c, domain.MsgKCPwdMismatch)
			case domain.ErrPasswordUpdateFailed:
				h.Response.Error(c, domain.MsgKCPwdUpdateError)
			default:
				h.Response.Error(c, domain.MsgKCPwdUpdateError)
			}
			return
		}

		response := UpdatePasswordResponse{
			Updated: true,
			Email:   email,
		}

		log.Success(logger.LogKeycloakPasswordUpdateOK, "email", email, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgKCPwdUpdated, response)
	}
}
