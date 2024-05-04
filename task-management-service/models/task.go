package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	UserId    uint      `json:"user_id" gorm:"index"`
	Name      string    `json:"name"`
	Deadline  time.Time `json:"deadline"`
	Completed bool      `json:"completed" gorm:"default:false"`
}

func GetAllTasksForUser(userId uint, tasks *[]Task) *gorm.DB {
	db.Where("user_id = ?", userId).Find(&tasks)
	return db
}

func GetAllTasksByDeadline(tasks *[]Task) *gorm.DB {
	now := time.Now()
	end := now.Add(24 * time.Hour) // End time is current time plus 24 hours

	db.Where("deadline BETWEEN ? AND ?", now, end).Find(&tasks)
	return db
}

func CompleteTask(taskId, userId int, t *Task) *gorm.DB {
	db.Model(&t).Where("ID = ?", taskId).Where("user_id = ?", userId).Update("completed", true)
	return db
}

func DeleteTaskById(userId, taskId int) *gorm.DB {
	db.Where("ID = ?", taskId).Where("user_id = ?", userId).Delete(&Task{})
	return db
}

func (t *Task) AddTask() *gorm.DB {
	db.Create(t)
	return db
}
