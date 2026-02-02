package belajar_go_lang_gorm

import "time"

type Address struct {
	ID        int64 `gorm:"primary_key;column:id"`
	UserId    string `gorm:"column:user_id"`
	Address   string  `gorm:"column:address"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`

	// implementasi belongs to (one to many)
	// digunakan untuk memudahkan pada saat pengambilan data Address, kita bisa mendapatkan informasi-
	// data user jika di butuhkan
	User User `gorm:"foreignKey:user_id;references:id"`
}

// menentukan nama table
func (w Address) TableName() string {
	return "addresses"
}