package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/jeevanrd/sms-service/routers"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	routers.CreateAppRouter(router)
	http.ListenAndServe(":7070", router)
}
