package helper

import "go.mongodb.org/mongo-driver/bson/primitive"

func IsIdentified(id string) bool {
	_, err := primitive.ObjectIDFromHex(id)
	return err == nil
}
