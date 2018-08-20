package routers

import (
	"github.com/jeevanrd/sms-service/database"
	"github.com/gorilla/mux"
	"github.com/jeevanrd/sms-service/utils"
	"github.com/jeevanrd/sms-service/smsservice"
	"sync"
	"os"
	"github.com/go-kit/kit/log"
	"github.com/golang/glog"
	"github.com/jeevanrd/sms-service/auth"
)

const DB_USER = "postgres"
const DB_DATABASE = "postgres"

func CreateAppRouter(router *mux.Router) *mux.Router {
	var Logger log.Logger
	Logger = log.NewJSONLogger(os.Stdout)
	Logger = &serializedLogger{Logger: Logger}
	Logger = log.With(Logger, "ts", log.DefaultTimestampUTC)

	repo,err := database.NewDatabaseRepo(DB_USER, DB_DATABASE)
	if(err !=nil) {
		glog.Fatal("Unable to connect to db")
	}
	redisClient := utils.GetCacheClient()
	is := smsservice.NewService(repo, utils.Cache{Client:redisClient})
	authFilter := auth.Auth{Repo:repo}

	router = smsservice.MakeHandler(is,Logger, router, authFilter.AuthHandler, authFilter.ContentTypeHandler)
	return router
}

type serializedLogger struct {
	mtx sync.Mutex
	log.Logger
}

func (l *serializedLogger) Log(keyvals ...interface{}) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	return l.Logger.Log(keyvals...)
}