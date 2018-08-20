package database

type Account struct {
	TableName struct{} `sql:"account"`
	Id     int
	AuthId   string
	Username string
}

type PhoneNumber struct {
	TableName struct{} `sql:"phone_number"`
	Id       int
	Number    string
	AccountId int
}
