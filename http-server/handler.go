package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/heaven-chp/common-library-go/json"
	"github.com/heaven-chp/common-library-go/log"
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
func testGet(w http.ResponseWriter, r *http.Request) {
	log.Debug("[%s] [%s] start", r.RequestURI, r.Method)
	defer log.Debug("[%s] [%s] end", r.RequestURI, r.Method)

	log.Debug("header-1 : (%s)", r.Header.Get("header-1"))

	log.Debug("id : (%s)", mux.Vars(r)["id"])

	log.Debug("param-1 : (%s)", r.URL.Query().Get("param-1"))
	log.Debug("param-2 : (%s)", r.URL.Query().Get("param-2"))
	log.Debug("param-3 : (%s)", r.URL.Query().Get("param-3"))

	body, err := json.ToString(Test{ID: mux.Vars(r)["id"], Field1: 1, Field2: "value-2"})
	if err != nil {
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
func testPost(w http.ResponseWriter, r *http.Request) {
	body, err := json.ToString(ResponseSuccess{Field1: "value-1"})
	if err != nil {
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
func testDelete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
