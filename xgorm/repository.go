// Package xgorm
/*
 * @Date: 2023-07-20 10:07:02
 * @LastEditTime: 2023-07-20 10:17:25
 * @Description:
 */
package xgorm

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/VNTools1/xutils/xcache"
	gormrepository "github.com/aklinkert/go-gorm-repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type NullType byte

const (
	_ NullType = iota
	// IsNull the same as `is null`
	IsNull
	// IsNotNull the same as `is not null`
	IsNotNull
)

type GormTransactionRepository interface {
	gormrepository.TransactionRepository
	FindWhere(target interface{}, filters map[string]interface{}, preloads ...string) error
	DeleteWhere(target interface{}, filters map[string]interface{}) error
	FindWhereBatch(target interface{}, filters map[string]interface{}, limit, offset int, orderBy string, preloads ...string) error
	FindWhereCount(target interface{}, filters map[string]interface{}) int64
	UpdateWhere(target interface{}, filters map[string]interface{}, updates interface{}, preloads ...string) error
}

type gormRepository struct {
	logger       *logrus.Logger
	db           *gorm.DB
	defaultJoins []string
	debug        bool
	useCache     bool
	cacheTtl     time.Duration
	cachePrefix  string
}

// NewGormRepository returns a new base repository that implements TransactionRepository
func NewGormRepository(db *gorm.DB, logger *logrus.Logger, debug bool, useCache bool, cacheTtl time.Duration, cachePrefix string, defaultJoins ...string) GormTransactionRepository {
	return &gormRepository{
		defaultJoins: defaultJoins,
		logger:       logger,
		db:           db,
		debug:        debug,
		useCache:     useCache,
		cacheTtl:     cacheTtl,
		cachePrefix:  cachePrefix,
	}
}

func (r *gormRepository) DB() *gorm.DB {
	return r.DBWithPreloads(nil)
}

func (r *gormRepository) GetAll(target interface{}, preloads ...string) error {
	r.logger.Debugf("Executing GetAll on %T", target)

	res := r.DBWithPreloads(preloads).
		Unscoped().
		Find(target)

	return r.HandleError(res)
}

func (r *gormRepository) GetBatch(target interface{}, limit, offset int, preloads ...string) error {
	r.logger.Debugf("Executing GetBatch on %T", target)

	res := r.DBWithPreloads(preloads).
		Unscoped().
		Limit(limit).
		Offset(offset).
		Find(target)

	return r.HandleError(res)
}

func (r *gormRepository) UpdateWhere(target interface{}, filters map[string]interface{}, updates interface{}, preloads ...string) error {
	r.logger.Debugf("Executing UpdateWhere on %T with filters = %+v ", target, filters)

	res := r.DBWithPreloads(preloads).
		Model(target).
		Where(filters).
		Updates(updates)

	return r.HandleError(res)
}

func (r *gormRepository) FindWhere(target interface{}, filters map[string]interface{}, preloads ...string) error {
	r.logger.Debugf("Executing FindWhere on %T with filters = %+v ", target, filters)
	cond, vals, err := r.whereBuild(filters)
	if err != nil {
		return err
	}
	res := r.DBWithPreloads(preloads).
		Where(cond, vals...).
		Order("id desc").
		Find(target)

	return r.HandleError(res)
}

func (r *gormRepository) FindWhereBatch(target interface{}, filters map[string]interface{}, limit, offset int, orderBy string, preloads ...string) error {
	r.logger.Debugf("Executing FindWhereBatch on %T with filters = %+v ", target, filters)
	cond, vals, err := r.whereBuild(filters)
	if err != nil {
		return err
	}
	if orderBy == "" {
		orderBy = "id desc"
	}
	res := r.DBWithPreloads(preloads).
		Where(cond, vals...).
		Limit(limit).
		Offset(offset).
		Order(orderBy).
		Find(target)

	return r.HandleError(res)
}

func (r *gormRepository) FindWhereCount(target interface{}, filters map[string]interface{}) int64 {
	r.logger.Debugf("Executing FindWhereCount on %T with filters = %+v ", target, filters)
	var total int64
	cond, vals, err := r.whereBuild(filters)
	if err != nil {
		return 0
	}
	r.DB().Where(cond, vals...).Find(target).Count(&total)
	return total
}

