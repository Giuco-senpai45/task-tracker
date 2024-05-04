package routes

import (
	"net/http"
	"tm-service/controller"
)

var RegisterTaskRoutes = func(router *http.ServeMux, controller *controller.TaskController) {
	router.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	router.HandleFunc("GET /tasks", controller.GetAllTasksForUser)
	router.HandleFunc("POST /tasks", controller.AddTask)
	router.HandleFunc("PUT /tasks/{id}", controller.CompleteTask)
	router.HandleFunc("DELETE /tasks/{id}", controller.DeleteTaskById)
}
