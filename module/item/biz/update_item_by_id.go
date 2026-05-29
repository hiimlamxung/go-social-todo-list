package biz

import (
	"context"
	"g09-social-todo-list/module/item/model"
)

type UpdateItemStorage interface {
	GetItem(ctx context.Context, cond map[string]any) (*model.TodoItem, error)
	UpdateItem(ctx context.Context, cond map[string]any, dataUpdate *model.TodoItemUpdate) error
}

type updateItemBiz struct {
	store UpdateItemStorage
}

func NewUpdateItemBiz(store UpdateItemStorage) *updateItemBiz {
	return &updateItemBiz{store: store}
}

func (biz *updateItemBiz) UpdateItemById(ctx context.Context, id int, dataUpdate *model.TodoItemUpdate) error {
	data, err := biz.store.GetItem(ctx, map[string]any{"id": id})
	if err != nil {
		return err
	}
	if data.Status == "Deleted" {
		return model.ErrorItemDeleted
	}

	if err := biz.store.UpdateItem(ctx, map[string]any{"id": id}, dataUpdate); err != nil {
		return err
	}

	return nil
}
