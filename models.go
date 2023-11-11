package main

type User struct {
	IDContract          string              `bson:"contract_id"`
	Subscriber          string              `bson:"subscriber"`
	Document            string              `bson:"document"`
	FirstName           string              `bson:"first_name"`
	LastName            string              `bson:"last_name"`
	Address             string              `bson:"address"`
	Phone               string              `bson:"phone"`
	Status              string              `bson:"status"`
	Subscription        string              `bson:"subscription"`
	InformationMikrotik MikrotikInformation `bson:"mikrotik_information"`
}

type MikrotikInformation struct {
	Secret        string `bson:"secret"`
	CallerID      string `bson:"caller_id"`
	Comment       string `bson:"comment"`
	RemoteAddress string `bson:"remote_address"`
	Profile       string `bson:"profile"`
	Bts           string `bson:"bts"`
	Host          string `bson:"host"`
}
