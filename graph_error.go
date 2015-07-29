package facebook

import (
	"encoding/json"
	"fmt"
)

// GraphError is a specific error that's returned from Facebook if there's an error with a request to the Graph API.
type GraphError struct {
	Code       int    `json:"code"`
	Subcode    int    `json:"error_subcode"`
	Message    string `json:"message"`
	Type       string `json:"type"`
	HTTPStatus int    `json:"http"`
}

// String returns a string representation of the error in the format: Error <code>: <message>.
func (e GraphError) String() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

func (e GraphError) Error() string {
	return e.String()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (e *GraphError) UnmarshalJSON(in []byte) (err error) {
	raw := struct {
		Error struct {
			Code       int    `json:"code"`
			Subcode    int    `json:"error_subcode"`
			Message    string `json:"message"`
			Type       string `json:"type"`
			HTTPStatus int    `json:"http"`
		} `json:"error"`
	}{}

	err = json.Unmarshal(in, &raw)

	if err == nil {
		*e = raw.Error
	}

	return err
}
