package task

import (
	"errors"
	"fmt"

	"tasklybe/pkg/db"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrValidation = errors.New("validation error")

func GetTasks(page int, limit int) (*[]Task, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	var total int64
	if err := db.DB.Model(&Task{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	var tasks []Task
	if err := db.DB.Order("created_at desc").Limit(limit).Offset(offset).Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return &tasks, total, nil
}

func GetTask(id string) (*Task, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", ErrValidation)
	}

	var task Task
	if err := db.DB.First(&task, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

func CreateTask(input CreateTaskRequest) (*Task, error) {
	task := Task{
		ID:     uuid.NewString(),
		UserID: input.UserID,
		Title:  input.Title,
		Desc:   input.Desc,
		Label:  input.Label,
	}

	if err := db.DB.Create(&task).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

func EditTask(id string, input EditTaskRequest) (*Task, error) {

	var task Task
	if err := db.DB.First(&task, "id = ?", id).Error; err != nil {
		return nil, err
	}

	task.Title = input.Title
	task.Desc = input.Desc
	task.Label = input.Label

	if err := db.DB.Save(&task).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

func DeleteTask(id string) error {
	tx := db.DB.Delete(&Task{}, "id = ?", id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
