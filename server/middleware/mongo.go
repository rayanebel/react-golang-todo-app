package middleware

import (
	"context"
	"fmt"
	"log"

	"github.com/rayanebel/react-golang-todo-app/server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDatabase *mongo.Client
var MongoTodoCollection *mongo.Collection

// create connection with mongo db
func InitMongoDB(connectionString string) {

	clientOptions := options.Client().ApplyURI(connectionString)
	MongoDatabase, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = MongoDatabase.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	MongoTodoCollection = MongoDatabase.Database("todos").Collection("tasks")

	fmt.Println("Collection instance created!")
}

func MongoNewData(task *models.TodoMongoTask) error {
	task.ID = primitive.NewObjectID()
	insertResult, err := MongoTodoCollection.InsertOne(context.Background(), task)

	if err != nil {
		return err
	}

	fmt.Println("Inserted a Single Record ", insertResult.InsertedID)
	return nil
}

func MongoGetAllData() (models.TodoMongo, error) {
	todos := models.TodoMongo{}
	cur, err := MongoTodoCollection.Find(context.Background(), bson.D{{}})
	if err != nil {
		return todos, err
	}
	for cur.Next(context.Background()) {
		var task models.TodoMongoTask
		e := cur.Decode(&task)
		if e != nil {
			log.Fatal(e)
		}
		todos.Tasks.Items = append(todos.Tasks.Items, task)
	}
	return todos, nil
}

func MongoFindDataByID(id string) models.TodoMongoTask {
	task := models.TodoMongoTask{}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}
	value := MongoTodoCollection.FindOne(context.Background(), bson.M{"id": objectID})
	err = value.Decode(&task)
	if err != nil {
		log.Fatal(err)
	}
	return task
}

func MongoUpdateData(updatedData bson.M, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)

	_, err = MongoTodoCollection.UpdateOne(context.TODO(), bson.M{"id": objectID}, updatedData)
	if err != nil {
		return err
	}
	return nil
}

func MongoDeleteData(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)

	_, err = MongoTodoCollection.DeleteOne(context.TODO(), bson.M{"id": objectID})

	if err != nil {
		return err
	}
	return nil
}

func MongoPurgeAllData() error {
	err := MongoTodoCollection.Drop(context.TODO())

	if err != nil {
		return err
	}
	return nil
}
