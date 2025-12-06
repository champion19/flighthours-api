package logger

// ============================================
// CENTRALIZED LOG MESSAGES
// ============================================
// All log messages are defined here as constants to:
// 1. Avoid hardcoded strings across the codebase
// 2. Enable easy maintenance and updates
// 3. Support future internationalization
// 4. Ensure consistency

// ============================================
// APPLICATION LIFECYCLE
// ============================================
const (
	LogAppStarting          = "Iniciando aplicación FlighHours Backend"
	LogAppConfigLoaded      = "Configuración cargada exitosamente"
	LogAppConfigError       = "Error cargando configuración"
	LogAppDatabaseConnected = "Conexión a base de datos establecida"
	LogAppDatabaseError     = "Error conectando a base de datos"
	LogAppDatabasePingOK    = "Ping a base de datos exitoso"
	LogAppDatabasePingError = "Error haciendo ping a base de datos"
	LogAppServerStarting    = "Servidor iniciando"
	LogAppServerListening   = "Servidor escuchando"
	LogAppServerStartError  = "Error iniciando servidor"
	LogAppShuttingDown      = "Apagando aplicación gracefully"
)

// ============================================
// REGISTRATION / AUTHENTICATION
// ============================================
const (
	LogRegRequestReceived   = "Nueva solicitud de registro recibida"
	LogRegProcessing        = "Procesando registro de usuario"
	LogRegJSONParseError    = "Error parseando JSON del request"
	LogRegJSONBindError     = "Error vinculando JSON del request"
	LogRegProcessError      = "Error en proceso de registro"
	LogRegIDEncodeError     = "Error ofuscando ID de usuario"
	LogRegSuccess           = "Usuario registrado exitosamente"
	LogRegKeycloakSync      = "Sincronizando usuario con Keycloak"
	LogRegKeycloakSyncError = "Error sincronizando con Keycloak"
)

// ============================================
// DATABASE OPERATIONS
// ============================================
const (
	LogDBQueryExecuting        = "Ejecutando query de base de datos"
	LogDBQuerySuccess          = "Query ejecutado exitosamente"
	LogDBQueryError            = "Error ejecutando query"
	LogDBTransactionStart      = "Iniciando transacción de base de datos"
	LogDBTransactionBegin      = "Comenzando transacción"
	LogDBTransactionBeginErr   = "Error comenzando transacción"
	LogDBTransactionCommit     = "Commit de transacción exitoso"
	LogDBTransactionCommitErr  = "Error haciendo commit de transacción"
	LogDBTransactionRollback   = "Rollback de transacción"
	LogDBConnectionPoolInfo    = "Información de connection pool"
	LogDBConnectionError       = "Error conectando a base de datos"
	LogDBConnectionEstablished = "Conexión a base de datos establecida"
)

// ============================================
// MESSAGING / CACHE
// ============================================
const (
	LogMsgCacheInit            = "Inicializando cache de mensajes"
	LogMsgCacheLoaded          = "Mensajes del sistema cargados en cache"
	LogMsgCacheLoadError       = "Error cargando mensajes del sistema desde BD"
	LogMsgCacheRefreshStart    = "Iniciando auto-refresh de cache de mensajes"
	LogMsgCacheRefreshDisabled = "Auto-refresh de cache deshabilitado"
	LogMsgCacheRefreshing      = "Auto-refrescando cache de mensajes desde BD"
	LogMsgCacheRefreshOK       = "Cache de mensajes refrescado exitosamente"
	LogMsgCacheRefreshError    = "Error durante auto-refresh de cache"
	LogMsgCacheRefreshStop     = "Deteniendo auto-refresh de cache de mensajes"
	LogMsgNotInCache           = "Mensaje no encontrado en cache, cargando desde BD"
	LogMsgNotInDB              = "Mensaje no encontrado en base de datos"
	LogMsgCachedFromDB         = "Mensaje cargado desde BD y cacheado"
	LogMsgInactive             = "Mensaje existe pero está desactivado"
)

// ============================================
// ROUTING / MIDDLEWARE
// ============================================
const (
	LogRouteConfiguring      = "Configurando rutas de la aplicación"
	LogRouteConfigured       = "Rutas configuradas correctamente"
	LogRouteValidatorInit    = "Inicializando validador de schemas"
	LogRouteValidatorOK      = "Validador de schema inicializado"
	LogRouteValidatorError   = "Error creando validador de schema"
	LogMiddlewareErrorCaught = "Error de negocio capturado"
	LogMiddlewareInternalErr = "Error interno del servidor"
)

// ============================================
// VALIDATION
// ============================================
const (
	LogValidationStart   = "Iniciando validación de request"
	LogValidationOK      = "Validación de request exitosa"
	LogValidationFailed  = "Validación de request fallida"
	LogValidationDetails = "Detalles de validación"
)

