package api

import (
	"encoding/json"
	"net/http"
)

func Write(w http.ResponseWriter, p interface{}) error {
	r, err := json.Marshal(p)
	if err != nil {
		return err
	}

	w.Write(r)
	return nil
}

func Read(r *http.Request, p interface{}) error {
	return json.NewDecoder(r.Body).Decode(p)
}
