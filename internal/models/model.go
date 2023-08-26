package models

type Token struct {
	RefreshHash string `bson:"refresh_hash"`
	RefreshExp  int64  `bson:"refresh_exp"`
	AccessExp   int64  `bson:"access_exp"`
}

type User struct {
	GUID  string `bson:"guid"`
	Token Token  `bson:"refresh_token"`
}
