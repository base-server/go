package handler

import (
	"net/http"
	"strings"

	"github.com/base-server/go/common/log"
	"github.com/common-library/go/json"
)

type ResponseSuccess struct {
	Field1 string `json:"field-1"`
}

type ResponseFailure struct {
	Cause string `json:"cause"`
}

type Test struct {
	ID     string `json:"id" example:"id-1"`
	Field1 int    `json:"field-1" example:"1"`
	Field2 string `json:"field-2" example:"value-2"`
}

func Get(w http.ResponseWriter, r *http.Request) {
	log.Log.Debug("handler start", "uri", r.RequestURI, "method", r.Method)
	defer log.Log.Debug("handler end", "uri", r.RequestURI, "method", r.Method)

	log.Log.Debug("header", "header-1", r.Header.Get("header-1"))

	log.Log.Debug("path", "id", strings.Split(r.URL.Path, "/")[3])

	log.Log.Debug("parameter", "param-1", r.URL.Query().Get("param-1"), "param-2", r.URL.Query().Get("param-2"), "param-3", r.URL.Query().Get("param-3"))

	if body, err := json.ToString(Test{ID: strings.Split(r.URL.Path, "/")[3], Field1: 1, Field2: "value-2"}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"Cause":"` + err.Error() + `"}`))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))
	}
}

func Post(w http.ResponseWriter, r *http.Request) {
	if body, err := json.ToString(ResponseSuccess{Field1: "value-1"}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"Cause":"` + err.Error() + `"}`))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
