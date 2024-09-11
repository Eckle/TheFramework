package httpcodec

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
)

type HttpCodec interface {
	Encode(w http.ResponseWriter, r *http.Request, data interface{}) error
	Decode(r *http.Request) (map[string]interface{}, error)
}

type JsonCodec struct{}

func (codec JsonCodec) Encode(w http.ResponseWriter, r *http.Request, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return err
	}

	return nil
}

func (codec JsonCodec) Decode(r *http.Request) (map[string]interface{}, error) {
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

type HtmxCodec struct{}

func (codec HtmxCodec) Encode(w http.ResponseWriter, r *http.Request, data interface{}) error {
	template, err := gonja.FromFile("./templates/homepage.html")
	if err != nil {
		return err
	}

	v, ok := data.(map[string]interface{})
	if !ok {
		return errors.New("data is not of type map[string]interface{}")
	}

	if err = template.Execute(w, exec.NewContext(v)); err != nil {
		return err
	}

	return nil
}

func (codec HtmxCodec) Decode(r *http.Request) (map[string]interface{}, error) {
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
