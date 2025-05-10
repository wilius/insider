package message

import (
	"context"
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"insider/types"
)

type repository struct {
	db *gorm.DB
}

func newRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, item *entity) (*entity, error) {
	result := r.db.WithContext(ctx).
		Clauses(clause.Returning{}).
		Create(&item)

	if result.Error != nil {
		return nil, result.Error
	}

	return item, nil
}

func (r *repository) List(ctx context.Context, filter *Filter) (*types.Pageable, error) {
	var items *[]entity

	dbQuery := r.db.
		WithContext(ctx).
		Order("id desc").
		Limit(filter.CalculateLimit()).
		Offset(filter.CalculateOffset())

	if filter.Query != nil {
		dbQuery = dbQuery.
			Where("name ilike @query", sql.Named("query", fmt.Sprintf("%%%s%%", *filter.Query)))
	}

	result := dbQuery.
		Find(&items)

	if result.Error != nil {
		return nil, result.Error
	}

	return types.MapToPageable(items, &filter.PagedFilter), nil
}

func (r *repository) FetchForSending(ctx context.Context, count uint) (*[]entity, error) {
	var items *[]entity

	result := r.db.
		Raw(`
           UPDATE notifications.message
			SET status = @upToDateStatus,
			    update_date = now()
			WHERE id IN (
				SELECT id 
				FROM notifications.message 
				WHERE status = @currentStatus
				ORDER BY create_date
				LIMIT @fetchCount
			)
			RETURNING *;
        `, map[string]interface{}{
			"upToDateStatus": Sending,
			"currentStatus":  Created,
			"fetchCount":     count,
		}).
		Find(&items)

	if result.Error != nil {
		return nil, result.Error
	}

	return items, nil
}
