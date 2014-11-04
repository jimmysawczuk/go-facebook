package facebook

// A specific error that's returned from Facebook if there's an error with a request to the Graph API.
type GraphError struct {
	Code       int    `json:"code"`
	Subcode    int    `json:"error_subcode"`
	Message    string `json:"message"`
	Type       string `json:"type"`
	HTTPStatus int    `json:"http"`
}

type rawGraphError struct {
	Error GraphError `json:"error"`
}
