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
	"strings"
)

const DB_USER = "postgres"
const DB_DATABASE = "postgres"
const DB_PASSWORD = "postgres"
const DB_HOST = "localhost:5432"



func CreateAppRouter(router *mux.Router) *mux.Router {
	var Logger log.Logger
	Logger = log.NewJSONLogger(os.Stdout)
	Logger = &serializedLogger{Logger: Logger}
	Logger = log.With(Logger, "ts", log.DefaultTimestampUTC)

	var repo *database.Repo
	var err error

	if(os.Getenv("DATABASE_URL") != "") {
		repo,err = GetDBConnectionCredentialsFromUrl(os.Getenv("DATABASE_URL"))
	} else {
		repo,err = database.NewDatabaseRepo(DB_USER, DB_PASSWORD, DB_DATABASE, DB_HOST)
	}

	if(err !=nil) {
		glog.Fatal("Unable to connect to db")
	}
	redisClient := utils.GetCacheClient()
	is := smsservice.NewService(repo, utils.LocalCache{Client:redisClient})
	authFilter := auth.Auth{Repo:repo}

	router = smsservice.MakeHandler(is,Logger, router, authFilter.AuthHandler, authFilter.ContentTypeHandler)
	return router
}

//to connect to heroku db
func GetDBConnectionCredentialsFromUrl(url string) (*database.Repo, error) {
	stringArr := strings.Split(url, "//")
	dbInfo := strings.Split(stringArr[1], "@")
	credentials := strings.Split(dbInfo[0], ":")
	dbNameAndHost := strings.Split(dbInfo[1], "/")
	return database.NewDatabaseRepo(credentials[0], credentials[1], dbNameAndHost[1], dbNameAndHost[0])
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