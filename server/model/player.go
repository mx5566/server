package model

type PlayerSimpleInfo struct {
	PlayerID  int64  `bson:"playerID"`
	Name      string `bson:"playerName"`
	Level     int32  `bson:"level"`
	Gold      int64  `bson:"gold"`
	AccountID int64  `bson:"accountID"`
}
