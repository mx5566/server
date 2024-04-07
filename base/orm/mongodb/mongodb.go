package mongodb

import (
	"context"
	"github.com/mx5566/logm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func NewMongoDB(ctx context.Context, appUri string) error {
	//连接到mongodb
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(appUri))
	if err != nil {
		logm.PanicfE("mongodb 数据库连接失败: %s", err.Error())
		return err
	}
	//检查连接
	err = c.Ping(ctx, nil)
	if err != nil {
		logm.PanicfE("mongodb 数据库ping失败: %s", err.Error())
		return err
	}
	client = c
	logm.InfofE("mongodb连接成功")
	return nil
}

type MongoDB[T any] struct {
	database   string
	collection string
}

func NewMGDB[T any](database string, collection string) *MongoDB[T] {
	if client == nil {
		logm.FatalfE("mongo.Client Not initialized !")
	}
	return &MongoDB[T]{
		database,
		collection,
	}
}

// 新增一条记录
func (mg *MongoDB[T]) InsertOne(ctx context.Context, value T) *mongo.InsertOneResult {
	result, err := mg.getCollection().InsertOne(ctx, value)
	if err != nil {
		panic(err)
	}
	return result
}

// 新增多条记录
func (mg *MongoDB[T]) InsertMultiple(ctx context.Context, data []T) *mongo.InsertManyResult {
	var array []interface{}
	for i := 0; i < len(data); i++ {
		array = append(array, data[i])
	}
	result, err := mg.getCollection().InsertMany(ctx, array)
	if err != nil {
		panic(err)
	}
	return result
}

// 根据字段名和值查询一条记录
func (mg *MongoDB[T]) FindOne(ctx context.Context, filter filter) (T, error) {
	var t T
	err := mg.getCollection().FindOne(ctx, filter).Decode(&t)
	if err != nil {
		return t, err
	}
	return t, nil
}

// 根据条件查询多条记录
func (mg *MongoDB[T]) Find(ctx context.Context, filter filter, limit int64) ([]T, error) {
	findOpts := options.Find()
	findOpts.SetLimit(limit)
	cursor, err := mg.getCollection().Find(ctx, filter, findOpts)
	var ts []T
	if err != nil {
		return ts, err
	}
	for cursor.Next(ctx) {
		var t T
		err := cursor.Decode(&t)
		if err != nil {
			return ts, err
		}
		ts = append(ts, t)
	}
	cursor.Close(ctx)
	return ts, nil
}

// 根据条件更新
func (mg *MongoDB[T]) UpdateOne(ctx context.Context, filter filter, update interface{}) (int64, error) {
	result, err := mg.getCollection().UpdateOne(ctx, filter, bson.M{"$set": update})
	return result.ModifiedCount, err
}

// 根据id更新
func (mg *MongoDB[T]) UpdateOneById(ctx context.Context, id string, update interface{}) (int64, error) {
	result, err := mg.getCollection().UpdateOne(ctx, filter{{Key: "_id", Value: mg.ObjectID(id)}}, update)
	return result.ModifiedCount, err
}

// 更新多个
func (mg *MongoDB[T]) UpdateMany(ctx context.Context, filter filter, update interface{}) (int64, error) {
	result, err := mg.getCollection().UpdateMany(ctx, filter, bson.D{{Key: "$set", Value: update}})
	return result.ModifiedCount, err
}

// 获取表
func (mg *MongoDB[T]) getCollection() *mongo.Collection {
	return client.Database(mg.database).Collection(mg.collection)
}

// 删除一条记录
func (mg *MongoDB[T]) DeleteOne(ctx context.Context, filter filter) (int64, error) {
	result, err := mg.getCollection().DeleteOne(ctx, filter)
	return result.DeletedCount, err
}

// 根据id删除一条记录
func (mg *MongoDB[T]) DeleteOneById(ctx context.Context, id string) (int64, error) {
	result, err := mg.getCollection().DeleteOne(ctx, filter{{Key: "_id", Value: mg.ObjectID(id)}})
	return result.DeletedCount, err
}

// 删除多条记录
func (mg *MongoDB[T]) DeleteMany(ctx context.Context, filter filter) (int64, error) {
	result, err := mg.getCollection().DeleteMany(ctx, filter)
	return result.DeletedCount, err
}

// objcetid
func (mg *MongoDB[T]) ObjectID(id string) primitive.ObjectID {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logm.PanicE(err)
	}
	return objectId
}
