package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TodoItem struct {
	Id          int        `json:"id" gorm:"column:id"`
	Title       string     `json:"title" gorm:"column:title"`
	Description string     `json:"description" gorm:"column:description"`
	Status      string     `json:"status" gorm:"column:status"`
	CreatedAt   *time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (TodoItem) TableName() string {
	return "todo_items"
}

type TodoItemCreation struct {
	Title       string  `json:"title" gorm:"column:title" binding:"required"`
	Description *string `json:"description" gorm:"column:description"`
}

func (TodoItemCreation) TableName() string {
	return TodoItem{}.TableName()
}

type TodoItemUpdate struct {
	Title       string  `json:"title" gorm:"column:title"`
	Description *string `json:"description" gorm:"column:description"`
	Status      string  `json:"status" gorm:"column:status"`
}

func (TodoItemUpdate) TableName() string {
	return TodoItem{}.TableName()
}

type Paging struct {
	Page  int   `json:"page" form:"page"`
	Limit int   `json:"limit" form:"limit"`
	Total int64 `json:"total" form:"-"`
}

func (p *Paging) Process() {
	if p.Page <= 1 {
		p.Page = 1
	}
	if p.Limit <= 1 || p.Limit > 200 {
		p.Limit = 10
	}
}

func main() {
	godotenv.Load() // Load environment variables from .env file

	// Connect to database
	dsn := os.Getenv("DB_DSN")
	log.Println("Database connection string:", dsn)
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
			items.POST("/", CreateItem(db))
			items.GET("/", ListItem(db))
			items.GET("/:id", GetItem(db))
			items.PUT("/:id", UpdateItem(db))
			items.DELETE("/:id", DeleteItem(db))
		}
	}

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.Run(":3000") // listens on 0.0.0.0:8080 by default
}

func CreateItem(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var itemData TodoItemCreation

		if err := c.ShouldBind(&itemData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := db.Create(&itemData).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Item created successfully",
			"data":    itemData,
		})
	}
}

func GetItem(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var itemData TodoItem

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := db.Where("id = ?", id).First(&itemData).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": itemData,
		})
	}
}

func UpdateItem(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		var updateData TodoItemUpdate
		if err := c.ShouldBind(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := db.Where("id = ?", id).Updates(&updateData).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": true,
		})
	}
}

func DeleteItem(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		deletedStatus := "Deleted"
		if err := db.Where("id = ?", id).Updates(&TodoItemUpdate{Status: deletedStatus}).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": true,
		})
	}
}

func ListItem(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var paging Paging

		if err := c.ShouldBind(&paging); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		paging.Process()
		var result []TodoItem

		db = db.Table(TodoItem{}.TableName()).Where("status != ?", "Deleted")
		// Count total items
		db.Count(&paging.Total)
		// Get items current page
		if err := db.Select("*").Offset((paging.Page - 1) * paging.Limit).
			Limit(paging.Limit).
			Find(&result).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data":   result,
			"paging": paging,
		})
	}
}
