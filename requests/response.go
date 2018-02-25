package requests

import (
	"encoding/json"
	"net/http"

	"github.com/dyweb/gommon/errors"
)

type Response struct {
	Res        *http.Response
	StatusCode int
	Data       []byte
}

func (res *Response) JSON(data interface{}) error {
	if err := json.Unmarshal(res.Data, &data); err != nil {
		return errors.Wrap(err, "error unmarshal json using map[string]string")
	}
	return nil
}

func (res *Response) JSONStringMap() (map[string]string, error) {
	var data map[string]string
	if err := json.Unmarshal(res.Data, &data); err != nil {
		return data, errors.Wrap(err, "error unmarshal json using map[string]string")
	}
	return data, nil
}
