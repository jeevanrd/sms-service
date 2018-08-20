package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/golang/glog"
	"github.com/jeevanrd/sms-service/routers"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	routers.CreateAppRouter(router)
	glog.Fatal(http.ListenAndServe(":7070", router))
}
