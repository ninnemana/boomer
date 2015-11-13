package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"net/http"
)

var (
	rndr = render.New(render.Options{
		Directory:     "dist/public",     // Specify what path to load the templates from.
		Extensions:    []string{".html"}, // Specify extensions to load for templates.
		IsDevelopment: true,              // Render will now recompile the templates on every HTML response.
	})
)

func Options(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	rw.Write([]byte("cors is allowed"))
}

func UserInterface(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	rndr.HTML(rw, http.StatusOK, "index", nil)
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
