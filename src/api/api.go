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

func WriteError(w http.ResponseWriter, status int, err error) error {
	w.WriteHeader(status)
	return Write(w, map[string]interface{}{
		"error": err.Error(),
	})
}

func Read(r *http.Request, p interface{}) error {
	return json.NewDecoder(r.Body).Decode(p)
}
