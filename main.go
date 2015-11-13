package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/ninnemana/boomer/api"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	router := httprouter.New()

	router.OPTIONS("/", api.Options)
	router.GET("/", api.UserInterface)
	router.POST("/bench", api.Boomer)

	var srv http.Server
	srv.Addr = ":8081"
	srv.Handler = cors.Default().Handler(router)

	var header string
	file, err := os.Open("header.txt")
	if err != nil {
		msg := ""
		if err != nil {
			msg = err.Error()
		}
		log.Printf("Failed to open the header file %s\n", msg)
	} else if file != nil {
		data, err := ioutil.ReadAll(file)
		if err == nil {
			header = string(data)
		}
	}

	http2.ConfigureServer(&srv, &http2.Server{})

	log.Printf("%s\n", header)
	log.Printf("Starting server on %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
