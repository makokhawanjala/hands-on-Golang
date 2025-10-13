package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
}

// Mock data for tasks
var tasks = []Task{
	{ID: "1", Title: "Morning Dev Practice", Description: "Spend 1 hour learning Go/Django", DueDate: time.Now(), Status: "Pending"},
	{ID: "2", Title: "Check Freelance Inbox", Description: "Reply to client emails and academic requests", DueDate: time.Now().AddDate(0, 0, 1), Status: "In Progress"},
	{ID: "3", Title: "Work on Azure Project", Description: "Continue cost optimization / Flask app setup", DueDate: time.Now().AddDate(0, 0, 2), Status: "Pending"},
	{ID: "4", Title: "Family Time (Kids)", Description: "Play and spend time with Aaliyah and Marcus", DueDate: time.Now().AddDate(0, 0, 3), Status: "Completed"},
	{ID: "5", Title: "Farm Check", Description: "Inspect goats and land for daily updates", DueDate: time.Now().AddDate(0, 0, 4), Status: "Pending"},
	{ID: "6", Title: "Content/Portfolio Update", Description: "Push code or update GitHub portfolio", DueDate: time.Now().AddDate(0, 0, 5), Status: "In Progress"},
	{ID: "7", Title: "Freelance Writing Task", Description: "Deliver at least one academic/freelance order", DueDate: time.Now().AddDate(0, 0, 6), Status: "Pending"},
	{ID: "8", Title: "Exercise/Health", Description: "Take a walk or light workout for 30 mins", DueDate: time.Now().AddDate(0, 0, 7), Status: "Completed"},
	{ID: "9", Title: "Evening Tech Study", Description: "Learn about WiFi hacking or Django models", DueDate: time.Now().AddDate(0, 0, 8), Status: "Pending"},
	{ID: "10", Title: "Plan Next Day", Description: "Review goals and prepare for tomorrow", DueDate: time.Now().AddDate(0, 0, 9), Status: "In Progress"},
}

// Function to get all the tasks
func getTasks(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"Tasks": tasks})
}

// Function to get a specific task
func GetTaskByID(ctx *gin.Context) {
	id := ctx.Param("id")

	for _, task := range tasks {
		if task.ID == id {
			ctx.JSON(http.StatusOK, task)
			return
		}
	}
	ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}

// Function to update a specific task (PUT/tasks/:id)
func UpdateTaskById(ctx *gin.Context) {
	id := ctx.Param("id")

	var updatedTask Task

	if err := ctx.ShouldBindJSON(&updatedTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, task := range tasks {
		if task.ID == id {
			// Update only the specified fields
			if updatedTask.Title != "" {
				tasks[i].Title = updatedTask.Title
			}
			if updatedTask.Description != "" {
				tasks[i].Description = updatedTask.Description
			}
			ctx.JSON(http.StatusOK, gin.H{"message": "Task updated"})
			return
		}
	}
	ctx.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
}

func getPing(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func main() {
	router := gin.Default()
	router.GET("/ping", getPing)
	router.GET("/tasks", getTasks)
	router.GET("/tasks/:id", GetTaskByID)
	router.PUT("/tasks/:id", UpdateTaskById)
	router.Run()
}
