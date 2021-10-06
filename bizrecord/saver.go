package bizrecord

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//MongoSaver ...
type MongoSaver struct {
	*mongo.Client
}

//MongoConfig ...
type MongoConfig struct {
	DNS string
}

//DefaultMongoConfig ...
func DefaultMongoConfig() MongoConfig {
	return MongoConfig{DNS: "mongodb://localhost:27017"}
}

//NewMongoSaver ...
func NewMongoSaver(conf MongoConfig) MongoSaver {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.DNS))
	if err != nil {
		panic(err)
	}
	return MongoSaver{client}
}

func (s MongoSaver) Save(recordIF RecordIF) bool {

	return false
}

func (s *MongoSaver) Load(f func(saverIF SaverIF) RecordIF) RecordIF {
	return f(s)
}

//MysqlSaver ...
type MysqlSaver struct {
}

func (s MysqlSaver) Save(recordIF RecordIF) bool {
	return false
}

func (s *MysqlSaver) Load(f func(saverIF SaverIF) RecordIF) RecordIF {
	return f(s)
}
