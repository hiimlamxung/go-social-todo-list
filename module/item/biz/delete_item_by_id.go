package biz

import (
	"context"
	"g09-social-todo-list/common"
	"g09-social-todo-list/module/item/model"
)

type DeleteItemStorage interface {
	GetItem(ctx context.Context, cond map[string]any) (*model.TodoItem, error)
	DeleteItem(ctx context.Context, cond map[string]any) error
}

type deleteItemBiz struct {
	store DeleteItemStorage
}

func NewDeleteItemBiz(store DeleteItemStorage) *deleteItemBiz {
	return &deleteItemBiz{store: store}
}

func (biz *deleteItemBiz) DeleteItemById(ctx context.Context, id int) error {
	data, err := biz.store.GetItem(ctx, map[string]any{"id": id})
	if err != nil {
		return common.ErrorCannotGetEntity(model.EntityName, err)
	}
	if data.Status == "Deleted" {
		return model.ErrorItemDeleted
	}

	if err := biz.store.DeleteItem(ctx, map[string]any{"id": id}); err != nil {
		return common.ErrorCannotDeleteEntity(model.EntityName, err)
	}

	return nil
}
