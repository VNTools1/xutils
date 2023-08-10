/*
 * @Date: 2023-07-20 10:09:47
 * @LastEditTime: 2023-07-20 10:16:49
 * @Description:
 */
package xgorm

import (
	"github.com/VNTools1/xutils/xcache"
	"gorm.io/gorm"
)

func NewGormWithCache(dialector gorm.Dialector, store xcache.Store, opts ...gorm.Option) (db *gorm.DB, err error) {
	db, err = gorm.Open(dialector, opts...)
	if err != nil {
		return
	}
	cacheConfig := &xcache.Config{
		Store:      store,
		Serializer: &xcache.DefaultJSONSerializer{},
	}
	cachePlugin := xcache.New(cacheConfig)
	err = db.Use(cachePlugin)
	return
}
