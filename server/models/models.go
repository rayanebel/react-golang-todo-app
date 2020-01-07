package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Used for badgerDB
type Todo struct {
	Tasks Tasks `json:"todo,omitempty"`
}

type Tasks struct {
	Items []Task `json:"tasks"`
}
type Task struct {
	ID     string `json:"id"`
	Name   string `json:"task"`
	Status bool   `json:"status"`
}

// Used for MongoDB
type TodoMongo struct {
	Tasks TodoMongoTasks `json:"todo,omitempty"`
}
type TodoMongoTasks struct {
	Items []TodoMongoTask `json:"tasks"`
}

type TodoMongoTask struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	Name   string             `json:"task"`
	Status bool               `json:"status"`
}
