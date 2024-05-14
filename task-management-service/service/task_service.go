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

func NewTaskService(db *gorm.DB, mp *producer.MessageProducer) *TaskService {
	models.New(db)
	return &TaskService{db: db, mp: mp}
}

func (ts *TaskService) AddTask(name string, deadline time.Time, userId int) (responses.TaskResponse, error) {
	task := models.Task{
		UserId:    uint(userId),
		Name:      name,
		Deadline:  deadline,
		Completed: false,
	}
	err := task.AddTask()
	errors.DBErrorCheck(err)

	return responses.TaskResponse{
		Task: task,
	}, nil
}

func (ts *TaskService) GetAllTasksForUser(id int) (responses.TaskListResponse, error) {
	tasks := []models.Task{}
	err := models.GetAllTasksForUser(uint(id), &tasks)
	errors.DBErrorCheck(err)

	return responses.TaskListResponse{
		Tasks: tasks,
	}, nil
}

func (ts *TaskService) CompleteTask(taskId, userId int) (responses.TaskResponse, error) {
	task := models.Task{}
	err := models.CompleteTask(taskId, userId, &task)
	errors.DBErrorCheck(err)

	return responses.TaskResponse{
		Task: task,
	}, nil
}

func (ts *TaskService) DeleteTaskById(taskId, userId int) error {
	err := models.DeleteTaskById(userId, taskId)
	errors.DBErrorCheck(err)

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
					msg := &producer.Message{
						TaskId:   int(task.ID),
						TaskName: task.Name,
						UserId:   int(task.UserId),
					}
					err := ts.mp.ProduceMessage(msg)
					if err != nil {
						log.Error("Error producing message: %v", err)
					}
				}
			}
		}
	}()
}
