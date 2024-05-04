package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tm-service/http/requests"
	"tm-service/service"
	"tm-service/utils"
	"tm-service/utils/log"
)

type TaskController struct {
	service *service.TaskService
}

func NewTaskController(service *service.TaskService) *TaskController {
	return &TaskController{service: service}
}

func getUserIdFromHeader(req *http.Request) int {
	userId := req.Header.Get("X-User-Id")

	id, err := strconv.Atoi(userId)
	if err != nil {
		log.Error("Error parsing user-id: %v", err)
		return 0
	}

	return id
}

func (tc *TaskController) AddTask(w http.ResponseWriter, r *http.Request) {
	userId := getUserIdFromHeader(r)
	if userId == 0 {
		http.Error(w, "Incorrect user id", http.StatusBadRequest)
		return
	}
	log.Info("Add task user-id: %v", userId)

	taskRequest := requests.TaskCreatePayload{}
	err := json.NewDecoder(r.Body).Decode(&taskRequest)
	if err != nil {
		http.Error(w, "Error decoding the body", http.StatusInternalServerError)
		return
	}

	res, err := tc.service.AddTask(taskRequest.Name, taskRequest.Deadline, userId)
	if err != nil {
		log.Error("Error adding task: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	taskJson, err := json.Marshal(res)
	if err != nil {
		log.Error("Error marshaling task: %v", err)
		http.Error(w, "Error converting tasks to JSON", http.StatusInternalServerError)
		return

	}
	w.Header().Set("Content-Type", "application/json")

	utils.Ok(taskJson, w)
}

func (tc *TaskController) GetAllTasksForUser(w http.ResponseWriter, r *http.Request) {
	userId := getUserIdFromHeader(r)
	if userId == 0 {
		http.Error(w, "Incorrect user id", http.StatusBadRequest)
		return
	}
	log.Info("Get all tasks user-id: %v", userId)

	res, err := tc.service.GetAllTasksForUser(userId)
	if err != nil {
		log.Error("Error getting tasks for user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tasksJson, err := json.Marshal(res)
	if err != nil {
		log.Error("Error marshaling tasks: %v", err)
		http.Error(w, "Error converting tasks to JSON", http.StatusInternalServerError)
		return

	}
	w.Header().Set("Content-Type", "application/json")

	utils.Ok(tasksJson, w)
}

func (tc *TaskController) CompleteTask(w http.ResponseWriter, r *http.Request) {
	taskId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Error("Incorrect path variable: %v", err)
		http.Error(w, "Incorrect path variable taskId", http.StatusBadRequest)
		return
	}
	userId := getUserIdFromHeader(r)
	if userId == 0 {
		http.Error(w, "Incorrect user id", http.StatusBadRequest)
		return
	}

	log.Warn("Got user id %v and taskId %v", userId, taskId)

	res, err := tc.service.CompleteTask(taskId, userId)
	if err != nil {
		log.Error("Error completing task: %v", err)
		http.Error(w, "Error completing task", http.StatusInternalServerError)
		return
	}

	taskJson, err := json.Marshal(res)
	if err != nil {
		log.Error("Error marshaling task: %v", err)
		http.Error(w, "Error converting task to JSON", http.StatusInternalServerError)
		return
	}

	utils.Ok(taskJson, w)
}

func (tc *TaskController) DeleteTaskById(w http.ResponseWriter, r *http.Request) {
	taskId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Error("Incorrect path variable: %v", err)
		http.Error(w, "Incorrect path variable taskId", http.StatusBadRequest)
		return
	}
	userId := getUserIdFromHeader(r)
	if userId == 0 {
		http.Error(w, "Incorrect user id", http.StatusBadRequest)
		return
	}

	err = tc.service.DeleteTaskById(taskId, userId)
	if err != nil {
		log.Error("Error deleting task: %v", err)
		http.Error(w, "Something went wrong task couldn't be deleted", http.StatusInternalServerError)
		return
	}

	utils.Ok(nil, w)
}
