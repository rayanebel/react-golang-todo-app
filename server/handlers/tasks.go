package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rayanebel/react-golang-todo-app/server/config"
	"github.com/rayanebel/react-golang-todo-app/server/middleware"
	"github.com/rayanebel/react-golang-todo-app/server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Content-Type", "application/json")
	config := config.Config
	if config.DatabaseEngine == "badger" {
		key := []byte("tasks")
		data, err := middleware.GetItemByKey(key)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		todos := &models.Todo{}
		err = json.Unmarshal(data, todos)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(todos)
	} else if config.DatabaseEngine == "mongo" {
		todos, err := middleware.MongoGetAllData()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(todos)
	}
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	config := config.Config
	if config.DatabaseEngine == "badger" {
		var p models.Task
		key := []byte("tasks")
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		p.ID = primitive.NewObjectID().Hex()
		data, err := middleware.GetItemByKey(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		todos := &models.Todo{}
		err = json.Unmarshal(data, todos)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		todos.Tasks.Items = append(todos.Tasks.Items, p)
		newTodos, _ := json.Marshal(todos)
		err = middleware.NewItemForKey(key, newTodos)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else if config.DatabaseEngine == "mongo" {
		var p models.TodoMongoTask
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = middleware.MongoNewData(&p)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Unable to create task...", http.StatusInternalServerError)
			return
		}
	}
}

func TaskComplete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	config := config.Config
	params := mux.Vars(r)

	ItemID, _ := params["id"]
	if _, ok := params["id"]; !ok {
		http.Error(w, "No ID has been provided", http.StatusBadRequest)
		return
	}

	if config.DatabaseEngine == "badger" {
		key := []byte("tasks")
		data, err := middleware.GetItemByKey(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		todos := &models.Todo{}
		err = json.Unmarshal(data, todos)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, ok := params["id"]; !ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var found bool
		for index, item := range todos.Tasks.Items {
			if item.ID == ItemID {
				todos.Tasks.Items[index].Status = true
				found = true
				break
			}
		}

		if !found {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}

		newTodos, _ := json.Marshal(todos)
		err = middleware.NewItemForKey(key, newTodos)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else if config.DatabaseEngine == "mongo" {
		updatedData := bson.M{"$set": bson.M{"status": true}}
		err := middleware.MongoUpdateData(updatedData, ItemID)
		if err != nil {
			log.Fatal(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func UndoTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	config := config.Config
	params := mux.Vars(r)

	ItemID, _ := params["id"]
	if _, ok := params["id"]; !ok {
		http.Error(w, "No ID has been provided", http.StatusBadRequest)
		return
	}

	if config.DatabaseEngine == "badger" {
		key := []byte("tasks")
		data, err := middleware.GetItemByKey(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		todos := &models.Todo{}
		err = json.Unmarshal(data, todos)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var found bool
		for index, item := range todos.Tasks.Items {
			if item.ID == ItemID {
				todos.Tasks.Items[index].Status = false
				found = true
				break
			}
		}

		if !found {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}

		newTodos, _ := json.Marshal(todos)
		err = middleware.NewItemForKey(key, newTodos)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else if config.DatabaseEngine == "mongo" {
		updatedData := bson.M{"$set": bson.M{"status": false}}
		err := middleware.MongoUpdateData(updatedData, ItemID)
		if err != nil {
			log.Fatal(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	config := config.Config
	params := mux.Vars(r)

	ItemID, _ := params["id"]
	if _, ok := params["id"]; !ok {
		http.Error(w, "No ID has been provided", http.StatusBadRequest)
		return
	}

	if config.DatabaseEngine == "badger" {
		key := []byte("tasks")
		data, err := middleware.GetItemByKey(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		todos := &models.Todo{}
		err = json.Unmarshal(data, todos)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		params := mux.Vars(r)

		if _, ok := params["id"]; !ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var found bool
		for index, item := range todos.Tasks.Items {
			if item.ID == ItemID {
				todos.Tasks.Items = append(todos.Tasks.Items[:index], todos.Tasks.Items[index+1:]...)
				found = true
				break
			}
		}

		if !found {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}

		newTodos, _ := json.Marshal(todos)
		err = middleware.NewItemForKey(key, newTodos)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else if config.DatabaseEngine == "mongo" {
		err := middleware.MongoDeleteData(ItemID)
		if err != nil {
			log.Fatal(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteAllTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	config := config.Config
	if config.DatabaseEngine == "badger" {
		key := []byte("tasks")
		err := middleware.NewItemForKey(key, []byte(`{}`))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if config.DatabaseEngine == "mongo" {
		err := middleware.MongoPurgeAllData()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
