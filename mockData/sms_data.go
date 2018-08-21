package mockData

import "github.com/jeevanrd/sms-service/database"


var Account1 = database.Account{
	Id:1,
	Username:"user1",
	AuthId:"p1",
}

var Account2 = database.Account{
	Id:2,
	Username:"user2",
	AuthId:"p2",
}

var Phone1 = database.PhoneNumber{
	Id:1,
	Number:"4924195509198",
	AccountId:1,
}

var Phone2 = database.PhoneNumber{
	Id:2,
	Number:"4924195509196",
	AccountId:2,
}

var Phone3 = database.PhoneNumber{
	Id:2,
	Number:"4924195509195",
	AccountId:2,
}