// ============================================
// KEYCLOAK / EXTERNAL SERVICES
// ============================================
const (
	LogKeycloakClientInit          = "Inicializando cliente Keycloak"
	LogKeycloakClientOK            = "Cliente Keycloak inicializado correctamente"
	LogKeycloakClientError         = "Error inicializando cliente Keycloak"
	LogKeycloakClientCreated       = "Cliente Keycloak creado exitosamente"
	LogKeycloakConfigNil           = "Configuración de Keycloak no puede ser nil"
	LogKeycloakAdminAuth           = "Autenticando admin de Keycloak"
	LogKeycloakAdminAuthError      = "Error autenticando admin de Keycloak"
	LogKeycloakAdminTokenInit      = "Inicializando token de admin"
	LogKeycloakAdminTokenInitError = "Error inicializando token de admin"
	LogKeycloakTokenRefresh        = "Refrescando token de admin de Keycloak"
	LogKeycloakTokenRefreshOK      = "Token de admin refrescado exitosamente"
	LogKeycloakTokenRefreshErr     = "Error refrescando token de admin de Keycloak"
	LogKeycloakTokenEnsure         = "Asegurando token válido"
	LogKeycloakTokenEnsureError    = "Error asegurando token válido"
	LogKeycloakUserLogin           = "Intentando login de usuario"
	LogKeycloakUserLoginOK         = "Login de usuario exitoso"
	LogKeycloakUserLoginError      = "Error en login de usuario"
	LogKeycloakUserLoginFailed     = "Login de usuario falló"
	LogKeycloakUserCreate          = "Creando usuario en Keycloak"
	LogKeycloakUserCreateOK        = "Usuario creado en Keycloak"
	LogKeycloakUserCreateError     = "Error creando usuario en Keycloak"
	LogKeycloakUserNil             = "Usuario no puede ser nil"
	LogKeycloakUserGet             = "Obteniendo usuario de Keycloak"
	LogKeycloakUserGetByEmail      = "Obteniendo usuario por email"
	LogKeycloakUserGetByEmailError = "Error obteniendo usuario por email"
	LogKeycloakUserGetByID         = "Obteniendo usuario por ID"
	LogKeycloakUserGetByIDError    = "Error obteniendo usuario por ID"
	LogKeycloakUserGetError        = "Error obteniendo usuario de Keycloak"
	LogKeycloakUserNotFound        = "Usuario no encontrado en Keycloak"
	LogKeycloakUserDelete          = "Eliminando usuario de Keycloak"
	LogKeycloakUserDeleteOK        = "Usuario eliminado de Keycloak"
	LogKeycloakUserDeleteError     = "Error eliminando usuario de Keycloak"
	LogKeycloakPasswordSet         = "Configurando password para usuario"
	LogKeycloakPasswordSetOK       = "Password configurado exitosamente"
	LogKeycloakPasswordSetError    = "Error configurando password"
	LogKeycloakRoleGet             = "Obteniendo rol"
	LogKeycloakRoleGetError        = "Error obteniendo rol"
	LogKeycloakRoleAssign          = "Asignando rol a usuario"
	LogKeycloakRoleAssignOK        = "Rol asignado exitosamente"
	LogKeycloakRoleAssignError     = "Error asignando rol a usuario"
	LogKeycloakUserTokenRefresh    = "Refrescando token de usuario"
	LogKeycloakUserTokenRefreshOK  = "Token de usuario refrescado exitosamente"
	LogKeycloakUserTokenRefreshErr = "Error refrescando token de usuario"
	LogKeycloakEmailEmpty          = "Email no puede estar vacío"
	LogKeycloakUserIDEmpty         = "ID de usuario no puede estar vacío"
	LogKeycloakPasswordEmpty       = "Password no puede estar vacío"
	LogKeycloakRoleNameEmpty       = "Nombre de rol no puede estar vacío"
	LogKeycloakRefreshTokenEmpty   = "Refresh token no puede estar vacío"
	LogKeycloakUsernameEmpty       = "Nombre de usuario no puede estar vacío"
)

// ============================================
// REPOSITORY / MESSAGE REPOSITORY
// ============================================
const (
	LogRepoMsgInit           = "Inicializando repositorio de mensajes"
	LogRepoMsgInitOK         = "Repositorio de mensajes inicializado"
	LogRepoMsgInitError      = "Error inicializando repositorio de mensajes"
	LogEmployeeRepoInit      = "Inicializando repositorio de empleados"
	LogEmployeeRepoInitOK    = "Repositorio de empleados inicializado"
	LogEmployeeRepoInitError = "Error creando repositorio de empleados"
	LogIDEncoderInit         = "Inicializando ID encoder"
	LogIDEncoderInitError    = "Error inicializando ID encoder"
)

