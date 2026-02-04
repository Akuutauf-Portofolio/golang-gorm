package belajar_go_lang_gorm

// implementasi auto increment, menambahkan tabel dan model baru
type UserLog struct {
	ID        string `gorm:"primary_key;column:id;autoIncrement"` 
	UserId  string `gorm:"column:user_id"`
	Action  string `gorm:"column:action"`

	// implementasi timestamp tracking
	// mengubah tipe data timestamp dari time.Time menjadi int64(big int)
	// dan menggunakan tag untuk default value create dan update nya menjadi mili
	// ketika menambahkan atau mengupdate data, maka otomatis created_at dan updated_at akan di isi milli
	// dan kita bisa konversi ke platform current millis untuk melihatnya
	CreatedAt int64 `gorm:"column:created_at;autoCreateTime:milli;"`
	UpdatedAt int64 `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
}

// membuat method baru untuk mengganti nama tabel (alias)
func (u *UserLog) TableName() string {
	return "user_logs"
}
