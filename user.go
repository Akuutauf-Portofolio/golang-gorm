package belajar_go_lang_gorm

import "time"

// gorm otomatis mengenali tabel dengan nama 'users'
// sehingga contoh kalau nama tabel / struct User => 'users' dan atau OrderDetail => 'order_details'
type User struct {
	ID        string `gorm:"primary_key;column:id;<-:create"` // kolom id datanya hanya boleh dicreate saja, tidak boleh di update
	Password  string `gorm:"column:password"`

	// field name sebagai embedded struct Name
	Name Name `gorm:"embedded"` // sebagai embedded, maka secara otomatis kolom di struct Name akan ditambahkan secara embedded disini

	// tidak perlu menggunakan autoCreateTime pun gorm sudah setting kolom ini sebagai created_at-
	// karena sudah diberikan nama kolom nya adalah 'CreatedAt'
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;<-:create"` // kolom created_at datanya hanya boleh dicreate saja, tidak boleh di update

	// tidak perlu menggunakan autoUpdateTime pun gorm sudah setting kolom ini sebagai updated_at-
	// karena sudah diberikan nama kolom nya adalah 'UpdatedAt' (termasuk ke attribute updated_at juga)
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	Information string `gorm:"-"` // di abaikan / tidak ada kolom nya di database
}

// membuat method baru untuk mengganti nama tabel (alias)
func (u *User) TableName() string {
	return "users"
}

// implementasi embedded struct
// membuat struct baru untuk embedded struct
type Name struct {
	FirstName string `gorm:"first_name"`
	MiddleName string `gorm:"middle_name"`
	LastName string `gorm:"last_name"`
}