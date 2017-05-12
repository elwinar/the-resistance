package main

import (
	"encoding/json"
	"net/http"
)

func write(w http.ResponseWriter, p interface{}) error {
	r, err := json.Marshal(p)
	if err != nil {
		return err
	}

	w.Write(r)
	return nil
}
