package gomini

import (
	"log"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func initDB() {
	var err error
	if DB, err = gorm.Open(postgres.Open(cfg.Dsn), &gorm.Config{}); err != nil {
		log.Println("database error:", err)
	}
}
func AutoMigrate(obj ...interface{}) {
	DB.AutoMigrate(obj...)
}
type validate struct {
	Err error
}

func Validate() *validate {
	return &validate{}
}
func (v *validate) Req(field string, val any) *validate {
	if val == "" {
		v.Err = errors.New(fmt.Sprintf("必须录入:%s", field))
	}
	return v
}
func (v *validate) Exists(field string, val any) *validate {
	result := DB.Take(val)
	if result.RowsAffected > 0 {
		v.Err = errors.New(fmt.Sprintf("不能重复录入:%s", field))
	}
	return v
}
