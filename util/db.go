package util

import (
	"myapp/config"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/google/uuid"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var Db = newDatabase()

func newDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Dialector{
		DSN:        config.Env.DB_DSN,
		DriverName: "libsql",
	}, &gorm.Config{},
	)

	if err != nil {
		panic("failed to connect database")
	}

	return db
}

type Entity struct {
	ID        uuid.UUID  `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

func (base *Entity) BeforeCreate(tx *gorm.DB) error {
	base.ID = uuid.New()
	return nil
}

func (base *Entity) BeforeSave(tx *gorm.DB) error {
	base.UpdatedAt = time.Now()
	return nil
}
