package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"todo-app/actor"
	"todo-app/storage"
)

var actorInstance *actor.Actor

func InitActor(ctx context.Context) {
	actorInstance = actor.NewActor(ctx)
}

func AddRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/create", createItemHandler)
	mux.HandleFunc("/update", updateItemHandler)
	mux.HandleFunc("/delete", deleteItemHandler)
	mux.HandleFunc("/get/{itemid}", getByIDHandler)
	mux.HandleFunc("/get", getListHandler)
	mux.HandleFunc("/about", aboutPageHandler)
	mux.HandleFunc("/list", dynamicListHandler)
}

func getListHandler(w http.ResponseWriter, r *http.Request) {
	if actorInstance == nil {
		http.Error(w, "Actor not initialized", http.StatusInternalServerError)
		return
	}
	items, err := actorInstance.ListAll(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	todos := make([]storage.Item, 0, len(items))
	for _, v := range items {
		todos = append(todos, storage.Item{
			ID:          v.ID,
			Description: v.Description,
			Status:      v.Status,
			Created:     v.Created,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func getByIDHandler(w http.ResponseWriter, r *http.Request) {
	if actorInstance == nil {
		http.Error(w, "Actor not initialized", http.StatusInternalServerError)
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Missing item ID", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}
	item, err := actorInstance.List(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(storage.Item{
		ID:          item.ID,
		Description: item.Description,
		Status:      item.Status,
		Created:     item.Created,
	})
}

func createItemHandler(w http.ResponseWriter, r *http.Request) {
	if actorInstance == nil {
		http.Error(w, "Actor not initialized", http.StatusInternalServerError)
		return
	}
	var todo storage.Item
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	item, err := actorInstance.Create(context.Background(), todo.Description, todo.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(storage.Item{
		ID:          item.ID,
		Description: item.Description,
		Status:      item.Status,
		Created:     item.Created,
	})
}

func updateItemHandler(w http.ResponseWriter, r *http.Request) {
	if actorInstance == nil {
		http.Error(w, "Actor not initialized", http.StatusInternalServerError)
		return
	}
	var todo storage.Item
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	item, err := actorInstance.Update(context.Background(), todo.ID, todo.Description, todo.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(storage.Item{
		ID:          item.ID,
		Description: item.Description,
		Status:      item.Status,
		Created:     item.Created,
	})
}

func deleteItemHandler(w http.ResponseWriter, r *http.Request) {
	if actorInstance == nil {
		http.Error(w, "Actor not initialized", http.StatusInternalServerError)
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Missing item ID", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}
	err = actorInstance.Delete(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"deleted": id})
}

func aboutPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "About Todo App")
}

func dynamicListHandler(w http.ResponseWriter, r *http.Request) {
	getListHandler(w, r)
}
