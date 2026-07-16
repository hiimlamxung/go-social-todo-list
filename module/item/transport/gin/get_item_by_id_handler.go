package gin

import (
	"g09-social-todo-list/common"
	"g09-social-todo-list/module/item/biz"
	"g09-social-todo-list/module/item/storage"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetItem(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		store := storage.NewSQLStore(db)     // tạo layer store
		business := biz.NewGetItemBiz(store) // tạo layer business

		data, err := business.GetItemById(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data))
	}
}
