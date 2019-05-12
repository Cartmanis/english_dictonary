package provider_db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"time"
)

type MongoClient struct {
	client     *mongo.Client
	ctx        context.Context
	cancelFunc *context.CancelFunc
	db         string
}

func NewStoreContext(url, dbName, login, pass string, interval time.Duration) (*MongoClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), interval*time.Second)
	auth := &options.Credential{Username: login, Password: pass}
	auth.AuthSource = dbName
	opts := options.ClientOptions{Auth: auth}
	client, err := mongo.Connect(ctx, opts.ApplyURI(url))
	if err != nil {
		return nil, err
	}
	//пинг до базы данных, чтобы убедиться, что подключение утсановлено
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return &MongoClient{
		client:     client,
		ctx:        ctx,
		cancelFunc: &cancel,
		db:         dbName,
	}, nil
}

func (m *MongoClient) Ping() bool {
	if err := m.client.Ping(m.ctx, nil); err != nil {
		return false
	}
	return true
}

func (m *MongoClient) Close() {
	if m != nil && m.client != nil && m.ctx != nil {
		if err := m.client.Disconnect(m.ctx); err != nil {
			fmt.Println("не удалось закрыть соединение mongoDb")
		}
	}
}

func (m *MongoClient) FindOne(filter interface{}, result interface{}, collectionName string) error {
	collection := m.client.Database(m.db).Collection(collectionName)
	return collection.FindOne(m.ctx, filter).Decode(result)
}

func (m *MongoClient) Find(filter interface{}, results interface{}, collName string) error {
	collection := m.client.Database(m.db).Collection(collName)
	cursor, err := collection.Find(m.ctx, filter)
	if err != nil {
		return err
	}
	defer cursor.Close(m.ctx)
	if err := cursor.All(m.ctx, results); err != nil {
		return err
	}
	if err := cursor.Err(); err != nil {
		return err
	}
	return nil
}

func (m *MongoClient) ListIndexes(collectionName string) ([]string, error) {
	collection := m.client.Database(m.db).Collection(collectionName)
	cur, err := collection.Indexes().List(m.ctx)
	if err != nil {
		return nil, err
	}
	defer cur.Close(m.ctx)
	type keys struct {
		Name string
	}
	listIdx := make([]string, 0)
	for cur.Next(m.ctx) {
		var elem keys
		if err := cur.Decode(&elem); err != nil {
			return nil, err
		}
		listIdx = append(listIdx, elem.Name)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return listIdx, nil
}

func (m *MongoClient) Count(filter interface{}, collectionName string) (int64, error) {
	collection := m.client.Database(m.db).Collection(collectionName)
	return collection.CountDocuments(m.ctx, filter)
}

func (m *MongoClient) InsertOne(document interface{}, collectionName string) (interface{}, error) {
	collection := m.client.Database(m.db).Collection(collectionName)
	res, err := collection.InsertOne(m.ctx, document)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("Запись добавить не удалось. Пустой овет")
	}
	return res.InsertedID, nil
}

func (m *MongoClient) CreateIndex(field, collectionName string, unique bool) (string, error) {
	collection := m.client.Database(m.db).Collection(collectionName)
	keys := bsonx.Doc{{Key: field, Value: bsonx.Int32(int32(1))}}
	index := mongo.IndexModel{}
	index.Keys = keys
	opts := &options.IndexOptions{Unique: &unique}
	name := "idx_" + field
	opts.Name = &name
	index.Options = opts
	return collection.Indexes().CreateOne(m.ctx, index)
}

func (m *MongoClient) InsertMany(documents []interface{}, collectionName string) ([]interface{}, error) {
	collection := m.client.Database(m.db).Collection(collectionName)
	res, err := collection.InsertMany(m.ctx, documents)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("Записи добавить не удалось. Пустой ответ")
	}
	return res.InsertedIDs, nil
}

func (m *MongoClient) UpdateOne(filter interface{}, update interface{}, keyUpdate string, collectionName string) (*mongo.UpdateResult, error) {
	collection := m.client.Database(m.db).Collection(collectionName)
	updateBson := bson.D{{keyUpdate, update}}

	filterBson, err := bson.Marshal(filter)
	if err != nil {
		return nil, err
	}
	return collection.UpdateOne(m.ctx, filterBson, updateBson)
}

func (m *MongoClient) DeleteOne(filter interface{}, collectionName string) (*mongo.DeleteResult, error) {
	collection := m.client.Database(m.db).Collection(collectionName)
	filterBson, err := bson.Marshal(filter)
	if err != nil {
		return nil, err
	}
	return collection.DeleteOne(m.ctx, filterBson)
}

func (m *MongoClient) DeleteMany(filter interface{}, collectionName string) (*mongo.DeleteResult, error) {
	collection := m.client.Database(m.db).Collection(collectionName)
	filterBson, err := bson.Marshal(filter)
	if err != nil {
		return nil, err
	}
	return collection.DeleteMany(m.ctx, filterBson)
}

func GetObjectId(id string) (interface{}, error) {
	return primitive.ObjectIDFromHex(id)
}
