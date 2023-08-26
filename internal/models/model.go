package models

type RefreshToken struct {
	RefreshTokenBCrypt string `bson:"refresh_token_bcrypt"`
	Exp                int64  `bson:"exp"`
}

type User struct {
	GUID         string       `bson:"_id"`
	RefreshToken RefreshToken `bson:"refresh_token"`
}
