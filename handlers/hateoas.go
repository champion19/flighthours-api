package handlers

type Link struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

type HateoasResource struct {
	Links []Link `json:"_links"`
}

func BuildAccountLinks(baseURL string, accountID string) []Link {
	return []Link{
		{
			Href:   baseURL + "/flighthours/api/v1/acounts/" + accountID,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   baseURL + "/flighthours/api/v1/acounts/" + accountID,
			Rel:    "update",
			Method: "PUT",
		},
		{
			Href:   baseURL + "/flighthours/api/v1/acounts/" + accountID,
			Rel:    "delete",
			Method: "DELETE",
		},
		{
			Href:   baseURL + "/flighthours/api/v1/acounts/",
			Rel:    "collection",
			Method: "GET",
		},
	}

}
