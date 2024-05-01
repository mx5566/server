package account

import "github.com/mx5566/server/server/model"

type AccountInfo struct {
	AccountID     int64  `bson:"accountId"`
	AccountName   string `bson:"accountName"`
	AccountPasswd string `bson:"accountPasswd"`
}

type Account struct {
	accountInfo   AccountInfo
	roleInfos     []*model.PlayerSimpleInfo
	PlayerID      int64
	GateClusterID uint32
}

func (a Account) PlayerLogin(playerID int64) bool {
	for _, v := range a.roleInfos {
		if v == nil {
			continue
		}

		if v.PlayerID == playerID {
			a.PlayerID = playerID
			return true
		}
	}

	return false
}
