package service

import (
	"time"
	"tm-service/http/responses"
	"tm-service/models"
	"tm-service/producer"
	"tm-service/utils/errors"
	"tm-service/utils/log"

	"gorm.io/gorm"
)

type TaskService struct {
	db *gorm.DB
	mp *producer.MessageProducer
}

func NewTaskService(db *gorm.DB, mp *producer.MessageProducer) (*TaskService, error) {
	if mp == nil {
		return nil, errors.ErrNonExistingService
	}

	models.New(db)
	return &TaskService{db: db, mp: mp}, nil
}

func (ts *TaskService) AddTask(name string, deadline time.Time, userId int) (responses.TaskResponse, error) {
	task := models.Task{
		UserId:    uint(userId),
		Name:      name,
		Deadline:  deadline,
		Completed: false,
	}
	dbErr := task.AddTask()
	errors.DBErrorCheck(dbErr)

	log.Info("Added task %v", task)
	ts.sendKafkaMessage(&task, "add_task")

	return responses.TaskResponse{
		Task: task,
	}, nil
}

func (ts *TaskService) GetAllTasksForUser(id int) (responses.TaskListResponse, error) {
	tasks := []models.Task{}
	dbErr := models.GetAllTasksForUser(uint(id), &tasks)
	errors.DBErrorCheck(dbErr)

	for _, task := range tasks {
		ts.sendKafkaMessage(&task, "get_tasks")
	}

	return responses.TaskListResponse{
		Tasks: tasks,
	}, nil
}

func (ts *TaskService) CompleteTask(taskId, userId int) (responses.TaskResponse, error) {
	task := models.Task{}
	dbErr := models.CompleteTask(taskId, userId, &task)
	errors.DBErrorCheck(dbErr)

	ts.sendKafkaMessage(&task, "complete_task")

	return responses.TaskResponse{
		Task: task,
	}, nil
}

func (ts *TaskService) DeleteTaskById(taskId, userId int) error {
	dbErr := models.DeleteTaskById(userId, taskId)
	errors.DBErrorCheck(dbErr)

	task := models.Task{
		Model: gorm.Model{
			ID: uint(taskId),
		},
		UserId: uint(userId),
	}
	ts.sendKafkaMessage(&task, "delete_task")

	return nil
}

func (ts *TaskService) CheckDeadlines() {
	interval := time.Duration(10) * time.Second

	ticker := time.NewTicker(interval)

	go func() {
		for {
			select {
			case <-ticker.C:
				log.Info("Checking deadlines")
				tasks := []models.Task{}
				dbErr := models.GetAllTasksByDeadline(&tasks)

				if dbErr.Error != nil {
					log.Error("Error getting tasks by deadline: %v", dbErr.Error)
					continue
				}
				log.Info("Got tasks by deadline %v", tasks)

				for _, task := range tasks {
					ts.sendKafkaMessage(&task, "deadline_task")
				}
			}
		}
	}()
}

func (ts *TaskService) sendKafkaMessage(task *models.Task, msgType string) {
	go func(task *models.Task, msgType string) {
		msg := &producer.Message{
			Type:     msgType,
			TaskId:   int(task.ID),
			TaskName: task.Name,
			UserId:   int(task.UserId),
		}
		err := ts.mp.ProduceMessage(msg)
		if err != nil {
			log.Error("Error producing message: %v", err)
		}
	}(task, msgType)
}
