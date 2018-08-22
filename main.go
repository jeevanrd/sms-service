package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/jeevanrd/sms-service/routers"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if(port  == "") {
		port = "7070"
	}
	router := mux.NewRouter().StrictSlash(true)
	routers.CreateAppRouter(router)
	http.ListenAndServe(":" + port, router)
}
