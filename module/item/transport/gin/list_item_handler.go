package gin

import (
	"g09-social-todo-list/common"
	"g09-social-todo-list/module/item/biz"
	"g09-social-todo-list/module/item/model"
	"g09-social-todo-list/module/item/storage"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListItem(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var queryString struct {
			common.Paging
			model.Filter
		}
		if err := c.ShouldBind(&queryString); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		queryString.Paging.Process()
		store := storage.NewSQLStore(db)      // tạo layer store
		business := biz.NewListItemBiz(store) // tạo layer business

		result, err := business.ListItem(c.Request.Context(), &queryString.Filter, &queryString.Paging)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, common.NewSuccessResponse(result, queryString.Paging, queryString.Filter))
	}
}
