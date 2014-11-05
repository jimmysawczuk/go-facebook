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

func (e GraphError) String() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

func (e GraphError) Error() string {
	return e.String()
}

func (e *GraphError) UnmarshalJSON(in []byte) (err error) {
	raw_error := struct {
		Error struct {
			Code       int    `json:"code"`
			Subcode    int    `json:"error_subcode"`
			Message    string `json:"message"`
			Type       string `json:"type"`
			HTTPStatus int    `json:"http"`
		} `json:"error"`
	}{}

	err = json.Unmarshal(in, &raw_error)

	if err == nil {
		*e = raw_error.Error
	}

	return err
}
