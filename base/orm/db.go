package orm

import (
	"context"
	"fmt"
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base/conf"
	"github.com/mx5566/server/base/orm/mongodb"
)

func OpenMongodb(db conf.DB) {
	ctx := context.Background()

	//连接数据库
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/?authSource=admin", db.User, db.Password, db.Ip, db.Port)
	err := mongodb.NewMongoDB(ctx, uri)
	if err != nil {
		logm.FatalfE("%s", err)
		return
	}
}
