package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"net/http"
)

var (
	rndr = render.New(render.Options{})
)

func Options(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	rw.Write([]byte("cors is allowed"))
}

func UserInterface(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	rw.Write([]byte("display the ui"))
}

func Boomer(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	var r Request
	err := json.NewDecoder(req.Body).Decode(&r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = r.validate()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = r.Do()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rndr.JSON(rw, http.StatusOK, r.Report)
}