func (r *gormRepository) DeleteWhere(target interface{}, filters map[string]interface{}) error {
	r.logger.Debugf("Executing Delete on %T with filters = %+v ", target, filters)
	cond, vals, err := r.whereBuild(filters)
	if err != nil {
		return err
	}
	res := r.db.Where(cond, vals...).Delete(target)
	return r.HandleError(res)
}

func (r *gormRepository) GetWhere(target interface{}, condition string, preloads ...string) error {
	r.logger.Debugf("Executing GetWhere on %T with %v ", target, condition)

	res := r.DBWithPreloads(preloads).
		Where(condition).
		Find(target)

	return r.HandleError(res)
}

func (r *gormRepository) GetWhereBatch(target interface{}, condition string, limit, offset int, preloads ...string) error {
	r.logger.Debugf("Executing GetWhere on %T with %v ", target, condition)

	res := r.DBWithPreloads(preloads).
		Where(condition).
		Limit(limit).
		Offset(offset).
		Find(target)

	return r.HandleError(res)
}

func (r *gormRepository) GetByField(target interface{}, field string, value interface{}, preloads ...string) error {
	r.logger.Debugf("Executing GetByField on %T with %v = %v", target, field, value)

	res := r.DBWithPreloads(preloads).
		Where(fmt.Sprintf("%v = ?", field), value).
		Find(target)

	return r.HandleError(res)
}

func (r *gormRepository) GetByFields(target interface{}, filters map[string]interface{}, preloads ...string) error {
	r.logger.Debugf("Executing GetByField on %T with filters = %+v", target, filters)

	db := r.DBWithPreloads(preloads)
	for field, value := range filters {
		db = db.Where(fmt.Sprintf("%v = ?", field), value)
	}

	res := db.Find(target)

	return r.HandleError(res)
}

func (r *gormRepository) GetByFieldBatch(target interface{}, field string, value interface{}, limit, offset int, preloads ...string) error {
	r.logger.Debugf("Executing GetByField on %T with %v = %v", target, field, value)

	res := r.DBWithPreloads(preloads).
		Where(fmt.Sprintf("%v = ?", field), value).
		Limit(limit).
		Offset(offset).
		Find(target)

	return r.HandleError(res)
}

func (r *gormRepository) GetByFieldsBatch(target interface{}, filters map[string]interface{}, limit, offset int, preloads ...string) error {
	r.logger.Debugf("Executing GetByField on %T with filters = %+v", target, filters)

	db := r.DBWithPreloads(preloads)
	for field, value := range filters {
		db = db.Where(fmt.Sprintf("%v = ?", field), value)
	}

	res := db.
		Limit(limit).
		Offset(offset).
		Find(target)

	return r.HandleError(res)
}

func (r *gormRepository) GetOneByField(target interface{}, field string, value interface{}, preloads ...string) error {
	r.logger.Debugf("Executing GetOneByField on %T with %v = %v", target, field, value)
	ctx := context.Background()
	if r.useCache {
		keyStr := fmt.Sprintf("%s%T_%s_%v", r.cachePrefix, target, field, value)
		ctx = xcache.NewExpiration(ctx, r.cacheTtl)
		ctx = xcache.NewKey(ctx, keyStr)
	}
	res := r.DBWithPreloads(preloads).WithContext(ctx).
		Where(fmt.Sprintf("%v = ?", field), value).
		First(target)

	return r.HandleOneError(res)
}

func (r *gormRepository) GetOneByFields(target interface{}, filters map[string]interface{}, preloads ...string) error {
	r.logger.Debugf("Executing FindOneByField on %T with filters = %+v", target, filters)
	ctx := context.Background()
	if r.useCache {
		keyStr := fmt.Sprintf("%s%T_%+v", r.cachePrefix, target, filters)
		ctx = xcache.NewExpiration(ctx, r.cacheTtl)
		ctx = xcache.NewKey(ctx, keyStr)
	}
	db := r.DBWithPreloads(preloads).WithContext(ctx)
	for field, value := range filters {
		db = db.Where(fmt.Sprintf("%v = ?", field), value)
	}

	res := db.First(target)
	return r.HandleOneError(res)
}

func (r *gormRepository) GetOneByID(target interface{}, id string, preloads ...string) error {
	r.logger.Debugf("Executing GetOneByID on %T with ID %v", target, id)
	ctx := context.Background()
	if r.useCache {
		keyStr := fmt.Sprintf("%s%T_%s", r.cachePrefix, target, id)
		ctx = xcache.NewExpiration(ctx, r.cacheTtl)
		ctx = xcache.NewKey(ctx, keyStr)
	}
	res := r.DBWithPreloads(preloads).WithContext(ctx).
		Where("id = ?", id).
		First(target)

	return r.HandleOneError(res)
}

