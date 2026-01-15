package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Link struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

type HATEOASResource struct {
	Links []Link `json:"_links"`
}

// GetBaseURL extrae la URL base de la petición (scheme + host)
func GetBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + c.Request.Host
}

// SetLocationHeader establece el Location header con la URL del recurso
func SetLocationHeader(c *gin.Context, baseURL, resource, resourceID string) {
	locationURL := BuildResourceURL(baseURL, resource, resourceID)
	c.Header("Location", locationURL)
}

// BuildResourceURL construye una URL completa para un recurso
// Ejemplo: BuildResourceURL(baseURL, "accounts", encodedID) → "http://host/flighthours/api/v1/accounts/xyz"
func BuildResourceURL(baseURL, resource, resourceID string) string {
	return fmt.Sprintf("%s/flighthours/api/v1/%s/%s", baseURL, resource, resourceID)
}

// BuildCollectionURL construye una URL completa para una colección
// Ejemplo: BuildCollectionURL(baseURL, "accounts") → "http://host/flighthours/api/v1/accounts"
func BuildCollectionURL(baseURL, resource string) string {
	return fmt.Sprintf("%s/flighthours/api/v1/%s", baseURL, resource)
}

// BuildResourceLinks construye links HATEOAS genéricos para un recurso
// resource: nombre del recurso (ej: "accounts", "transactions")
// resourceID: ID ya ofuscado del recurso
func BuildResourceLinks(baseURL, resource, resourceID string) []Link {
	resourceURL := BuildResourceURL(baseURL, resource, resourceID)
	collectionURL := BuildCollectionURL(baseURL, resource)

	return []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   resourceURL,
			Rel:    "update",
			Method: "PUT",
		},
		{
			Href:   resourceURL,
			Rel:    "delete",
			Method: "DELETE",
		},
		{
			Href:   collectionURL,
			Rel:    "collection",
			Method: "GET",
		},
	}
}

// BuildAccountLinks construye links específicos para cuentas (wrapper para compatibilidad)
func BuildAccountLinks(baseURL string, accountID string) []Link {
	return BuildResourceLinks(baseURL, "accounts", accountID)
}

// BuildMessageLinks construye links HATEOAS para un mensaje específico
func BuildMessageLinks(baseURL string, messageID string) []Link {
	return BuildResourceLinks(baseURL, "messages", messageID)
}

// BuildMessageCreatedLinks construye links para un mensaje recién creado
func BuildMessageCreatedLinks(baseURL string, messageID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "messages", messageID)
	collectionURL := BuildCollectionURL(baseURL, "messages")

	return []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   resourceURL,
			Rel:    "update",
			Method: "PUT",
		},
		{
			Href:   resourceURL,
			Rel:    "delete",
			Method: "DELETE",
		},
		{
			Href:   collectionURL,
			Rel:    "list",
			Method: "GET",
		},
	}
}

// BuildMessageUpdatedLinks construye links para un mensaje actualizado
func BuildMessageUpdatedLinks(baseURL string, messageID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "messages", messageID)
	collectionURL := BuildCollectionURL(baseURL, "messages")

	return []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   resourceURL,
			Rel:    "delete",
			Method: "DELETE",
		},
		{
			Href:   collectionURL,
			Rel:    "list",
			Method: "GET",
		},
	}
}

// BuildMessageListLinks construye links para la lista de mensajes
func BuildMessageListLinks(baseURL string) []Link {
	collectionURL := BuildCollectionURL(baseURL, "messages")

	return []Link{
		{
			Href:   collectionURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   collectionURL,
			Rel:    "create",
			Method: "POST",
		},
	}
}

// ============================================================================
// AIRLINE HATEOAS LINKS
// ============================================================================

// BuildAirlineLinks construye links HATEOAS para una aerolínea específica
func BuildAirlineLinks(baseURL string, airlineID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "airlines", airlineID)
	collectionURL := BuildCollectionURL(baseURL, "airlines")

	return []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   resourceURL + "/activate",
			Rel:    "activate",
			Method: "PATCH",
		},
		{
			Href:   resourceURL + "/deactivate",
			Rel:    "deactivate",
			Method: "PATCH",
		},
		{
			Href:   collectionURL,
			Rel:    "collection",
			Method: "GET",
		},
	}
}

// BuildAirlineListLinks construye links para la lista de aerolíneas
func BuildAirlineListLinks(baseURL string) []Link {
	collectionURL := BuildCollectionURL(baseURL, "airlines")

	return []Link{
		{
			Href:   collectionURL,
			Rel:    "self",
			Method: "GET",
		},
	}
}

