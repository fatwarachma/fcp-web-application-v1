package api

import (
	"a21hc3NpZ25tZW50/model"
	"a21hc3NpZ25tZW50/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TaskAPI interface {
	AddTask(c *gin.Context)
	UpdateTask(c *gin.Context)
	DeleteTask(c *gin.Context)
	GetTaskByID(c *gin.Context)
	GetTaskList(c *gin.Context)
	GetTaskListByCategory(c *gin.Context)
}

type taskAPI struct {
	taskService service.TaskService
}

func NewTaskAPI(taskRepo service.TaskService) *taskAPI {
	return &taskAPI{taskRepo}
}

func (t *taskAPI) AddTask(c *gin.Context) {
	var newTask model.Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	err := t.taskService.Store(&newTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "add task success"})
}

func (t *taskAPI) UpdateTask(c *gin.Context) {
	taskID, result := strconv.Atoi(c.Param("id"))
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var makeTaskNew model.Task
	if result := c.ShouldBindJSON(&makeTaskNew); result != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error()})
		return
	}

	var id int
	makeTaskNew.ID = taskID
	result = t.taskService.Update(id, &makeTaskNew)
	if result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "update task success"})
}

func (t *taskAPI) DeleteTask(c *gin.Context) {
	taskID, result := strconv.Atoi(c.Param("id"))
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	result = t.taskService.Delete(taskID)
	if result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "delete task success"})
}

func (t *taskAPI) GetTaskByID(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid task ID"})
		return
	}

	task, err := t.taskService.GetByID(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (t *taskAPI) GetTaskList(c *gin.Context) {
	tasks, result := t.taskService.GetList()
	if result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (t *taskAPI) GetTaskListByCategory(c *gin.Context) {
	category, result := strconv.Atoi(c.Param("id"))
	tasks, result := t.taskService.GetTaskCategory(category)
	if result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
