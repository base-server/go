package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/heaven-chp/base-server-go/http-server/log"
	"github.com/heaven-chp/common-library-go/json"
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

// @Summary get test
// @Description get test
// @Accept json
// @Produce json
// @Param header-1 header string true "header-1 description" default(value-1)
// @Param param_1 query string true "param-1 description" Enums(1, 2, 3)
// @Param param_2 query string true "param-2 description" Enums(A, B, C, D) default(A)
// @Param param_3 query string true "param-3 description" default(AAA)
// @Param id path string true "id" default(id_1)
// @Success 200 {object} Test
// @Failure default {object} ResponseFailure
// @Router /v1/test/{id} [get]
// @tags test
func Get(w http.ResponseWriter, r *http.Request) {
	log.Server.Debug("handler start", "uri", r.RequestURI, "method", r.Method)
	defer log.Server.Debug("handler end", "uri", r.RequestURI, "method", r.Method)

	log.Server.Debug("header", "header-1", r.Header.Get("header-1"))

	log.Server.Debug("path", "id", mux.Vars(r)["id"])

	log.Server.Debug("parameter", "param-1", r.URL.Query().Get("param-1"), "param-2", r.URL.Query().Get("param-2"), "param-3", r.URL.Query().Get("param-3"))

	if body, err := json.ToString(Test{ID: mux.Vars(r)["id"], Field1: 1, Field2: "value-2"}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"Cause":"` + err.Error() + `"}`))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))
	}
}

// @Summary post test
// @Description post test
// @Accept json
// @Produce json
// @Param request body Test true "country selection"
// @Success 200 {object} ResponseSuccess
// @Failure default {object} ResponseFailure
// @Router /v1/test [post]
// @tags test
func Post(w http.ResponseWriter, r *http.Request) {
	if body, err := json.ToString(ResponseSuccess{Field1: "value-1"}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"Cause":"` + err.Error() + `"}`))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))
	}
}

// @Summary delete test
// @Description delete test
// @Accept json
// @Produce json
// @Param header-1 header string true "header-1 description" default(value-1)
// @Param id path string true "id" default(id_1)
// @Success 204
// @Failure default {object} ResponseFailure
// @Router /v1/test/{id} [delete]
// @tags test
func Delete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