// BuildAirlineStatusLinks construye links para respuesta de cambio de status
func BuildAirlineStatusLinks(baseURL string, airlineID string, isActive bool) []Link {
	resourceURL := BuildResourceURL(baseURL, "airlines", airlineID)
	collectionURL := BuildCollectionURL(baseURL, "airlines")

	links := []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
	}

	// Si está activo, mostrar link para desactivar y viceversa
	if isActive {
		links = append(links, Link{
			Href:   resourceURL + "/deactivate",
			Rel:    "deactivate",
			Method: "PATCH",
		})
	} else {
		links = append(links, Link{
			Href:   resourceURL + "/activate",
			Rel:    "activate",
			Method: "PATCH",
		})
	}

	links = append(links, Link{
		Href:   collectionURL,
		Rel:    "collection",
		Method: "GET",
	})

	return links
}

// ============================================================================
// AIRPORT HATEOAS LINKS
// ============================================================================

// BuildAirportLinks construye links HATEOAS para un aeropuerto específico
func BuildAirportLinks(baseURL string, airportID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "airports", airportID)
	collectionURL := BuildCollectionURL(baseURL, "airports")

	return []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   resourceURL + "/activate",
			Rel:    "activate",
			Method: "PATCH",
		},
		{
			Href:   resourceURL + "/deactivate",
			Rel:    "deactivate",
			Method: "PATCH",
		},
		{
			Href:   collectionURL,
			Rel:    "collection",
			Method: "GET",
		},
	}
}

// BuildAirportListLinks construye links para la lista de aeropuertos
func BuildAirportListLinks(baseURL string) []Link {
	collectionURL := BuildCollectionURL(baseURL, "airports")

	return []Link{
		{
			Href:   collectionURL,
			Rel:    "self",
			Method: "GET",
		},
	}
}

// BuildAirportStatusLinks construye links para respuesta de cambio de status
func BuildAirportStatusLinks(baseURL string, airportID string, isActive bool) []Link {
	resourceURL := BuildResourceURL(baseURL, "airports", airportID)
	collectionURL := BuildCollectionURL(baseURL, "airports")

	links := []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
	}

	// Si está activo, mostrar link para desactivar y viceversa
	if isActive {
		links = append(links, Link{
			Href:   resourceURL + "/deactivate",
			Rel:    "deactivate",
			Method: "PATCH",
		})
	} else {
		links = append(links, Link{
			Href:   resourceURL + "/activate",
			Rel:    "activate",
			Method: "PATCH",
		})
	}

	links = append(links, Link{
		Href:   collectionURL,
		Rel:    "collection",
		Method: "GET",
	})

	return links
}

// ============================================================================
// EMPLOYEE HATEOAS LINKS
// ============================================================================

// BuildEmployeeLinks construye links HATEOAS para un empleado específico
func BuildEmployeeLinks(baseURL string, employeeID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "employees", employeeID)
	collectionURL := BuildCollectionURL(baseURL, "employees")

	return []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   resourceURL,
			Rel:    "update",
			Method: "PUT",
		},
		{
			Href:   resourceURL,
			Rel:    "delete",
			Method: "DELETE",
		},
		{
			Href:   collectionURL,
			Rel:    "collection",
			Method: "GET",
		},
	}
}

// BuildEmployeeMeLinks construye links HATEOAS para el endpoint /employees/me
func BuildEmployeeMeLinks(baseURL string) []Link {
	meURL := baseURL + "/flighthours/api/v1/employees/me"

	return []Link{
		{
			Href:   meURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   meURL,
			Rel:    "update",
			Method: "PUT",
		},
		{
			Href:   meURL,
			Rel:    "delete",
			Method: "DELETE",
		},
	}
}

// ============================================================================
// DAILY LOGBOOK HATEOAS LINKS
// ============================================================================

// BuildDailyLogbookLinks construye links HATEOAS para una bitácora diaria específica
func BuildDailyLogbookLinks(baseURL string, logbookID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "daily-logbooks", logbookID)
	collectionURL := BuildCollectionURL(baseURL, "daily-logbooks")

	return []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   resourceURL,
			Rel:    "update",
			Method: "PUT",
		},
		{
			Href:   resourceURL,
			Rel:    "delete",
			Method: "DELETE",
		},
		{
			Href:   resourceURL + "/activate",
			Rel:    "activate",
			Method: "PATCH",
		},
		{
			Href:   resourceURL + "/deactivate",
			Rel:    "deactivate",
			Method: "PATCH",
		},
		{
			Href:   collectionURL,
			Rel:    "collection",
			Method: "GET",
		},
	}
}