func (r *gormRepository) Create(target interface{}) error {
	r.logger.Debugf("Executing Create on %T", target)

	res := r.db.Create(target)
	return r.HandleError(res)
}

func (r *gormRepository) CreateTx(target interface{}, tx *gorm.DB) error {
	r.logger.Debugf("Executing Create on %T", target)

	res := tx.Create(target)
	return r.HandleError(res)
}

func (r *gormRepository) Save(target interface{}) error {
	r.logger.Debugf("Executing Save on %T", target)

	res := r.db.Save(target)
	return r.HandleError(res)
}

func (r *gormRepository) SaveTx(target interface{}, tx *gorm.DB) error {
	r.logger.Debugf("Executing Save on %T", target)

	res := tx.Save(target)
	return r.HandleError(res)
}

func (r *gormRepository) Delete(target interface{}) error {
	r.logger.Debugf("Executing Delete on %T", target)

	res := r.db.Delete(target)
	return r.HandleError(res)
}

func (r *gormRepository) DeleteTx(target interface{}, tx *gorm.DB) error {
	r.logger.Debugf("Executing Delete on %T", target)

	res := tx.Delete(target)
	return r.HandleError(res)
}

func (r *gormRepository) HandleError(res *gorm.DB) error {
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		err := fmt.Errorf("error: %w", res.Error)
		r.logger.Error(err)
		return err
	}

	return nil
}

func (r *gormRepository) HandleOneError(res *gorm.DB) error {
	if err := r.HandleError(res); err != nil {
		return err
	}

	if res.RowsAffected != 1 {
		return gormrepository.ErrNotFound
	}

	return nil
}

func (r *gormRepository) DBWithPreloads(preloads []string) *gorm.DB {
	dbConn := r.db

	for _, join := range r.defaultJoins {
		dbConn = dbConn.Joins(join)
	}

	for _, preload := range preloads {
		dbConn = dbConn.Preload(preload)
	}

	if r.debug {
		dbConn = dbConn.Debug()
	}

	return dbConn
}

// sql build where
func (r *gormRepository) whereBuild(where map[string]interface{}) (whereSQL string, vals []interface{}, err error) {
	for k, v := range where {
		ks := strings.Split(k, " ")
		if len(ks) > 3 {
			return "", nil, fmt.Errorf("error in query condition: %s. ", k)
		}

		if whereSQL != "" {
			if len(ks) == 3 && strings.ToLower(ks[2]) == "or" {
				whereSQL += " OR "
			} else {
				whereSQL += " AND "
			}
		}
		strings.Join(ks, ",")
		switch len(ks) {
		case 1:
			//fmt.Println(reflect.TypeOf(v))
			switch v := v.(type) {
			case NullType:
				if v == IsNotNull {
					whereSQL += fmt.Sprint("`", k, "`", " IS NOT NULL")
				} else {
					whereSQL += fmt.Sprint("`", k, "`", " IS NULL")
				}
			default:
				whereSQL += fmt.Sprint("`", k, "`", "=?")
				vals = append(vals, v)
			}
		case 2, 3:
			k = ks[0]
			switch strings.ToLower(ks[1]) {
			case "=":
				whereSQL += fmt.Sprint("`", k, "`", "=?")
				vals = append(vals, v)
			case ">":
				whereSQL += fmt.Sprint("`", k, "`", ">?")
				vals = append(vals, v)
			case ">=":
				whereSQL += fmt.Sprint("`", k, "`", ">=?")
				vals = append(vals, v)
			case "<":
				whereSQL += fmt.Sprint("`", k, "`", "<?")
				vals = append(vals, v)
			case "<=":
				whereSQL += fmt.Sprint("`", k, "`", "<=?")
				vals = append(vals, v)
			case "!=":
				whereSQL += fmt.Sprint("`", k, "`", "!=?")
				vals = append(vals, v)
			case "<>":
				whereSQL += fmt.Sprint("`", k, "`", "!=?")
				vals = append(vals, v)
			case "in":
				whereSQL += fmt.Sprint("`", k, "`", " in (?)")
				vals = append(vals, v)
			case "like":
				whereSQL += fmt.Sprint("`", k, "`", " like ?")
				vals = append(vals, v)
			}
		}
	}
	return
}
