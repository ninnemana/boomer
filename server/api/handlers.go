package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

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

}