// BuildDailyLogbookListLinks construye links para la lista de bitácoras diarias
func BuildDailyLogbookListLinks(baseURL string) []Link {
	collectionURL := BuildCollectionURL(baseURL, "daily-logbooks")

	return []Link{
		{
			Href:   collectionURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   collectionURL,
			Rel:    "create",
			Method: "POST",
		},
	}
}

// BuildDailyLogbookStatusLinks construye links para respuesta de cambio de status
func BuildDailyLogbookStatusLinks(baseURL string, logbookID string, isActive bool) []Link {
	resourceURL := BuildResourceURL(baseURL, "daily-logbooks", logbookID)
	collectionURL := BuildCollectionURL(baseURL, "daily-logbooks")

	links := []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
	}

	// Si está activo, mostrar link para desactivar y viceversa
	if isActive {
		links = append(links, Link{
			Href:   resourceURL + "/deactivate",
			Rel:    "deactivate",
			Method: "PATCH",
		})
	} else {
		links = append(links, Link{
			Href:   resourceURL + "/activate",
			Rel:    "activate",
			Method: "PATCH",
		})
	}

	links = append(links, Link{
		Href:   collectionURL,
		Rel:    "collection",
		Method: "GET",
	})

	return links
}

// BuildDailyLogbookCreatedLinks construye links para una bitácora recién creada
func BuildDailyLogbookCreatedLinks(baseURL string, logbookID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "daily-logbooks", logbookID)
	collectionURL := BuildCollectionURL(baseURL, "daily-logbooks")

	return []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   resourceURL,
			Rel:    "update",
			Method: "PUT",
		},
		{
			Href:   resourceURL,
			Rel:    "delete",
			Method: "DELETE",
		},
		{
			Href:   collectionURL,
			Rel:    "list",
			Method: "GET",
		},
	}
}

// BuildDailyLogbookDeletedLinks construye links para respuesta de eliminación
func BuildDailyLogbookDeletedLinks(baseURL string) []Link {
	collectionURL := BuildCollectionURL(baseURL, "daily-logbooks")

	return []Link{
		{
			Href:   collectionURL,
			Rel:    "list",
			Method: "GET",
		},
		{
			Href:   collectionURL,
			Rel:    "create",
			Method: "POST",
		},
	}
}

// ============================================================================
// AIRCRAFT REGISTRATION HATEOAS LINKS
// ============================================================================

// BuildAircraftRegistrationLinks construye links HATEOAS para una matrícula específica
func BuildAircraftRegistrationLinks(baseURL string, registrationID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "aircraft-registrations", registrationID)
	collectionURL := BuildCollectionURL(baseURL, "aircraft-registrations")

	return []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   resourceURL,
			Rel:    "update",
			Method: "PUT",
		},
		{
			Href:   collectionURL,
			Rel:    "collection",
			Method: "GET",
		},
	}
}

// BuildAircraftRegistrationListLinks construye links para la lista de matrículas
func BuildAircraftRegistrationListLinks(baseURL string) []Link {
	collectionURL := BuildCollectionURL(baseURL, "aircraft-registrations")

	return []Link{
		{
			Href:   collectionURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   collectionURL,
			Rel:    "create",
			Method: "POST",
		},
	}
}

// BuildAircraftRegistrationCreatedLinks construye links para una matrícula recién creada
func BuildAircraftRegistrationCreatedLinks(baseURL string, registrationID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "aircraft-registrations", registrationID)
	collectionURL := BuildCollectionURL(baseURL, "aircraft-registrations")

	return []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   resourceURL,
			Rel:    "update",
			Method: "PUT",
		},
		{
			Href:   collectionURL,
			Rel:    "list",
			Method: "GET",
		},
	}
}

// ============================================================================
// AIRCRAFT MODEL HATEOAS LINKS
// ============================================================================

// BuildAircraftModelLinks construye links HATEOAS para un modelo de aeronave específico
func BuildAircraftModelLinks(baseURL string, modelID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "aircraft-models", modelID)
	collectionURL := BuildCollectionURL(baseURL, "aircraft-models")

	return []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   collectionURL,
			Rel:    "collection",
			Method: "GET",
		},
	}
}

// BuildAircraftModelListLinks construye links para la lista de modelos de aeronave
func BuildAircraftModelListLinks(baseURL string) []Link {
	collectionURL := BuildCollectionURL(baseURL, "aircraft-models")

	return []Link{
		{
			Href:   collectionURL,
			Rel:    "self",
			Method: "GET",
		},
	}
}

