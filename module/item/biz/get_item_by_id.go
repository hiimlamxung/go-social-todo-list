package biz

import (
	"context"
	"g09-social-todo-list/module/item/model"
)

type GetItemStorage interface {
	GetItem(ctx context.Context, cond map[string]any) (*model.TodoItem, error)
}

type getItemBiz struct {
	store GetItemStorage
}

func NewGetItemBiz(store GetItemStorage) *getItemBiz {
	return &getItemBiz{store: store}
}

func (biz *getItemBiz) GetItemById(ctx context.Context, id int) (*model.TodoItem, error) {
	data, err := biz.store.GetItem(ctx, map[string]any{"id": id})
	if err != nil {
		return nil, err
	}
	return data, nil
}
