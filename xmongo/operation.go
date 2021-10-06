package xmongo

import (
	"bufio"
	"bytes"
	"context"

	jsoniter "github.com/json-iterator/go"

	"go.mongodb.org/mongo-driver/x/bsonx"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCOll(client *mongo.Client, dbName, collName string) *mongo.Collection {
	collection := client.Database(dbName).Collection(collName)
	return collection
}

func InsertOne(ctx context.Context, coll *mongo.Collection, data interface{}) (*mongo.InsertOneResult, error) {
	res, err := coll.InsertOne(ctx, data)
	return res, err
}

func InsertMany(coll *mongo.Collection, data []interface{}) (*mongo.InsertManyResult, error) {

	res, err := coll.InsertMany(context.Background(), data)
	return res, err
}

func FindOne(coll *mongo.Collection, filter interface{}) *mongo.SingleResult {

	return coll.FindOne(context.Background(), filter)
}

func Find(coll *mongo.Collection, filter interface{}) (*mongo.Cursor, error) {

	return coll.Find(context.Background(), filter)
}

func FindWithOption(coll *mongo.Collection, filter interface{},
	findOptions *options.FindOptions) (*mongo.Cursor, error) {

	return coll.Find(context.Background(), filter, findOptions)
}

func Distinct(coll *mongo.Collection, fieldName string, filter interface{}) ([]interface{}, error) {

	return coll.Distinct(context.Background(), fieldName, filter)
}

func UpdateOne(coll *mongo.Collection, filter interface{}, data interface{}) (*mongo.UpdateResult, error) {

	return coll.UpdateOne(context.Background(), filter, data)
}

func UpdateMany(coll *mongo.Collection, filter interface{}, data interface{}) (*mongo.UpdateResult, error) {

	return coll.UpdateMany(context.Background(), filter, data)
}

func UpdateByID(coll *mongo.Collection, id interface{}, data interface{}) (*mongo.UpdateResult, error) {

	return coll.UpdateByID(context.Background(), id, data)
}

func UpdateOneWithSession(client *mongo.Client, coll *mongo.Collection, filter interface{}, data interface{}) error {

	var (
		session mongo.Session
		err     error
		ctx     = context.Background()
	)
	session, err = client.StartSession()
	if err != nil {
		return err
	}
	if err = session.StartTransaction(); err != nil {
		return err
	}
	f := func(sessionContext mongo.SessionContext) error {
		_, err = coll.UpdateOne(sessionContext, filter, data)
		if err != nil {
			return err
		}
		err = session.CommitTransaction(sessionContext)
		if err != nil {
			return err
		}
		return nil
	}
	err = mongo.WithSession(ctx, session, f)
	if err != nil {
		return err
	}
	session.EndSession(ctx)
	return nil
}

func UpdateManyWithSession(client *mongo.Client, coll *mongo.Collection, filter interface{}, data interface{}) error {

	var (
		session mongo.Session
		err     error
		ctx     = context.Background()
	)
	session, err = client.StartSession()
	if err != nil {
		return err
	}
	if err = session.StartTransaction(); err != nil {
		return err
	}
	f := func(sessionContext mongo.SessionContext) error {
		_, err = coll.UpdateMany(sessionContext, filter, data)
		if err != nil {
			return err
		}
		err = session.CommitTransaction(sessionContext)
		if err != nil {
			return err
		}
		return nil
	}
	err = mongo.WithSession(ctx, session, f)
	if err != nil {
		return err
	}
	session.EndSession(ctx)
	return nil
}

func UpdateByIDWithSession(client *mongo.Client, coll *mongo.Collection, id interface{}, data interface{}) error {

	var (
		session mongo.Session
		err     error
		ctx     = context.Background()
	)
	session, err = client.StartSession()
	if err != nil {
		return err
	}
	if err = session.StartTransaction(); err != nil {
		return err
	}
	f := func(sessionContext mongo.SessionContext) error {
		_, err = coll.UpdateByID(sessionContext, id, data)
		if err != nil {
			return err
		}
		err = session.CommitTransaction(sessionContext)
		if err != nil {
			return err
		}
		return nil
	}
	err = mongo.WithSession(ctx, session, f)
	if err != nil {
		return err
	}
	session.EndSession(ctx)
	return nil
}

func DeleteOne(coll *mongo.Collection, filter interface{}) (*mongo.DeleteResult, error) {

	return coll.DeleteOne(context.Background(), filter)
}

func DeleteMany(coll *mongo.Collection, filter interface{}) (*mongo.DeleteResult, error) {

	return coll.DeleteMany(context.Background(), filter)
}

func Count(coll *mongo.Collection, filter interface{}) (int64, error) {

	return coll.CountDocuments(context.Background(), filter)
}

func ChangeStreamClient(client *mongo.Client, coll *mongo.Collection) {
	//
}

func ChangeStreamCollection(client *mongo.Client, coll *mongo.Collection) {
	//
}

func ChangeStreamDB(client *mongo.Client, coll *mongo.Collection) {
	//
}

//UploadGridFS ...
func UploadGridFS(filename string, data interface{}, db *mongo.Database, bucketOptions *options.BucketOptions) error {
	bucket, err := gridfs.NewBucket(db, bucketOptions)
	if err != nil {
		return err
	}
	opts := options.GridFSUpload()
	opts.SetMetadata(bsonx.Doc{{Key: "content-type", Value: bsonx.String("application/json")}})
	var upLoadStream *gridfs.UploadStream
	if upLoadStream, err = bucket.OpenUploadStream(filename, opts); err != nil {
		return err
	}
	str, err := jsoniter.MarshalToString(data)
	if err != nil {
		return err
	}
	if _, err = upLoadStream.Write([]byte(str)); err != nil {
		return err
	}
	upLoadStream.Close()
	return nil
}

//DownLoadGridFS ...
func DownLoadGridFS(fileID interface{}, db *mongo.Database, bucketOptions *options.BucketOptions) (string, error) {
	bucket, err := gridfs.NewBucket(db, bucketOptions)
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	if _, err = bucket.DownloadToStream(fileID, w); err != nil {
		return "", err
	}
	return b.String(), err
}