// ============================================================================
// ROUTE HATEOAS LINKS
// ============================================================================

// BuildRouteLinks construye links HATEOAS para una ruta específica
func BuildRouteLinks(baseURL string, routeID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "routes", routeID)
	collectionURL := BuildCollectionURL(baseURL, "routes")

	return []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   collectionURL,
			Rel:    "collection",
			Method: "GET",
		},
	}
}

// BuildRouteListLinks construye links para la lista de rutas
func BuildRouteListLinks(baseURL string) []Link {
	collectionURL := BuildCollectionURL(baseURL, "routes")

	return []Link{
		{
			Href:   collectionURL,
			Rel:    "self",
			Method: "GET",
		},
	}
}

// ============================================================================
// AIRLINE ROUTE HATEOAS LINKS
// ============================================================================

// BuildAirlineRouteLinks construye links HATEOAS para una ruta aerolínea específica
func BuildAirlineRouteLinks(baseURL string, airlineRouteID string, isActive bool) []Link {
	resourceURL := BuildResourceURL(baseURL, "airline-routes", airlineRouteID)
	collectionURL := BuildCollectionURL(baseURL, "airline-routes")

	links := []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
	}

	// Si está activo, mostrar link para desactivar y viceversa
	if isActive {
		links = append(links, Link{
			Href:   resourceURL + "/deactivate",
			Rel:    "deactivate",
			Method: "PATCH",
		})
	} else {
		links = append(links, Link{
			Href:   resourceURL + "/activate",
			Rel:    "activate",
			Method: "PATCH",
		})
	}

	links = append(links, Link{
		Href:   collectionURL,
		Rel:    "collection",
		Method: "GET",
	})

	return links
}

// BuildAirlineRouteListLinks construye links para la lista de rutas aerolínea
func BuildAirlineRouteListLinks(baseURL string) []Link {
	collectionURL := BuildCollectionURL(baseURL, "airline-routes")

	return []Link{
		{
			Href:   collectionURL,
			Rel:    "self",
			Method: "GET",
		},
	}
}

// BuildAirlineRouteStatusLinks construye links para respuesta de cambio de status
func BuildAirlineRouteStatusLinks(baseURL string, airlineRouteID string, isActive bool) []Link {
	resourceURL := BuildResourceURL(baseURL, "airline-routes", airlineRouteID)
	collectionURL := BuildCollectionURL(baseURL, "airline-routes")

	links := []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
	}

	// Si está activo, mostrar link para desactivar y viceversa
	if isActive {
		links = append(links, Link{
			Href:   resourceURL + "/deactivate",
			Rel:    "deactivate",
			Method: "PATCH",
		})
	} else {
		links = append(links, Link{
			Href:   resourceURL + "/activate",
			Rel:    "activate",
			Method: "PATCH",
		})
	}

	links = append(links, Link{
		Href:   collectionURL,
		Rel:    "collection",
		Method: "GET",
	})

	return links
}

// ============================================================================
// DAILY LOGBOOK DETAIL HATEOAS LINKS
// ============================================================================

// BuildDailyLogbookDetailLinks construye links HATEOAS para un detalle de bitácora
// Retorna un mapa para usar directamente en el DTO response
func BuildDailyLogbookDetailLinks(c *gin.Context, detailID string) map[string]string {
	baseURL := GetBaseURL(c)
	resourceURL := BuildResourceURL(baseURL, "daily-logbook-details", detailID)

	return map[string]string{
		"self":   resourceURL,
		"update": resourceURL,
		"delete": resourceURL,
	}
}

// BuildDailyLogbookDetailLinksArray construye links HATEOAS como array
func BuildDailyLogbookDetailLinksArray(baseURL string, detailID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "daily-logbook-details", detailID)

	return []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   resourceURL,
			Rel:    "update",
			Method: "PUT",
		},
		{
			Href:   resourceURL,
			Rel:    "delete",
			Method: "DELETE",
		},
	}
}

// BuildDailyLogbookDetailListLinks construye links para la lista de detalles de una bitácora
func BuildDailyLogbookDetailListLinks(baseURL string, logbookID string) []Link {
	logbookURL := BuildResourceURL(baseURL, "daily-logbooks", logbookID)
	detailsURL := logbookURL + "/details"

	return []Link{
		{
			Href:   detailsURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   detailsURL,
			Rel:    "create",
			Method: "POST",
		},
		{
			Href:   logbookURL,
			Rel:    "logbook",
			Method: "GET",
		},
	}
}
