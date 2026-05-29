package main

import (
	"log"
	"net/http"
	"os"
	"time"

	ginitem "g09-social-todo-list/module/item/transport/gin"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load() // Load environment variables from .env file

	// Connect to database
	dsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	db = db.Debug()
	log.Println("Database connected successfully", db)

	// Lấy *sql.DB nằm dưới GORM để cấu hình connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalln(err)
	}
	sqlDB.SetMaxOpenConns(100)                 // tối đa 25 kết nối đồng thời
	sqlDB.SetMaxIdleConns(25)                  // giữ sẵn 25 kết nối rảnh
	sqlDB.SetConnMaxLifetime(5 * time.Minute)  // kết nối sống tối đa 5 phút
	sqlDB.SetConnMaxIdleTime(10 * time.Minute) // rảnh quá 10 phút thì đóng

	// ====================================================
	router := gin.Default()

	v1 := router.Group("/v1")
	{
		items := v1.Group("/items")
		{
			items.POST("/", ginitem.CreateItem(db))
			items.GET("/", ginitem.ListItem(db))
			items.GET("/:id", ginitem.GetItem(db))
			items.PUT("/:id", ginitem.UpdateItem(db))
			items.DELETE("/:id", ginitem.DeleteItem(db))
		}
	}

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.Run(":3000") // listens on 0.0.0.0:8080 by default
}
