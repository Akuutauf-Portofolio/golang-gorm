package belajar_go_lang_gorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// implementasi database connection
// membuat function untuk koneksi ke database
func OpenConnection() *gorm.DB {
	// membuat destinasi untuk database yang dituju
	dialect := mysql.Open("root:@tcp(localhost:3306)/belajar_golang_gorm?charset=utf8mb4&parseTime=True&loc=Local")
	db, err := gorm.Open(dialect, &gorm.Config{})

	// mengecek error
	if err != nil {
		panic(err)
	}

	return db
}

// membuat variabel global
var db = OpenConnection()

// membuat kode uji untuk menguji konek database
func TestOpenConnection(t *testing.T) {
	// melakukan perbandingan dengan assert untuk mengecek apakah koneksi ditemukan atau tidak
	assert.NotNil(t, db)
}

// implementasi raw sql : execute sql
func TestExecuteSQL(t *testing.T) {
	// untuk memanipulasi data (insert, update, delete) gunakan function Exec pada gorm.DB
	err := db.Exec("insert into sample(id, name) values (?, ?)", "1", "Taufik").Error
	assert.Nil(t, err) // memastikan tidak ada error pada query

	err = db.Exec("insert into sample(id, name) values (?, ?)", "2", "Ilham").Error
	assert.Nil(t, err) // memastikan tidak ada error pada query 

	err = db.Exec("insert into sample(id, name) values (?, ?)", "3", "Dimas").Error
	assert.Nil(t, err) // memastikan tidak ada error pada query
}

// membuat struct untuk data samples
type Sample struct {
	Id string 
	Name string
}

// implementasi raw sql : query sql
func TestRawSQL(t *testing.T) {
	// mengambil sebuah data dari tabel sample
	// membuat variabel baru untuk menampung sebuah data sample
	var sample Sample

	// melakukan raw sql (untuk sebuah data)
	// function scan digunakan untuk mengirimkan data yang diambil ke bentuk struct
	err := db.Raw("select id, name from sample where id = ?", "1").Scan(&sample).Error
	assert.Nil(t, err)
	assert.Equal(t, "Taufik", sample.Name) // membandingkan isi data pada kolom name

	// mengambil lebih dari satu data dari tabel sample
	// membuat variabel slice baru untuk menampung kumpulan data sample
	var samples []Sample

	// melakukan raw sql (untuk lebih dari satu data)
	err = db.Raw("select id, name from sample").Scan(&samples).Error
	assert.Nil(t, err)
	assert.Equal(t, 3, len(samples)) // membandingkan jumlah data pada tabel sample
}

// implementasi sql row dan sql rows
func TestSqlRow(t *testing.T) {
	// melakukan select dengan method Raw()
	// method Rows() mengembalikan baris data (row) dan error
	rows, err := db.Raw("select id, name from sample").Rows()
	assert.Nil(t, err)

	// jangan lupa menututp rows jika sudah selesai digunakan, agar tidak memory leak
	defer rows.Close()

	// menyiapkan variabel slice kosong
	var samples []Sample

	// jika kita ingin menampilkan data rows bisa gunakan iterasi
	for rows.Next() {
		// menyiapkan data tiap kolom
		var id, name string 

		// mengambil data kolom untuk setiap baris (gunakan pointer)
		// urutan scan disesuaikan dengan query pada Raw() diatas
		err := rows.Scan(&id, &name)
		assert.Nil(t, err) // pastikan tidak ada error setiap pengambilan data per barisnya
		
		// menambahkan data ke variabel samples
		samples = append(samples, Sample{
			Id: id,
			Name: name,
		})
	}

	assert.Equal(t, 3, len(samples)) // membandingkan jumlah data pada table sample
}

// implementasi scan rows
func TestScanRow(t *testing.T) {
	// melakukan select dengan method Raw()
	// method Rows() mengembalikan baris data (row) dan error
	rows, err := db.Raw("select id, name from sample").Rows()
	assert.Nil(t, err)

	// jangan lupa menututp rows jika sudah selesai digunakan, agar tidak memory leak
	defer rows.Close()

	// menyiapkan variabel slice kosong
	var samples []Sample

	// jika kita ingin menampilkan data rows bisa gunakan iterasi
	for rows.Next() {
		// solusi lebih cepat dibanding manual loop data rows
		// dengan menggunakan method ScanRows
		// row di setiap iterasi akan diambil dan dimasukkan ke dalam samples
		err := db.ScanRows(rows, &samples)
		assert.Nil(t, err)
	}

	assert.Equal(t, 3, len(samples)) // membandingkan jumlah data pada table sample
}