package database

import "gopkg.in/pg.v4"

type Repo struct {
	db         *pg.DB
}

type Repository interface {
	GetAccount(username string, authId string) (Account, error)
	GetPhoneNumber(accountId int, number string) (PhoneNumber, error)
}

func (r *Repo)  GetAccount(username string, authId string) (Account, error) {
	var account Account
	err := r.db.Model(&account).Where("auth_id = ? AND username = ?", authId, username).First()
	return account, err
}

func (r *Repo)  GetPhoneNumber(accountId int, number string) (PhoneNumber, error) {
	var phone PhoneNumber
	err :=  r.db.Model(&phone).Where("account_id = ? AND number = ?", accountId, number).First()
	return phone, err
}

func NewDatabaseRepo(dbUser,password,dbName,dbhost string) (*Repo, error) {
	Repo := &Repo{
		db:   GetDbConnectionWithOptions(dbUser,password,dbName,dbhost),
	}
	return Repo, nil
}