// ============================================
// HTTP REQUESTS
// ============================================
const (
	LogHTTPRequestIncoming = "Request HTTP entrante"
	LogHTTPResponseSent    = "Respuesta HTTP enviada"
	LogHTTPClientIP        = "IP del cliente"
	LogHTTPMethod          = "Método HTTP"
	LogHTTPPath            = "Path HTTP"
	LogHTTPStatus          = "Status HTTP"
)

// ============================================
// EMPLOYEE / USER SERVICES
// ============================================
const (
	LogEmployeeCreating              = "Creando empleado en base de datos"
	LogEmployeeCreated               = "Empleado creado exitosamente"
	LogEmployeeCreateError           = "Error creando empleado"
	LogEmployeeSaving                = "Guardando empleado en base de datos"
	LogEmployeeSaveError             = "Error guardando empleado en base de datos"
	LogEmployeeUpdating              = "Actualizando empleado"
	LogEmployeeUpdated               = "Empleado actualizado"
	LogEmployeeUpdateError           = "Error actualizando empleado"
	LogEmployeeUpdateKeycloakID      = "Actualizando ID de Keycloak del empleado"
	LogEmployeeUpdateKeycloakIDError = "Error actualizando ID de Keycloak en base de datos"
	LogEmployeeSearching             = "Buscando empleado"
	LogEmployeeGetByEmail            = "Obteniendo empleado por email"
	LogEmployeeGetByEmailError       = "Error obteniendo empleado por email"
	LogEmployeeGetByID               = "Obteniendo empleado por ID"
	LogEmployeeGetByIDError          = "Error obteniendo empleado por ID"
	LogEmployeeFound                 = "Empleado encontrado"
	LogEmployeeNotFound              = "Empleado no encontrado"
	LogEmployeeSearchError           = "Error buscando empleado"
	LogEmployeeRegistering           = "Registrando empleado"
	LogEmployeeRegisterError         = "Error registrando empleado"
	LogEmployeeRegisterSuccess       = "Empleado registrado exitosamente"
	LogEmployeeExists                = "El empleado ya existe"
	LogEmployeeDeleting              = "Eliminando empleado de base de datos"
	LogEmployeeDeleteError           = "Error eliminando empleado de base de datos"
	LogEmployeeLocating              = "Localizando empleado"
	LogEmployeeLocateError           = "Error localizando empleado"
)

// ============================================
// DEPENDENCY INJECTION
// ============================================
const (
	LogDepInit          = "Inicializando dependencias"
	LogDepInitComplete  = "Dependencias inicializadas exitosamente"
	LogDepInitError     = "Error inicializando dependencias"
	LogDepWiringService = "Inyectando servicio"
	LogDepConfigError   = "Error cargando configuración de dependencias"
)

// ============================================
// SECURITY / ENCODING
// ============================================
const (
	LogIDEncode               = "Codificando ID"
	LogIDEncodeOK             = "ID codificado exitosamente"
	LogIDEncodeError          = "Error codificando ID"
	LogIDDecode               = "Decodificando ID"
	LogIDDecodeOK             = "ID decodificado exitosamente"
	LogIDDecodeError          = "Error decodificando ID"
	LogIDEncoderInvalidUUID   = "UUID inválido"
	LogIDEncoderHashidsCreate = "Error creando hashids"
	LogIDEncoderEncodingError = "Error encodeando UUID"
	LogIDEncoderEmptyID       = "ID ofuscado no puede estar vacío"
	LogIDEncoderDecodingError = "Error decodeando ID ofuscado"
	LogIDEncoderInvalidFormat = "ID ofuscado tiene formato incorrecto"
	LogIDEncoderUUIDError     = "Error reconstruyendo UUID"
	LogIDEncoderMinLengthWarn = "MinLength es igual a 36, lo cual es el valor por defecto"
)

// ============================================
// UTILS / HELPERS
// ============================================
const (
	LogUtilsModuleRootSearch = "Buscando raíz del módulo"
	LogUtilsModuleRootFound  = "Raíz del módulo encontrada"
	LogUtilsModuleRootError  = "No se pudo encontrar la raíz del módulo"
	LogUtilsCurrentDirError  = "No se pudo determinar el directorio actual"
	LogUtilsPathResolved     = "Ruta resuelta exitosamente"
	LogUtilsPathError        = "Error resolviendo ruta"
)

// ============================================
// MIDDLEWARE / VALIDATORS
// ============================================
const (
	LogMiddlewareValidationStart    = "Iniciando validación de request body"
	LogMiddlewareValidationOK       = "Validación exitosa"
	LogMiddlewareValidationFailed   = "Validación de request fallida"
	LogMiddlewareBodyReadError      = "Error leyendo body del request"
	LogMiddlewareJSONParseError     = "Error parseando JSON del body"
	LogMiddlewareSchemaError        = "Error de validación de schema"
	LogMiddlewareResponseCacheError = "Error obteniendo mensaje de cache"
	LogMiddlewareResponseSuccess    = "Respuesta enviada exitosamente"
	LogMiddlewareErrorHandling      = "Manejando error de middleware"
)

