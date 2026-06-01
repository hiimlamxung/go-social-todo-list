package storage

import (
	"context"
	"g09-social-todo-list/common"
	"g09-social-todo-list/module/item/model"
)

func (s *sqlStore) ListItem(
	ctx context.Context,
	filter *model.Filter,
	paging *common.Paging, moreKeys ...string,
) ([]model.TodoItem, error) {
	var result []model.TodoItem
	db := s.db.Table(model.TodoItem{}.TableName())
	if f := filter; f != nil {
		if v := f.Status; v != "" {
			db = db.Where("status = ?", v)
		} else {
			db = db.Where("status != ?", "Deleted")
		}
	}
	// Tính total items
	if err := db.Count(&paging.Total).Error; err != nil {
		return nil, err
	}

	// Lấy items current page
	if err := db.Offset((paging.Page - 1) * paging.Limit).Order("created_at DESC").Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
