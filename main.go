package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type tasks struct {
	gorm.Model
	Title string `gorm:"unique_index:idx_name_code"`
	List  string `gorm:"unique_index:idx_name_code"`
}

var db *gorm.DB

func initDB() {

	db, err := gorm.Open("mysql", "todos:todos@tcp(192.168.1.200:3306)/todos?charset=utf8mb4&parseTime=True&loc=Local")
	db.LogMode(true)

	if err != nil {
		panic(err)
	}

	if !db.HasTable(&tasks{}) {
		db.CreateTable(&tasks{})
		db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&tasks{})
	}
}

func setupRoutes(r *gin.Engine) {

	v1 := r.Group("api/v1")
	{
		v1.POST("/tasks", postTask)
		v1.GET("/tasks", getTasks)
		v1.GET("/tasks/:list", getListTasks)
		v1.DELETE("/tasks", deleteTask)
	}

}

func main() {
	initDB()
	defer db.Close()

	r := gin.Default()

	setupRoutes(r)

	r.Run(":8910")
}

func postTask(c *gin.Context) {
	var task tasks
	c.Bind(&task)

	if task.Title != "" && task.List != "" {
		err := db.Create(&task)

		if err.Error == nil {
			c.JSON(201, gin.H{"created": task})
		} else {
			c.JSON(409, gin.H{"already exists": task})
		}

	} else {
		c.JSON(422, gin.H{"error": "Fields are empty"})
	}
}

func getTasks(c *gin.Context) {
	var tasks []tasks
	db.Find(&tasks)
	c.JSON(200, tasks)
}

func getListTasks(c *gin.Context) {
	list := c.Params.ByName("list")

	var tasks []tasks
	db.Where("list = ?", list).Find(&tasks)

	c.JSON(200, tasks)
}

func deleteTask(c *gin.Context) {
	var task tasks
	c.Bind(&task)
	db.Where("list = ? and title = ?", task.List, task.Title).Unscoped().Delete(tasks{})

	c.JSON(201, gin.H{"deleted": task})
}
