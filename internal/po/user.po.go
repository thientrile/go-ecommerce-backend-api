package po

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UUID uuid.UUID `gorm:"column:uuid; type: varchar(255);not null; unique;index:idx_uuid"`
	Username string `gorm:"column:username; type: varchar(255);not null; unique;index:idx_username"`
	IsActive bool ` gorm:"column:is_active; type:boolean"`
	Roles [] Role `gorm:"many2many:go_user_roles"`
}
func (r *User) TableName() string{
	return "go_db_user"
}