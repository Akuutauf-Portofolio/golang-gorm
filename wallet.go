package belajar_go_lang_gorm

import "time"

type Wallet struct {
	ID        string `gorm:"primary_key;column:id"`
	UserId    string `gorm:"column:user_id"`
	Balance   int64  `gorm:"column:balance"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`

	
	// implementasi belongs to (one to one)
	// digunakan untuk memudahkan pada saat pengambilan data Wallet, kita bisa mendapatkan informasi-
	// data user jika di butuhkan
	User *User `gorm:"foreignKey:user_id;references:id"`
}

// menentukan nama table
func (w Wallet) TableName() string {
	return "wallets"
}