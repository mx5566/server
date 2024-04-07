package mongodb

import "go.mongodb.org/mongo-driver/bson"

// 定义过滤器
type filter bson.D

// 匹配字段值大于指定值的文档
func (f filter) GT(key string, value interface{}) filter {
	f = append(f, bson.E{Key: key, Value: bson.D{bson.E{Key: "$gt", Value: value}}})
	return f
}

// 匹配字段值大于等于指定值的文档
func (f filter) GTE(key string, value interface{}) filter {
	f = append(f, bson.E{Key: key, Value: bson.D{{Key: "$gte", Value: value}}})
	return f
}

// 匹配字段值等于指定值的文档
func (f filter) EQ(key string, value interface{}) filter {
	f = append(f, bson.E{Key: key, Value: bson.D{{Key: "$eq", Value: value}}})
	return f
}

// 匹配字段值小于指定值的文档
func (f filter) LT(key string, value interface{}) filter {
	f = append(f, bson.E{Key: key, Value: bson.D{{Key: "$lt", Value: value}}})
	return f
}

// 匹配字段值小于等于指定值的文档
func (f filter) LET(key string, value interface{}) filter {
	f = append(f, bson.E{Key: key, Value: bson.D{{Key: "$let", Value: value}}})
	return f
}

// 匹配字段值不等于指定值的文档，包括没有这个字段的文档
func (f filter) NE(key string, value interface{}) filter {
	f = append(f, bson.E{Key: key, Value: bson.D{{Key: "$ne", Value: value}}})
	return f
}

// 匹配字段值等于指定数组中的任何值
func (f filter) IN(key string, value ...interface{}) filter {
	f = append(f, bson.E{Key: key, Value: bson.D{{Key: "$in", Value: value}}})
	return f
}

// 字段值不在指定数组或者不存在
func (f filter) NIN(key string, value ...interface{}) filter {
	f = append(f, bson.E{Key: key, Value: bson.D{{Key: "$nin", Value: value}}})
	return f
}

// 创建一个条件查询对象
func Newfilter() filter {
	return filter{}
}