// ============================================
// SCHEMA VALIDATION
// ============================================
const (
	LogSchemaValidatorInit      = "Inicializando validador de esquemas"
	LogSchemaValidatorInitOK    = "Validador de esquemas inicializado"
	LogSchemaValidatorInitError = "Error inicializando validador de esquemas"
	LogSchemaReading            = "Leyendo esquema JSON"
	LogSchemaReadError          = "Error leyendo esquema JSON"
	LogSchemaEmpty              = "Esquema JSON vacío o nulo"
	LogSchemaCompiling          = "Compilando esquema JSON"
	LogSchemaCompileError       = "Error compilando esquema JSON"
	LogSchemaValidating         = "Validando datos contra esquema"
	LogSchemaValidationOK       = "Validación de esquema exitosa"
	LogSchemaValidationError    = "Error en validación de esquema"
	LogSchemaFileNotFound       = "Archivo de esquema no encontrado"
	LogSchemaNotInitialized     = "Validador no inicializado"
)

// ============================================
// CONFIG LOADING
// ============================================
const (
	LogConfigFileNotFound    = "Archivo de configuración no encontrado"
	LogConfigFallback        = "Usando configuración por defecto"
	LogConfigReading         = "Leyendo archivo de configuración"
	LogConfigReadError       = "Error leyendo archivo de configuración"
	LogConfigParsing         = "Parseando configuración JSON"
	LogConfigParseError      = "Error parseando configuración JSON"
	LogConfigEnvOverride     = "Sobreescribiendo configuración con variables de entorno"
	LogConfigValidating      = "Validando configuración"
	LogConfigValidationError = "Error validando configuración"
)

// ============================================
// MESSAGES MODULE / SYSTEM MESSAGES
// ============================================
const (
	LogMessageServiceInit       = "Inicializando servicio de mensajes"
	LogMessageServiceInitOK     = "Servicio de mensajes inicializado"
	LogMessageCreate            = "Creando mensaje del sistema"
	LogMessageCreateOK          = "Mensaje creado exitosamente"
	LogMessageCreateError       = "Error creando mensaje"
	LogMessageCreateProcessing  = "Procesando creación de mensaje"
	LogMessageUpdate            = "Actualizando mensaje del sistema"
	LogMessageUpdateOK          = "Mensaje actualizado exitosamente"
	LogMessageUpdateError       = "Error actualizando mensaje"
	LogMessageUpdateProcessing  = "Procesando actualización de mensaje"
	LogMessageDelete            = "Eliminando mensaje del sistema"
	LogMessageDeleteOK          = "Mensaje eliminado exitosamente"
	LogMessageDeleteError       = "Error eliminando mensaje"
	LogMessageDeleteProcessing  = "Procesando eliminación de mensaje"
	LogMessageGet               = "Obteniendo mensaje del sistema"
	LogMessageGetOK             = "Mensaje obtenido exitosamente"
	LogMessageGetError          = "Error obteniendo mensaje"
	LogMessageList              = "Listando mensajes del sistema"
	LogMessageListOK            = "Mensajes listados exitosamente"
	LogMessageListError         = "Error listando mensajes"
	LogMessageValidation        = "Validando mensaje"
	LogMessageValidationOK      = "Mensaje validado exitosamente"
	LogMessageValidationError   = "Error validando mensaje"
	LogMessageCodeDuplicate     = "Código de mensaje duplicado"
	LogMessageTxBegin           = "Iniciando transacción para mensaje"
	LogMessageTxBeginOK         = "Transacción iniciada para mensaje"
	LogMessageTxCommit          = "Confirmando transacción de mensaje"
	LogMessageTxCommitOK        = "Transacción confirmada exitosamente"
	LogMessageTxCommitError     = "Error confirmando transacción"
	LogMessageTxRollback        = "Ejecutando rollback de transacción"
	LogMessageTxRollbackOK      = "Rollback ejecutado exitosamente"
	LogMessageTxRollbackError   = "Error ejecutando rollback"
	LogMessageCacheRefresh      = "Refrescando cache de mensajes"
	LogMessageCacheRefreshOK    = "Cache de mensajes refrescado"
	LogMessageCacheRefreshError = "Error refrescando cache de mensajes"
	LogMessageInvalidID         = "ID inválido"
	LogMessageIDEncodeError     = "Error ofuscando ID"
	LogMessageIDDecodeError     = "Error decodificando ID"
)
