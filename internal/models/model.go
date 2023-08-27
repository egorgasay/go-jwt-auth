package models

type TokenData struct {
	GUID        string `bson:"guid"`
	RefreshHash string `bson:"refresh_hash"`
	RefreshExp  int64  `bson:"refresh_exp"`
	AccessExp   int64  `bson:"access_exp"`
}
