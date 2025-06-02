package po

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	ID       int64  `gorm:"column:id; type: int;not null;primaryKey;autoIncrement;comment: 'Role ID'"`
	RoleName string `gorm:"column:role_name; type: varchar(255);not null; unique;index:idx_role_name;comment: 'Role Name'"`
	RoleNote string `gorm:"column:role_note:type:text;"`
}

func (r *Role) TableName() string {
	return "go_db_role"
}
