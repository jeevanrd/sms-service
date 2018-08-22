package database

import "gopkg.in/pg.v4"

func GetDbConnection(dbUser string, dbName string) *pg.DB {
	db := pg.Connect(&pg.Options{
		User: dbUser,
		Password: dbUser,
		Database: dbName,
	})
	return db
}


func GetDbConnectionWithOptions(dbUser string, password string, dbName string, hostName string) *pg.DB {
	db := pg.Connect(&pg.Options{
		User:dbUser,
		Password:password,
		Database:dbName,
		Addr:hostName,
	})
	return db
}

