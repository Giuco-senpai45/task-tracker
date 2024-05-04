package responses

import "tm-service/models"

type TaskResponse struct {
	Task models.Task `json:"task"`
}

type TaskListResponse struct {
	Tasks []models.Task `json:"tasks"`
}
