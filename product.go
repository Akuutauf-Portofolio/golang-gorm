package belajar_go_lang_gorm

import "time"

type Product struct {
	ID        string `gorm:"primary_key;column:id"`
	Name	  string `gorm:"column:name"`
	Price     int64  `gorm:"column:price"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`

	// implementasi relasi many to many
	// menambahkan relasi ke tabel penghubung menuju ke tabel user sebagai many to many

	// many 2 many : menunjukkan tabel penghubung antara product dan user
	// foreignKey:id : menunjukkan id (field primary key) di tabel sekarang (product)
	// joinForeignKey:product_id : menunjukkan foreign key pada tabel penghubung yang menghubungkan ke product
	// references:id : menunjukkan id (field primary key di tabel lain (user)
	// joinReferences:user_id : menunjukkan foreign key pada tabel penghubung yang menghubungkan ke user
	LikedByUsers []User `gorm:"many2many:user_like_product;foreignKey:id;joinForeignKey:product_id;references:id;joinReferences:user_id"`
}

// menentukan nama table
func (w Product) TableName() string {
	return "products"
}