package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Result bson.M

func (r *Result) GetId() (id primitive.ObjectID) {
	res, ok := r.Value()["_id"]
	if ok {
		id, ok = res.(primitive.ObjectID)
		if ok {
			return id
		}
	}
	return id
}

func (r *Result) SetId(id primitive.ObjectID) {
	(*r)["_id"] = id
}

func (r *Result) Value() map[string]interface{} {
	return *r
}

func (r *Result) SetValue(key string, value interface{}) {
	(*r)[key] = value
}

func (r *Result) GetValue(key string) (value interface{}) {
	return (*r)[key]
}

func (r *Result) GetTaskId() (id primitive.ObjectID) {
	res := r.GetValue("_tid")
	if res == nil {
		return id
	}
	id, _ = res.(primitive.ObjectID)
	return id
}

func (r *Result) SetTaskId(id primitive.ObjectID) {
	r.SetValue("_tid", id)
}
