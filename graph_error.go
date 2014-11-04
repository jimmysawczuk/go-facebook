package facebook

import (
	"encoding/json"
	"fmt"
)

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

func NewGraphError(data map[string]interface{}) (fge GraphError, err error) {
	b, _ := json.Marshal(data)
	err = json.Unmarshal(b, &fge)

	return fge, err
}

func (e GraphError) String() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

func (e GraphError) Error() string {
	return e.String()
}

func parseError(status_code int, buf []byte) (GraphError, error) {
	raw_error := rawGraphError{}
	err := json.Unmarshal(buf, &raw_error)

	if err == nil {
		return raw_error.Error, nil
	} else {
		return GraphError{}, err
	}
}
