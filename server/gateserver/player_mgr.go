package gateserver

type Player struct {
	PlayerID     int64  // 玩家的数据库ID
	ConnID       uint32 // 客户端连接的socketid
	GameServerID uint32 // 所在的游戏服务器ID
	AccountID    int64  // 账号ID数据库ID
}

// 管理所有登录的玩家
type PlayerMgr struct {
	players map[int64]*Player
}
