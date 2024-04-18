package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"testing"
)

type Test struct {
	Id            primitive.ObjectID `bson:"_id"`
	Title         string             `bson:"title"`
	Author        string             `bson:"author"`
	YearPublished int64              `bson:"year_published"`
}

func TestMongo(t *testing.T) {
	ctx := context.Background()

	//连接数据库
	err := NewMongoDB(ctx, "mongodb://mengxiang:123456@localhost:27017/?authSource=admin")
	if err != nil {
		log.Fatalf("%s", err)
		return
	}

	//设置使用的库和表
	mgdb := NewMGDB[Test]("test", "favorite_books")

	//插入单条
	insertOneResult := mgdb.InsertOne(ctx, Test{
		Id:            primitive.NewObjectID(),
		Title:         "test",
		Author:        "author test",
		YearPublished: 9999,
	})

	log.Printf("插入单条记录: %v \n", insertOneResult.InsertedID)

	//插入多条
	var tests []Test
	for i := 1; i < 100; i++ {
		tests = append(tests, Test{
			Id:            primitive.NewObjectID(),
			Title:         "test_" + fmt.Sprintf("%d", i),
			Author:        "author test " + fmt.Sprintf("%d", i),
			YearPublished: int64(i),
		})
	}
	insertMultipleResult := mgdb.InsertMultiple(ctx, tests)

	log.Printf("插入多条记录: %v \n", insertMultipleResult.InsertedIDs)

	//查询
	filter := Newfilter().EQ("title", "test").EQ("author", "author test")
	result, err := mgdb.FindOne(ctx, filter)
	if err != nil {
		log.Fatalf("%s", err)
	}
	buf, err := json.Marshal(result)
	fmt.Printf("查询单条记录: %s\n  ", string(buf))

	//查询
	filter = Newfilter().GT("year_published", 5).LT("year_published", 10)
	results, err := mgdb.Find(ctx, filter, 10)
	if err != nil {
		log.Fatalf("%s", err)
	}
	buf, err = json.Marshal(results)
	fmt.Printf("查询多条记录: %v\n  ", string(buf))

	//单条记录更新
	filter = Newfilter().EQ("year_published", 9999)
	updateCount, err := mgdb.UpdateOne(ctx, filter, map[string]interface{}{
		"author": "test 00021",
	})
	if err != nil {
		log.Fatalf("%s", err)

	}
	fmt.Printf("更新数量 : %d\n", updateCount)

	//批量更新
	filter = Newfilter().IN("year_published", 11, 12, 13)
	updateCount, err = mgdb.UpdateMany(ctx, filter, map[string]interface{}{
		"author": "update author",
	})

	if err != nil {
		log.Fatalf("%s", err)
	}
	fmt.Printf("批量更新数量 : %d\n", updateCount)

	//单条数据删除
	filter = Newfilter().EQ("year_published", 15)
	deleteCount, err := mgdb.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatalf("%s", err)
	}
	fmt.Printf("单条数据删除数量 : %d\n", deleteCount)

	//多条数据删除
	filter = Newfilter().IN("year_published", 16, 17, 18)
	deleteCount, err = mgdb.DeleteMany(ctx, filter)
	if err != nil {
		log.Fatalf("%s", err)
	}
	fmt.Printf("多条数据删除数量 : %d\n", deleteCount)

}
