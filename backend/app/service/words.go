package service

import (
	"fmt"
	"github.com/cartmanis/english_dictonary/backend/app/provider_db"
)

const (
	english = "english"
)

func ForgetOrRecall(forget bool, idWord string, m *provider_db.MongoClient) (int64, error) {
	if m == nil {
		return -1, fmt.Errorf("не проиницилизованный mongo client")
	}
	idWordObject, err := provider_db.GetObjectId(idWord)
	if err != nil {
		return -1, err
	}
	filter := map[string]interface{}{"_id": idWordObject}
	update := map[string]bool{"forget": forget}
	result, err := m.UpdateOne(filter, update, "$set", english)
	if err != nil {
		return -1, err
	}
	return result.MatchedCount, nil
}

func DeleteWord(idWord string, m *provider_db.MongoClient) (int64, error) {
	if m == nil {
		return -1, fmt.Errorf("не проиницилизованный mongo client")
	}
	idWordObject, err := provider_db.GetObjectId(idWord)
	if err != nil {
		return -1, err
	}
	filter := map[string]interface{}{"_id": idWordObject}
	result, err := m.DeleteOne(filter, english)
	if err != nil {
		return -1, err
	}
	return result.DeletedCount, nil
}
