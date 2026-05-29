package model

import (
	"errors"
	"g09-social-todo-list/common"
	"strings"
)

var (
	ErrorTitleCannotBeEmpty = errors.New("title cannot be empty")
	ErrorItemDeleted = errors.New("item is deleted")
)

type TodoItem struct {
	common.SQLModel
	Title       string `json:"title" gorm:"column:title"`
	Description string `json:"description" gorm:"column:description"`
	Status      string `json:"status" gorm:"column:status"`
}

func (TodoItem) TableName() string {
	return "todo_items"
}

type TodoItemCreation struct {
	Title       string  `json:"title" gorm:"column:title" binding:"required"`
	Description *string `json:"description" gorm:"column:description"`
}

func (i *TodoItemCreation) Validate() error {
	i.Title = strings.TrimSpace(i.Title)
	if i.Title == "" {
		return ErrorTitleCannotBeEmpty
	}
	return nil
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
