package middleware

import (
	"g09-social-todo-list/common"

	"github.com/gin-gonic/gin"
)

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Header("Content-Type", "application/json")
				if appErr, ok := err.(*common.AppError); ok {
					c.AbortWithStatusJSON(appErr.StatusCode, appErr)
					return
				}

				appErr := common.ErrInternal(err.(error))
				c.AbortWithStatusJSON(appErr.StatusCode, appErr)

				/* Bật lại panic lần nữa để dẫn tới middleware Recover gốc của Gin ()
				* Mục đích là để hiển thị lại stack trace ở log theo hành vi mặc định của Gin,
				* nếu ko có thì log sẽ ko có stack trace
				 */
				panic(err)
			}
		}()
		c.Next()
	}
}
