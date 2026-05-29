package biz

import (
	"context"
	"g09-social-todo-list/module/item/model"
)

type CreateNewItemStore interface {
	CreateItem(ctx context.Context, data *model.TodoItemCreation) error
}

type createNewItemBiz struct {
	store CreateNewItemStore
}

func NewCreateNewItemBiz(store CreateNewItemStore) *createNewItemBiz {
	return &createNewItemBiz{store: store}
}

func (biz *createNewItemBiz) CreateNewItem(ctx context.Context, data *model.TodoItemCreation) error {
	if err := data.Validate(); err != nil {
		return err
	}
	if err := biz.store.CreateItem(ctx, data); err != nil {
		return err
	}
	return nil
}
