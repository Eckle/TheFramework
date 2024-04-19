package httpcodec

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func Encode(w http.ResponseWriter, r *http.Request, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return err
	}

	return nil
}

func Decode(r *http.Request) (map[string]interface{}, error) {
	content, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})

	err = json.Unmarshal(content, &data)
	if err != nil {
		return nil, errors.New("could not parse request body")
	}

	return data, nil
}
