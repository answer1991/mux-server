package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func UnmarshalRequestBodyToJson(r *http.Request, body interface{}) (err error) {
	bytes, err := ioutil.ReadAll(r.Body)

	if nil != err {
		return err
	}

	err = json.Unmarshal(bytes, body)

	if nil != err {
		return err
	}

	return nil
}
