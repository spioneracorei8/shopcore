package domain

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

const RUN_NUMBER_COLLECTION = "run_number"

type RunNumber struct {
	Id      *bson.ObjectID `bson:"_id" json:"_id"`
	Prefix  string         `bson:"prefix" json:"prefix"`
	Running int            `bson:"running" json:"running"`
}

func (rn *RunNumber) GenObjectID() {
	id := bson.NewObjectID()
	rn.Id = &id
}
