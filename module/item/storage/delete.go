package storage

import (
	"context"
	"g09-social-todo-list/module/item/model"
)

func (s *sqlStore) DeleteItem(ctx context.Context, cond map[string]any) error {
	DeleteStatus := "Deleted"
	if err := s.db.Table(model.TodoItem{}.
		TableName()).
		Where(cond).
		Updates(model.TodoItemUpdate{Status: DeleteStatus}).Error; err != nil {
		return err
	}
	return nil
}
