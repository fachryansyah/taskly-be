package task

import (
	"errors"
	"strconv"
	"tasklybe/internal/dto"
	"tasklybe/internal/validation"

	"github.com/gofiber/fiber/v2"
)

type paginationResponse struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
	Pages int   `json:"pages"`
}

// GetTasks godoc
// @Summary Get tasks
// @Description Get all tasks
// @Tags task
// @Produce json
// @Success 200 {array} Task
// @Router /task [get]
// @Param page query int true "Page number"
// @Param limit query int true "Limit per page"
func HandleGetTasks(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	tasks, total, err := GetTasks(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	pages := 0
	if limit > 0 {
		pages = int((total + int64(limit) - 1) / int64(limit))
	}

	return c.JSON(dto.ResponseWrapper[[]Task]{
		Data:    tasks,
		Success: true,
		Message: "Success! tasks found.",
		Pagination: &dto.PaginationResponse{
			Page:      page,
			Limit:     limit,
			Total:     total,
			TotalPage: pages,
		},
	})
}

// GetTask godoc
// @Summary Get task
// @Description Get task by id
// @Tags task
// @Produce json
// @Success 200 {object} Task
// @Router /task/{id} [get]
// @Param id path string true "Task ID"
func HandleGetTask(c *fiber.Ctx) error {
	id := c.Params("id")
	task, err := GetTask(id)
	if err != nil {
		if IsNotFound(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "task not found"})
		}
		if errors.Is(err, ErrValidation) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(dto.ResponseWrapper[Task]{
		Data:       task,
		Success:    true,
		Message:    "Success! task found.",
		Pagination: nil,
	})
}

// CreateTask godoc
// @Summary Create task
// @Description Create a new task
// @Tags task
// @Produce json
// @Success 200 {object} Task
// @Router /task [post]
// @Param req body CreateTaskRequest true "Create Task Request"
func HandleCreateTask(c *fiber.Ctx) error {
	var req CreateTaskRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ResponseWrapper[Task]{
			Data:    nil,
			Success: false,
			Message: "Failed! Validation error.",
			Error:   validation.FormatValidationError(err),
		})
	}

	userId := c.Locals("userId").(string)
	task, err := CreateTask(userId, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ResponseWrapper[Task]{
			Data:    nil,
			Success: false,
			Message: "Failed! Something went wrong.",
		})
	}

	return c.JSON(dto.ResponseWrapper[Task]{
		Data:    task,
		Success: true,
		Message: "Success! task created.",
	})
}

// EditTask godoc
// @Summary Edit task
// @Description Edit task by id
// @Tags task
// @Produce json
// @Success 200 {object} Task
// @Router /task/{id} [put]
// @Param id path string true "Task ID"
// @Param req body EditTaskRequest true "Edit Task Request"
func HandleEditTask(c *fiber.Ctx) error {
	id := c.Params("id")

	var req EditTaskRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ResponseWrapper[Task]{
			Data:    nil,
			Success: false,
			Message: "Failed! Validation error.",
			Error:   validation.FormatValidationError(err),
		})
	}

	task, err := EditTask(id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ResponseWrapper[Task]{
			Data:    nil,
			Success: false,
			Message: "Failed! Something went wrong.",
		})
	}

	return c.JSON(dto.ResponseWrapper[Task]{
		Data:    task,
		Success: true,
		Message: "Success! task updated.",
	})
}

// DeleteTask godoc
// @Summary Delete task
// @Description Delete task by id
// @Tags task
// @Produce json
// @Success 200 {object} Task
// @Router /task/{id} [delete]
// @Param id path string true "Task ID"
func HandleDeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := DeleteTask(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ResponseWrapper[Task]{
			Data:    nil,
			Success: false,
			Message: "Failed! Something went wrong.",
		})
	}

	return c.JSON(dto.ResponseWrapper[Task]{
		Success: true,
		Message: "Success! task deleted.",
	})
}
