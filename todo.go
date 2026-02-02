package belajar_go_lang_gorm

import (
	"gorm.io/gorm"
)

// implementasi model struct
type Todo struct {
	// menggunakan struct bawaan GORM (Model struct), sehingga tidak perlu mendefinisikan-
	// id, created at, updated at, dan delete at
	// id sudah menggunakan default auto increment, namun jika tabel yang kita buat tidak-
	// mengimplementasikan auto increment, maka bisa mendefinisikan field satu persatu
	gorm.Model
	UserId  string `gorm:"column:user_id"`
	Title  string `gorm:"column:title"`
	Description  string `gorm:"column:description"`
}

// // membuat struct dengan cara yang normal
// type Todo struct {
// 	ID        int64 `gorm:"primary_key;column:id;autoIncrement"`
// 	UserId  string `gorm:"column:user_id"`
// 	Title  string `gorm:"column:title"`
// 	Description  string `gorm:"column:description"`
// 	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"` 
// 	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`

// 	// implementasi soft delete
// 	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"` 
// }

// membuat method baru untuk mengganti nama tabel (alias)
func (t *Todo) TableName() string {
	return "todos"
}
