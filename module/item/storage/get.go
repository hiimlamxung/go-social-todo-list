package storage

import (
	"context"
	"g09-social-todo-list/module/item/model"
)

func (s *sqlStore) GetItem(ctx context.Context, cond map[string]any) (*model.TodoItem, error) {
	var dataModel model.TodoItem
	if err := s.db.Where(cond).First(&dataModel).Error; err != nil {
		return nil, err
	}
	return &dataModel, nil
}
