package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

func getAllTasks(w http.ResponseWriter, _ *http.Request) {
	tasksData, err := json.Marshal(tasks)

	if err != nil {
		http.Error(w, fmt.Errorf("JSON encoding error: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(tasksData)

	if err != nil {
		http.Error(w, fmt.Errorf("writing tasks data error %w", err).Error(), http.StatusInternalServerError)
	}
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var taskData Task
	var buffer bytes.Buffer

	if _, err := buffer.ReadFrom(r.Body); err != nil {
		http.Error(w, fmt.Errorf("body getting error: %w", err).Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buffer.Bytes(), &taskData); err != nil {
		http.Error(w, fmt.Errorf("JSON encoding error: %w", err).Error(), http.StatusBadRequest)
		return
	}

	if _, ok := tasks[taskData.ID]; ok {
		http.Error(w, fmt.Sprintf("task with id = %s already exists", taskData.ID), http.StatusBadRequest)
		return
	}

	tasks[taskData.ID] = taskData

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func getTaskById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task, ok := tasks[id]
	if !ok {
		http.Error(w, fmt.Sprintf("failed to find task with id = %s ", id), http.StatusBadRequest)
		return
	}

	taskData, err := json.Marshal(task)
	if err != nil {
		http.Error(w, fmt.Errorf("JSON encoding error: %w", err).Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(taskData)

	if err != nil {
		http.Error(w, fmt.Errorf("writing tasks data error: %w", err).Error(), http.StatusBadRequest)
	}
}

func deleteTaskById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, ok := tasks[id]
	if !ok {
		http.Error(w, fmt.Sprintf("failed to find task with id = %s", id), http.StatusBadRequest)
		return
	}

	delete(tasks, id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", getAllTasks)
	r.Post("/tasks", createTask)
	r.Get("/tasks/{id}", getTaskById)
	r.Delete("/tasks/{id}", deleteTaskById)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
