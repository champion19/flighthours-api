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

		log.Success("register employee success",
			result.Employee.ToLogger(),
			"encoded_id", encodedID,
			"client_ip", c.ClientIP())

		// Record Prometheus metric for employee registration
		middleware.RecordEmployeeRegistration()
		h.Response.Success(c, domain.MsgUserRegistered)
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
				h.Response.Error(c, domain.MsgKCUserNotFound)
			case domain.ErrEmailAlreadyVerified:
				h.Response.Warning(c, domain.MsgKCEmailAlreadyVerified)
			default:
				h.Response.Error(c, domain.MsgKCVerifEmailError)
			}
			return
		}
		h.Response.Success(c, domain.MsgKCVerifEmailResent, req.Email)
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
		h.Response.Success(c, domain.MsgKCPwdResetSent)
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

			// Handle specific errors with appropriate messages
			switch err {
			case domain.ErrorEmailNotVerified:
				// Email not verified - verification email was resent automatically
				h.Response.Error(c, domain.MsgKCLoginEmailNotVerified)
			case domain.ErrUserNotFound:
				// User not found - return generic unauthorized for security
				h.Response.Error(c, domain.MsgUnauthorized)
			default:
				// Other errors (invalid credentials, etc.)
				h.Response.Error(c, domain.MsgUnauthorized)
			}
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
		h.Response.SuccessWithData(c, domain.MsgKCLoginSuccess, response)
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

		var req VerifyEmailRequest
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

		response := VerifyEmailResponse{
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

// GetEmployeeByID godoc
// @Summary      Obtener empleado por ID
// @Description  Obtiene la información de un empleado por su ID. Acepta tanto UUID como ID ofuscado. No expone la contraseña del usuario.
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID del empleado (UUID o ID ofuscado)"
// @Success      200  {object}  middleware.APIResponse{data=EmployeeResponse}  "Empleado encontrado"
// @Failure      400  {object}  middleware.APIResponse  "ID inválido"
// @Failure      404  {object}  middleware.APIResponse  "Empleado no encontrado"
// @Failure      500  {object}  middleware.APIResponse  "Error interno del servidor"
// @Router       /employees/{id} [get]
func (h handler) GetEmployeeByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		// Obtener el ID del path parameter
		inputID := c.Param("id")
		if inputID == "" {
			log.Error(logger.LogMessageIDDecodeError, "error", "empty id parameter", "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		log.Info(logger.LogEmployeeGetByID, "input_id", inputID, "client_ip", c.ClientIP())

		var employeeUUID string
		var responseID string

		// Detectar si es un UUID válido o un ID ofuscado
		if isValidUUID(inputID) {
			// Es un UUID directo, usarlo tal cual
			employeeUUID = inputID
			// Codificar el UUID para la respuesta (mantener consistencia)
			encodedID, err := h.EncodeID(inputID)
			if err != nil {
				log.Warn(logger.LogIDEncodeError, "uuid", inputID, "error", err)
				// Si no se puede codificar, usar el UUID en la respuesta
				responseID = inputID
			} else {
				responseID = encodedID
			}
			log.Debug(logger.LogEmployeeGetByID, "detected_format", "UUID", "uuid", employeeUUID)
		} else {
			// Es un ID ofuscado, decodificarlo
			uuid, err := h.DecodeID(inputID)
			if err != nil {
				h.HandleIDDecodingError(c, inputID, err)
				return
			}
			employeeUUID = uuid
			responseID = inputID // Mantener el ID ofuscado original
			log.Debug(logger.LogEmployeeGetByID, "detected_format", "encoded", "decoded_uuid", employeeUUID)
		}

		// Obtener el empleado del servicio
		employee, err := h.EmployeeService.GetEmployeeByID(c, employeeUUID)
		if err != nil {
			log.Error(logger.LogEmployeeGetByIDError, "uuid", employeeUUID, "error", err, "client_ip", c.ClientIP())
			switch err {
			case domain.ErrPersonNotFound, domain.ErrNotFoundUserById:
				h.Response.Error(c, domain.MsgUserNotFound)
			default:
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		if employee == nil {
			log.Warn(logger.LogEmployeeNotFound, "uuid", employeeUUID, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgUserNotFound)
			return
		}

		// Convertir a EmployeeResponse (sin contraseña) usando FromDomain
		response := FromDomain(employee, responseID)

		log.Success(logger.LogEmployeeGetByIDOK, "uuid", employeeUUID, "email", employee.Email, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgUserFound, response)
	}
}

// isValidUUID verifica si un string es un UUID válido
func isValidUUID(str string) bool {
	// UUID tiene formato: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx (36 caracteres con guiones)
	if len(str) != 36 {
		return false
	}
	// Verificar posiciones de los guiones
	if str[8] != '-' || str[13] != '-' || str[18] != '-' || str[23] != '-' {
		return false
	}
	// Verificar que los demás caracteres sean hexadecimales
	for i, c := range str {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			continue // Saltar guiones
		}
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

// UpdateEmployee godoc
// @Summary      Actualizar información de empleado
// @Description  Actualiza la información general de un empleado (nombre, airline, bp, fechas, rol, active).
// @Description  Los campos email y password NO se modifican ya que se manejan en endpoints separados.
// @Description  Si el campo active cambia, se sincroniza el estado enabled/disabled con Keycloak.
// @Description  Si el campo role cambia, se actualiza el rol asignado en Keycloak.
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        id       path      string                   true   "ID del empleado (UUID o ID ofuscado)"
// @Param        request  body      UpdateEmployeeRequest    true   "Datos a actualizar"
// @Success      200      {object}  middleware.APIResponse{data=UpdateEmployeeResponse}  "Empleado actualizado exitosamente"
// @Failure      400      {object}  middleware.APIResponse  "Error de validación - JSON inválido o ID inválido"
// @Failure      404      {object}  middleware.APIResponse  "Empleado no encontrado"
// @Failure      500      {object}  middleware.APIResponse  "Error interno del servidor"
// @Router       /employees/{id} [put]
func (h handler) UpdateEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		// Step 1: Get employee ID from path
		inputID := c.Param("id")
		if inputID == "" {
			log.Error(logger.LogMessageIDDecodeError, "error", "empty id parameter", "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		log.Info(logger.LogEmployeeUpdateRequest, "input_id", inputID, "client_ip", c.ClientIP())

		// Step 2: Decode ID (obfuscated or UUID)
		var employeeUUID string
		var responseID string

		if isValidUUID(inputID) {
			employeeUUID = inputID
			encodedID, err := h.EncodeID(inputID)
			if err != nil {
				log.Warn(logger.LogIDEncodeError, "uuid", inputID, "error", err)
				responseID = inputID
			} else {
				responseID = encodedID
			}
		} else {
			uuid, err := h.DecodeID(inputID)
			if err != nil {
				h.HandleIDDecodingError(c, inputID, err)
				return
			}
			employeeUUID = uuid
			responseID = inputID
		}

		// Step 3: Parse request body
		var req UpdateEmployeeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogRegJSONParseError, "error", err, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}

		// Step 4: Get current employee to preserve protected fields
		currentEmployee, err := h.EmployeeService.GetEmployeeByID(c, employeeUUID)
		if err != nil {
			log.Error(logger.LogEmployeeGetByIDError, "uuid", employeeUUID, "error", err, "client_ip", c.ClientIP())
			switch err {
			case domain.ErrPersonNotFound, domain.ErrNotFoundUserById:
				h.Response.Error(c, domain.MsgUserNotFound)
			default:
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		if currentEmployee == nil {
			log.Warn(logger.LogEmployeeNotFound, "uuid", employeeUUID, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgUserNotFound)
			return
		}

		// Step 5: Decode airline ID if it's obfuscated
		// The airline field can come as: empty, UUID, or obfuscated ID
		if req.Airline != "" {
			if !isValidUUID(req.Airline) {
				// It's an obfuscated ID, decode it to UUID
				decodedAirlineID, err := h.DecodeID(req.Airline)
				if err != nil {
					log.Error(logger.LogMessageIDDecodeError, "airline_id", req.Airline, "error", err, "client_ip", c.ClientIP())
					h.Response.Error(c, domain.MsgInvalidForeignKey)
					return
				}
				req.Airline = decodedAirlineID
				log.Debug(logger.LogEmployeeUpdateRequest, "decoded_airline", decodedAirlineID, "original", req.Airline)
			}
		}

		// Step 6: Build updated employee data (preserving email, password, keycloak_user_id)
		updatedEmployee := req.ToUpdateData(currentEmployee)

		// Step 6: Call interactor to update
		err = h.Interactor.UpdateEmployee(c, employeeUUID, updatedEmployee)
		if err != nil {
			log.Error(logger.LogEmployeeUpdateError, "uuid", employeeUUID, "error", err, "client_ip", c.ClientIP())
			switch err {
			case domain.ErrPersonNotFound, domain.ErrNotFoundUserById:
				h.Response.Error(c, domain.MsgUserNotFound)
			case domain.ErrUserCannotUpdate:
				h.Response.Error(c, domain.MsgUserUpdateError)
			case domain.ErrKeycloakUpdateFailed:
				h.Response.Error(c, domain.MsgUserKeycloakUpdateError)
			case domain.ErrRoleUpdateFailed:
				h.Response.Error(c, domain.MsgUserRoleUpdateError)
			// Data validation errors - 400/422, not 500
			case domain.ErrInvalidForeignKey:
				h.Response.Error(c, domain.MsgInvalidForeignKey)
			case domain.ErrDataTooLong:
				h.Response.Error(c, domain.MsgDataTooLong)
			case domain.ErrDuplicateUser:
				h.Response.Error(c, domain.MsgUserDuplicate)
			case domain.ErrInvalidData:
				h.Response.Error(c, domain.MsgInvalidData)
			default:
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		// Step 7: Return success response
		response := UpdateEmployeeResponse{
			ID:      responseID,
			Updated: true,
		}

		log.Success(logger.LogEmployeeUpdateComplete, "uuid", employeeUUID, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgUserUpdated, response)
	}
}

// ChangePassword godoc
// @Summary      Cambiar contraseña de usuario autenticado
// @Description  Permite a un usuario cambiar su contraseña conociendo la contraseña actual. Este flujo no requiere salir de la API ni tokens por email.
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        request  body      ChangePasswordRequest  true  "Email, contraseña actual y nueva contraseña"
// @Success      200      {object}  middleware.APIResponse{data=ChangePasswordResponse}  "Contraseña cambiada exitosamente"
// @Failure      400      {object}  middleware.APIResponse  "Error de validación - Contraseñas no coinciden"
// @Failure      401      {object}  middleware.APIResponse  "Contraseña actual incorrecta"
// @Failure      404      {object}  middleware.APIResponse  "Usuario no encontrado"
// @Failure      500      {object}  middleware.APIResponse  "Error interno del servidor"
// @Router       /auth/change-password [post]
func (h handler) ChangePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		var req ChangePasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogRegJSONParseError, "error", err)
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}

		// Validate new passwords match
		if req.NewPassword != req.ConfirmPassword {
			log.Warn(logger.LogKeycloakChangePasswordMismatch, "email", req.Email, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgKCPwdChangeNewMismatch)
			return
		}

		log.Info(logger.LogKeycloakChangePassword, "email", req.Email, "client_ip", c.ClientIP())

		email, err := h.Interactor.ChangePassword(c, req.Email, req.CurrentPassword, req.NewPassword, req.ConfirmPassword)
		if err != nil {
			switch err {
			case domain.ErrInvalidCurrentPassword:
				h.Response.Error(c, domain.MsgKCPwdCurrentInvalid)
			case domain.ErrUserNotFound:
				h.Response.Error(c, domain.MsgKCUserNotFound)
			case domain.ErrPasswordMismatch:
				h.Response.Error(c, domain.MsgKCPwdChangeNewMismatch)
			case domain.ErrPasswordUpdateFailed:
				h.Response.Error(c, domain.MsgKCPwdChangeError)
			default:
				h.Response.Error(c, domain.MsgKCPwdChangeError)
			}
			return
		}

		response := ChangePasswordResponse{
			Changed: true,
			Email:   email,
		}

		log.Success(logger.LogKeycloakChangePasswordOK, "email", email, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgKCPwdChanged, response)
	}
}

// DeleteEmployee godoc
// @Summary      Eliminar empleado
// @Description  Elimina un empleado del sistema (BD y Keycloak). Esta operación es irreversible.
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID del empleado (UUID o ID ofuscado)"
// @Success      200  {object}  middleware.APIResponse  "Empleado eliminado exitosamente"
// @Failure      400  {object}  middleware.APIResponse  "ID inválido"
// @Failure      404  {object}  middleware.APIResponse  "Empleado no encontrado"
// @Failure      500  {object}  middleware.APIResponse  "Error interno del servidor"
// @Router       /employees/{id} [delete]
func (h handler) DeleteEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		// Get ID from URL parameter
		inputID := c.Param("id")
		if inputID == "" {
			log.Error(logger.LogMessageIDDecodeError, "error", "empty id parameter", "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		log.Info(logger.LogEmployeeDeleting, "input_id", inputID, "client_ip", c.ClientIP())

		var employeeUUID string

		// Detect if it's a valid UUID or an obfuscated ID
		if isValidUUID(inputID) {
			employeeUUID = inputID
		} else {
			// Try to decode obfuscated ID
			decodedID, err := h.IDEncoder.Decode(inputID)
			if err != nil {
				log.Error(logger.LogMessageIDDecodeError, "input_id", inputID, "error", err, "client_ip", c.ClientIP())
				h.Response.Error(c, domain.MsgValIDInvalid)
				return
			}
			employeeUUID = decodedID
		}

		log.Debug(logger.LogEmployeeDeletingDB, "employee_id", employeeUUID, "original_id", inputID, "client_ip", c.ClientIP())

		// Call interactor to delete
		if err := h.Interactor.DeleteEmployee(c, employeeUUID); err != nil {
			switch err {
			case domain.ErrPersonNotFound:
				h.Response.Error(c, domain.MsgUserNotFound)
			case domain.ErrUserCannotDelete:
				h.Response.Error(c, domain.MsgUserCannotDelete)
			default:
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		log.Success(logger.LogEmployeeDeleteComplete, "employee_id", employeeUUID, "client_ip", c.ClientIP())
		h.Response.Success(c, domain.MsgUserDeleted)
	}
}

// GetMe godoc
// @Summary      Obtener información del empleado autenticado
// @Description  Obtiene la información del empleado actualmente autenticado usando el token JWT. No expone la contraseña.
// @Tags         employees
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  middleware.APIResponse{data=EmployeeResponse}  "Empleado encontrado"
// @Failure      401  {object}  middleware.APIResponse  "No autenticado"
// @Failure      500  {object}  middleware.APIResponse  "Error interno del servidor"
// @Router       /employees/me [get]
func (h handler) GetMe() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		// Get authenticated employee from JWT middleware
		employee, ok := middleware.GetAuthenticatedUser(c)
		if !ok || employee == nil {
			log.Error(logger.LogEmployeeGetByIDError, "error", "no authenticated user in context", "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgUnauthorized)
			return
		}

		log.Info(logger.LogEmployeeGetByID, "employee_id", employee.ID, "email", employee.Email, "client_ip", c.ClientIP())

		// Encode ID for response
		encodedID, err := h.EncodeID(employee.ID)
		if err != nil {
			log.Warn(logger.LogIDEncodeError, "uuid", employee.ID, "error", err)
			encodedID = employee.ID // Use raw ID if encoding fails
		}

		// Convert to EmployeeResponse (without password)
		response := FromDomain(employee, encodedID)

		log.Success(logger.LogEmployeeGetByIDOK, "employee_id", employee.ID, "email", employee.Email, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgUserFound, response)
	}
}

// UpdateMe godoc
// @Summary      Actualizar información del empleado autenticado
// @Description  Actualiza la información del empleado actualmente autenticado. Los campos email y password NO se modifican.
// @Tags         employees
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      UpdateEmployeeRequest    true   "Datos a actualizar"
// @Success      200      {object}  middleware.APIResponse{data=UpdateEmployeeResponse}  "Empleado actualizado exitosamente"
// @Failure      400      {object}  middleware.APIResponse  "Error de validación"
// @Failure      401      {object}  middleware.APIResponse  "No autenticado"
// @Failure      500      {object}  middleware.APIResponse  "Error interno del servidor"
// @Router       /employees/me [put]
func (h handler) UpdateMe() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		// Get authenticated employee from JWT middleware
		currentEmployee, ok := middleware.GetAuthenticatedUser(c)
		if !ok || currentEmployee == nil {
			log.Error(logger.LogEmployeeUpdateError, "error", "no authenticated user in context", "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgUnauthorized)
			return
		}

		employeeUUID := currentEmployee.ID
		log.Info(logger.LogEmployeeUpdateRequest, "employee_id", employeeUUID, "email", currentEmployee.Email, "client_ip", c.ClientIP())

		// Parse request body
		var req UpdateEmployeeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogRegJSONParseError, "error", err, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}

		// Decode airline ID if it's obfuscated
		// The airline field can come as: empty, UUID, or obfuscated ID
		if req.Airline != "" {
			if !isValidUUID(req.Airline) {
				// It's an obfuscated ID, decode it to UUID
				decodedAirlineID, err := h.DecodeID(req.Airline)
				if err != nil {
					log.Error(logger.LogMessageIDDecodeError, "airline_id", req.Airline, "error", err, "client_ip", c.ClientIP())
					h.Response.Error(c, domain.MsgInvalidForeignKey)
					return
				}
				req.Airline = decodedAirlineID
				log.Debug(logger.LogEmployeeUpdateRequest, "decoded_airline", decodedAirlineID, "original", req.Airline)
			}
		}

		// Build updated employee data (preserving email, password, keycloak_user_id)
		updatedEmployee := req.ToUpdateData(currentEmployee)

		// Call interactor to update
		err := h.Interactor.UpdateEmployee(c, employeeUUID, updatedEmployee)
		if err != nil {
			log.Error(logger.LogEmployeeUpdateError, "employee_id", employeeUUID, "error", err, "client_ip", c.ClientIP())
			switch err {
			case domain.ErrPersonNotFound, domain.ErrNotFoundUserById:
				h.Response.Error(c, domain.MsgUserNotFound)
			case domain.ErrUserCannotUpdate:
				h.Response.Error(c, domain.MsgUserUpdateError)
			case domain.ErrKeycloakUpdateFailed:
				h.Response.Error(c, domain.MsgUserKeycloakUpdateError)
			case domain.ErrRoleUpdateFailed:
				h.Response.Error(c, domain.MsgUserRoleUpdateError)
			case domain.ErrInvalidForeignKey:
				h.Response.Error(c, domain.MsgInvalidForeignKey)
			case domain.ErrDataTooLong:
				h.Response.Error(c, domain.MsgDataTooLong)
			case domain.ErrDuplicateUser:
				h.Response.Error(c, domain.MsgUserDuplicate)
			case domain.ErrInvalidData:
				h.Response.Error(c, domain.MsgInvalidData)
			default:
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		// Encode ID for response
		encodedID, err := h.EncodeID(employeeUUID)
		if err != nil {
			log.Warn(logger.LogIDEncodeError, "uuid", employeeUUID, "error", err)
			encodedID = employeeUUID
		}

		response := UpdateEmployeeResponse{
			ID:      encodedID,
			Updated: true,
		}

		log.Success(logger.LogEmployeeUpdateComplete, "employee_id", employeeUUID, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgUserUpdated, response)
	}
}

// DeleteMe godoc
// @Summary      Eliminar cuenta del empleado autenticado
// @Description  Elimina la cuenta del empleado actualmente autenticado (BD y Keycloak). Esta operación es irreversible.
// @Tags         employees
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  middleware.APIResponse  "Cuenta eliminada exitosamente"
// @Failure      401  {object}  middleware.APIResponse  "No autenticado"
// @Failure      500  {object}  middleware.APIResponse  "Error interno del servidor"
// @Router       /employees/me [delete]
func (h handler) DeleteMe() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		// Get authenticated employee from JWT middleware
		employee, ok := middleware.GetAuthenticatedUser(c)
		if !ok || employee == nil {
			log.Error(logger.LogEmployeeDeleteDBError, "error", "no authenticated user in context", "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgUnauthorized)
			return
		}

		employeeUUID := employee.ID
		log.Info(logger.LogEmployeeDeleting, "employee_id", employeeUUID, "email", employee.Email, "client_ip", c.ClientIP())

		// Call interactor to delete
		if err := h.Interactor.DeleteEmployee(c, employeeUUID); err != nil {
			switch err {
			case domain.ErrPersonNotFound:
				h.Response.Error(c, domain.MsgUserNotFound)
			case domain.ErrUserCannotDelete:
				h.Response.Error(c, domain.MsgUserCannotDelete)
			default:
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		log.Success(logger.LogEmployeeDeleteComplete, "employee_id", employeeUUID, "client_ip", c.ClientIP())
		h.Response.Success(c, domain.MsgUserDeleted)
	}
}
