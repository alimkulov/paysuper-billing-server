package models

import (
	"github.com/bxcodec/faker"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
)

// InitFakeCustomProviders is initialize all custom providers for faker library.
func InitFakeCustomProviders() {
	faker.AddProvider("objectId", FakeProviderObjectId)
	faker.AddProvider("objectIdString", FakeProviderObjectIdString)
	faker.AddProvider("objectIdPointer", FakeProviderObjectIdPointer)
}

// FakeProviderObjectIdPointer returns value with pointer as mongo ObjectID
func FakeProviderObjectIdPointer(_ reflect.Value) (interface{}, error) {
	pointer := primitive.NewObjectID()
	return &pointer, nil
}

// FakeProviderObjectId returns value as mongo ObjectID
func FakeProviderObjectId(_ reflect.Value) (interface{}, error) {
	return primitive.NewObjectID(), nil
}

// FakeProviderObjectIdString returns value as string of mongo ObjectID
func FakeProviderObjectIdString(_ reflect.Value) (interface{}, error) {
	return primitive.NewObjectID().Hex(), nil
}
