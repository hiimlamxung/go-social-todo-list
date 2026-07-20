package upload

import (
	"fmt"
	"g09-social-todo-list/common"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Upload(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		fileHeader, err := c.FormFile("file")
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}
		c.SaveUploadedFile(fileHeader, fmt.Sprintf("./static/%d.%s", time.Now().UTC().UnixNano(), fileHeader.Filename))
	}
}
