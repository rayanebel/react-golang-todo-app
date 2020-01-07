package router

import (
	"github.com/rayanebel/react-golang-todo-app/server/handlers"

	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/task", handlers.GetAllTask).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/task", handlers.CreateTask).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/task/{id}", handlers.TaskComplete).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/undoTask/{id}", handlers.UndoTask).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deleteTask/{id}", handlers.DeleteTask).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/deleteAllTask", handlers.DeleteAllTask).Methods("DELETE", "OPTIONS")
	return router
}
