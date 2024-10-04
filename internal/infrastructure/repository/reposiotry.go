package repository

import (
	"context"
	"time"

	"github.com/AsaHero/just-ask-bot/pkg/database/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CtxGorm string

const CtxGormKey CtxGorm = "ctx-gorm"

func FromContext(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	db, ok := ctx.Value(CtxGormKey).(*gorm.DB)
	if ok {
		return db
	}
	return defaultDB.WithContext(ctx)
}

// BaseRepository interface now includes generic type T which must satisfy Entity interface
type BaseRepository[T any] interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	FindAll(ctx context.Context, limit, page uint64, orderBy string, filter map[string]any) (uint64, []*T, error)
	FindOne(ctx context.Context, filter map[string]any) (*T, error)
	Create(ctx context.Context, e *T) error
	Update(ctx context.Context, e *T) error
	UpdateDataWhere(ctx context.Context, data map[string]any, filter map[string]any) error
	Upsert(ctx context.Context, clauseName string, toUpdateColumns []string, e *T) error
	BatchCreate(ctx context.Context, entities []*T) error
}

type baseRepository[T any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) BaseRepository[T] {
	return &baseRepository[T]{
		db: db,
	}
}

func (r *baseRepository[T]) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		ctx := context.WithValue(ctx, CtxGormKey, tx)
		return fn(ctx)
	})
}

// Implement additional methods using generic type T
func (r *baseRepository[T]) FindAll(ctx context.Context, limit, page uint64, orderBy string, filter map[string]any) (uint64, []*T, error) {
	var model T
	var results []*T
	db := FromContext(ctx, r.db)

	// Calculate offset
	if limit != 0 {
		offset := (page - 1) * limit
		db = db.Offset(int(offset)).Limit(int(limit))
	}

	if orderBy != "" {
		db = db.Order(orderBy)
	}

	// Apply filtering, pagination, and find operation
	for key, value := range filter {
		switch v := value.(type) {
		case []time.Time: // Handle date range
			if len(v) == 2 {
				db = db.Where(key+" BETWEEN ? AND ?", v[0], v[1])
			}
		case time.Time: // Handle exact date
			db = db.Where(key+" = ?", v)
		case postgres.TimeCondition: // Handle complex date conditions
			for cond, val := range v {
				db = db.Where(key+" "+string(cond)+" ?", val)
			}
		default: // Handle all other data types
			db = db.Where(key, value)
		}
	}
	result := db.Find(&results)
	if result.Error != nil {
		return 0, nil, postgres.Error(result.Error, "FindAll", &model)
	}

	// Clone the DB session for count to avoid reusing modified `db`
	countDB := FromContext(ctx, r.db)

	// Reapply filtering for count
	for key, value := range filter {
		switch v := value.(type) {
		case []time.Time:
			if len(v) == 2 {
				countDB = countDB.Where(key+" BETWEEN ? AND ?", v[0], v[1])
			}
		case time.Time:
			countDB = countDB.Where(key+" = ?", v)
		case postgres.TimeCondition:
			for cond, val := range v {
				countDB = countDB.Where(key+" "+string(cond)+" ?", val)
			}
		default:
			countDB = countDB.Where(key, value)
		}
	}

	// Count total records matching the filter
	var total int64
	countDB.Model(&model).Count(&total)

	return uint64(total), results, nil
}
func (r *baseRepository[T]) FindOne(ctx context.Context, filter map[string]any) (*T, error) {
	var result T
	db := FromContext(ctx, r.db)
	if err := db.Where(filter).First(&result).Error; err != nil {
		return nil, postgres.Error(err, "FindOne", &result)
	}
	return &result, nil
}

func (r *baseRepository[T]) Create(ctx context.Context, e *T) error {
	db := FromContext(ctx, r.db)
	err := db.Model(e).Create(e).Error
	if err != nil {
		return postgres.Error(err, "Create", e)
	}

	return nil
}

func (r *baseRepository[T]) Update(ctx context.Context, e *T) error {
	db := FromContext(ctx, r.db)
	err := db.Save(e).Error
	if err != nil {
		return postgres.Error(err, "Create", e)
	}

	return nil
}

func (r *baseRepository[T]) UpdateDataWhere(ctx context.Context, data map[string]any, filter map[string]any) error {
	db := FromContext(ctx, r.db)

	var model *T

	err := db.Model(model).Where(filter).Updates(data).Error
	if err != nil {
		return postgres.Error(err, "UpdateDataWhere", model)
	}

	return nil
}

func (r *baseRepository[T]) Upsert(ctx context.Context, clauseName string, toUpdateColumns []string, e *T) error {
	query := FromContext(ctx, r.db)

	if err := query.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: clauseName}},
		DoUpdates: clause.AssignmentColumns(toUpdateColumns),
	}).Create(e).Error; err != nil {
		return postgres.Error(err, "Upsert", e)
	}

	return nil
}

// BatchCreate inserts multiple records in a single database transaction.
func (r *baseRepository[T]) BatchCreate(ctx context.Context, entities []*T) error {
	db := FromContext(ctx, r.db) // Get the database context

	var model *T

	// Create records in a transaction
	result := db.Create(entities) // Using Gorm's Create to handle batch insert
	if result.Error != nil {
		return postgres.Error(result.Error, "BatchCreate", model) // Return an error with context
	}

	return nil
}
