package db

import "gorm.io/gorm"

type MafiaGormInterface interface {
	Where(query interface{}, args ...interface{}) (tx *gorm.DB)
	Save(value interface{}) (tx *gorm.DB)
	First(dest interface{}, conds ...interface{}) (tx *gorm.DB)
}
