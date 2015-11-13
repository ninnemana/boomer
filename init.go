package boom

import (
	"api"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"

	"net/http"
)

func init() {
	router := httprouter.New()

	router.OPTIONS("/", api.Options)
	router.GET("/", api.UserInterface)
	router.POST("/bench", api.Boomer)

	http.Handle("/", cors.Default().Handler(router))
}
