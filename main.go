package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type tasks struct {
	gorm.Model
	Title string `gorm:"unique_index:idx_name_code"`
	List  string `gorm:"unique_index:idx_name_code"`
}

type report struct {
	Tasks []tasks
	Lists map[string][]tasks
}

type periodic func()

var db *gorm.DB
var mqttClient mqtt.Client

func initDB() {

	mysqlConnection := os.Getenv("DB")

	var err error
	db, err = gorm.Open("mysql", mysqlConnection)

	if err != nil {
		panic(err)
	}

	db.LogMode(true)

	if !db.HasTable(&tasks{}) {
		db.CreateTable(&tasks{})
		db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&tasks{})
	}
}

func initMqtt() {

	broker := os.Getenv("MQTT_BROKER")
	username := os.Getenv("MQTT_USERNAME")
	password := os.Getenv("MQTT_PASSWORD")

	opts := mqtt.NewClientOptions()
	opts = opts.AddBroker(broker).SetUsername(username).SetPassword(password)

	mqttClient = mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}

}

func mqttPublish(topic string, message string) {

	if mqttClient == nil {
		fmt.Println("MQTT Client not initialized")
	}

	if token := mqttClient.Publish(topic, 0, false, message); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}

}

func createReport() report {
	var ts []tasks
	db.Find(&ts)

	lists := map[string][]tasks{}
	rep := report{ts, lists}

	for _, t := range ts {
		rep.Lists[t.List] = append(rep.Lists[t.List], t)
	}

	return rep
}

func sendReport() {

	rep := createReport()

	mqttPublish("todoodle/tasks", strconv.Itoa(len(rep.Tasks)))
	mqttPublish("todoodle/lists", strconv.Itoa(len(rep.Lists)))

	for l, t := range rep.Lists {

		topic := fmt.Sprintf("todoodle/lists/%s", l)
		oldest := fmt.Sprintf("%s/oldest", topic)
		newest := fmt.Sprintf("%s/newest", topic)
		mqttPublish(topic, strconv.Itoa(len(t)))
		if len(t) > 0 {
			mqttPublish(oldest, t[0].Title)
			mqttPublish(newest, t[len(t)-1].Title)
		}
	}
}

func initTicker(fn periodic) {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
			fn()
		}
	}()
}

func setupRoutes(r *gin.Engine) {

	v1 := r.Group("api/v1")
	{
		v1.POST("/tasks", postTask)
		v1.GET("/tasks", getTasks)
		v1.DELETE("/tasks", deleteTask)
	}

}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func main() {

	fmt.Println("Starting!")

	initDB()
	defer db.Close()

	initMqtt()

	r := gin.Default()
	r.Use(cors())

	setupRoutes(r)

	r.Run(":8910")

	initTicker(sendReport)
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

func deleteTask(c *gin.Context) {
	var task tasks
	c.Bind(&task)
	db.Where("list = ? and title = ?", task.List, task.Title).Unscoped().Delete(tasks{})

	c.JSON(200, gin.H{"deleted": task})
}
