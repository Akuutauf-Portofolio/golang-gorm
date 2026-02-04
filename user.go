package belajar_go_lang_gorm

import (
	"time"

	"gorm.io/gorm"
)

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

	// implementasi one to one (has one)
	Wallet Wallet `gorm:"foreignKey:user_id;references:id"`

	// implementasi one to many (has many)
	Addresses []Address `gorm:"foreignKey:user_id;references:id"`

	// implementasi relasi many to many
	// menambahkan relasi ke tabel penghubung menuju ke tabel product sebagai many to many
	
	// many 2 many : menunjukkan tabel penghubung antara user dan product
	// foreignKey:id : menunjukkan id (field primary key) di tabel sekarang (user)
	// joinForeignKey:user_id : menunjukkan foreign key pada tabel penghubung yang menghubungkan ke user
	// references:id : menunjukkan id (field primary key di tabel lain (product)
	// joinReferences:product_id : menunjukkan foreign key pada tabel penghubung yang menghubungkan ke product
	LikeProducts []Product `gorm:"many2many:user_like_product;foreignKey:id;joinForeignKey:user_id;references:id;joinReferences:product_id"`
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

// implementasi hook - untuk Before Create (operasi create/insert)
func (u *User) BeforeCreate(db *gorm.DB) error {
	// didalam function ini akan menjalankan operasi kustom
	// jika terjadi error maka akan dibatalkan, sesuai dengan konsep hook
	
	// mengecek jika data user yang dibuat id nya menggunakan " "
	if u.ID == "" {
		// contoh mengubah id nya menjadi kustom
		// mengatur format waktu nya (tahun-bulan-tanggal-jam-menit-detik)
		u.ID = "user-" + time.Now().Format("20060102150405") 
	}
	
	return nil
}