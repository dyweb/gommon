package requests

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// FIXME: remove response, this is no longer needed I think ...
type Response struct {
	Res  *http.Response
	Text []byte
}

func (res *Response) JSON(data interface{}) error {
	// TODO: may change to decoder
	if err := json.Unmarshal(res.Text, &data); err != nil {
		return errors.Wrap(err, "error unmarshal json using map[string]string")
	}
	return nil
}

func (res *Response) JSONStringMap() (map[string]string, error) {
	var data map[string]string
	if err := json.Unmarshal(res.Text, &data); err != nil {
		return data, errors.Wrap(err, "error unmarshal json using map[string]string")
	}
	return data, nil
}
