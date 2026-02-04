package belajar_go_lang_gorm

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

// implementasi database connection
// membuat function untuk koneksi ke database
func OpenConnection() *gorm.DB {
	// membuat destinasi untuk database yang dituju
	dialect := mysql.Open("root:@tcp(localhost:3306)/belajar_golang_gorm?charset=utf8mb4&parseTime=True&loc=Local")
	db, err := gorm.Open(dialect, &gorm.Config{
		// implementasi logger
		// menambahkan logger untuk memunculkan informasi log query sql
		Logger: logger.Default.LogMode(logger.Info),

		// implementasi performance
		// tips 1 : matikan auto transaction
		SkipDefaultTransaction: true,

		// tips 2 : cache prepared statement
		PrepareStmt: true,
	})

	// mengecek error
	if err != nil {
		panic(err)
	}

	// implementasi connection pool
	sqlDB, err := db.DB()

	// mengecek error
	if err != nil {
		panic(err)
	}

	// mengatur connection pool
	// Batas maksimal koneksi ke database yang boleh aktif bersamaan.
	// Kalau sudah 100 koneksi dipakai:
	// - request berikutnya akan menunggu
	// - bukan bikin koneksi baru
	sqlDB.SetMaxOpenConns(100)

	// Jumlah koneksi yang disimpan dalam kondisi siap pakai (nganggur).
	// Saat ada request baru:
	// - ambil dari sini dulu
	// - jadi lebih cepat daripada buat koneksi baru
	sqlDB.SetMaxIdleConns(10)

	// Umur maksimal sebuah koneksi.
	// Setelah 30 menit:
	// - koneksi dipensiunkan
	// - diganti koneksi baru
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	// Batas waktu koneksi boleh menganggur.
	// Kalau 5 menit tidak dipakai:
	// - koneksi akan ditutup
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

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

// implementasi create
// membuat pengujian untuk membuat user baru
func TestCreateUser(t *testing.T) {
	// membuat objek user baru dari struct user dan name
	user := User {
		ID: "1",
		Password: "rahasia",
		Name: Name{
			FirstName: "Taufik",
			MiddleName: "H",
			LastName: "Hidayat",
		},
		Information: "data ini akan diabaikan",
		// untuk kolom created at dan updated at akan secara otomatis dibuat kan oleh gorm
	}

	// melakukan insert data ke database menggunakan method Insert()
	response := db.Create(&user) // data yang dikirimkan berupa pointer

	// melakukan assert dari hasil response
	assert.Nil(t, response.Error)

	// memastikan data yang ditambahkan ke database
	assert.Equal(t, int64(1), response.RowsAffected)
}

// implementasi batch insert (create)
func TestBatchInsert(t *testing.T) {
	// menyiapkan tempat data user
	var users []User

	// melakukan perulangan, dimulai dari data id ke 2 (karena data dengan id 1 sudah ada)
	for i := 2; i < 10; i++ {
		// melakukan append ke slices (menambahkan data baru ke slice)
		users = append(users, User{
			ID: strconv.Itoa(i),
			Password: "rahasia",
			Name: Name{
				FirstName: "User " + strconv.Itoa(i),
			},
		})
	}

	// setelah ditambahkan ke slice
	// lakukan batch insert ke database
	result := db.Create(&users)

	// melakukan perbandingan dengan assert
	assert.Nil(t, result.Error)
	assert.Equal(t, int64(8), result.RowsAffected)
}

// implementasi transaction (success)
func TestTransactions(t *testing.T) {
	// membuat transaksi baru
	// ketika membuat transaction, kita tidak perlu mendefinisikan begin dan commit
	// method transaction juga membutuhkan parameter function callback
	err := db.Transaction(func(tx *gorm.DB) error {
		// menambahkan data baru
		err := tx.Create(&User{ID:"10", Password: "rahasia", Name: Name{FirstName: "User 10"}}).Error

		// mengecek jika terjadi error pada saat insert maka return error
		if err != nil {
			return err
		}
		
		// menambahkan data baru
		err = tx.Create(&User{ID:"11", Password: "rahasia", Name: Name{FirstName: "User 11"}}).Error

		// mengecek jika terjadi error pada saat insert maka return error
		if err != nil {
			return err
		}
		
		// menambahkan data baru
		err = tx.Create(&User{ID:"12", Password: "rahasia", Name: Name{FirstName: "User 12"}}).Error

		// mengecek jika terjadi error pada saat insert maka return error
		if err != nil {
			return err
		}

		return nil
	})

	// memastikan transaction berhasil tanpa error
	assert.Nil(t, err)
}

// implementasi transaction (error)
func TestTransactionsError(t *testing.T) {
	// membuat transaksi baru
	// ketika membuat transaction, kita tidak perlu mendefinisikan begin dan commit
	// method transaction juga membutuhkan parameter function callback
	err := db.Transaction(func(tx *gorm.DB) error {
		// menambahkan data baru
		err := tx.Create(&User{ID:"13", Password: "rahasia", Name: Name{FirstName: "User 13"}}).Error

		// mengecek jika terjadi error pada saat insert maka return error
		if err != nil {
			return err
		}
		
		// menambahkan data baru
		err = tx.Create(&User{ID:"11", Password: "rahasia", Name: Name{FirstName: "User 11"}}).Error

		// mengecek jika terjadi error pada saat insert maka return error
		if err != nil {
			return err
		}

		// jika terdapat error pada saat insert data diatas, maka transaction akan di rollback
		return nil
	})

	// memastikan transaction error (karena data dengan id 11 sudah ada)
	assert.NotNil(t, err)
}

// implementasi transaction (manual dan sukses)
func TestManualTransactionSuccess(t *testing.T) {
	// membuat transaksi manual baru
	// ketika membuat transaction manual, kita perlu mendefinisikan begin dan commit
	tx := db.Begin()

	// jika kita menentukan untuk melakukan transaksi manual, gunakan rollback
	// rollback kita set dengan defer, untuk berjaga-jaga semisal terjadi error pada transaksi,-
	// maka langsung lakukan rollback
	// jika transaksi error maka langsung rollback, namun setelah transaksi selesai dan tidak ada error,-
	// maka tetap di rollback, namun data yang di rollback tidak ada. sehingga berjalan normal seperti biasanya
	defer tx.Rollback()

	// menambahkan data baru 
	err := tx.Create(&User{ID:"13", Password: "rahasia", Name: Name{FirstName: "User 13"}}).Error

	// menambahkan data baru
	err = tx.Create(&User{ID:"14", Password: "rahasia", Name: Name{FirstName: "User 14"}}).Error

	// jika tidak ada error sama sekali pada proses transaksi, maka lakukan commit disini
	if err == nil {
		tx.Commit()
	}
}

// implementasi transaction (manual dan gagal)
func TestManualTransactionFailed(t *testing.T) {
	// membuat transaksi manual baru
	// ketika membuat transaction manual, kita perlu mendefinisikan begin dan commit
	tx := db.Begin()

	// jika kita menentukan untuk melakukan transaksi manual, gunakan rollback
	// rollback kita set dengan defer, untuk berjaga-jaga semisal terjadi error pada transaksi,-
	// maka langsung lakukan rollback
	// jika transaksi error maka langsung rollback, namun setelah transaksi selesai dan tidak ada error,-
	// maka tetap di rollback, namun data yang di rollback tidak ada. sehingga berjalan normal seperti biasanya
	defer tx.Rollback()

	// menambahkan data baru 
	err := tx.Create(&User{ID:"16", Password: "rahasia", Name: Name{FirstName: "User 16"}}).Error

	// menambahkan data baru
	err = tx.Create(&User{ID:"13", Password: "rahasia", Name: Name{FirstName: "User 13"}}).Error
	
	// menambahkan data baru
	err = tx.Create(&User{ID:"17", Password: "rahasia", Name: Name{FirstName: "User 17"}}).Error

	// jika tidak ada error sama sekali pada proses transaksi, maka lakukan commit disini
	if err == nil {
		tx.Commit()
	}
}

// implementasi query (single object) first dan last
func TestQuerySingleObject(t *testing.T) {
	// membuat variabel untuk menampung hasil query select
	user := User{}

	// mengambil data pertama
	err := db.First(&user).Error

	// mengecek dengan assert, pastikan tidak ada error
	assert.Nil(t, err)

	// memastikan data pertama id nya adalah "1"
	assert.Equal(t, "1", user.ID)

	// mengkosongkan variabel user untuk menampung data terakhir
	user = User{}

	// mengambil data paling akhir
	err = db.Last(&user).Error

	// mengecek dengan assert, pastikan tidak ada error
	assert.Nil(t, err)

	// memastikan data terakhir id nya adalah "9"
	// kenapa 9, karena data yang kita miliki varchar paling belakang adalah 9 bukan 16 (kecuali kalau id nya integer)
	assert.Equal(t, "9", user.ID)
}

// implementasi query (single object), inline condition
func TestQuerySingleObjectInlineCondition(t *testing.T) {
	// membuat variabel untuk menyimpan data query
	user := User{}

	// menggunakan Take untuk mengambil data dengan inline condition (where)
	// disarankan pakai take, kalau memang mengambil hanya satu buah data dengan kondisi tertentu (inline condition)
	err := db.Take(&user, "id = ?", 5).Error // mengambil data dengan id karakter 5 paling depan (string)

	// mengecek dengan assert, pastikan tidak ada error
	assert.Nil(t, err)
	assert.Equal(t, "5", user.ID)
	assert.Equal(t, "User 5", user.Name.FirstName)
}

// implementasi query all objects
func TestQueryAllObjects(t *testing.T) {
	// membuat variabel slice untuk menyimpan data hasil query yang datanya nanti lebih dari satu
	var users []User

	// mengambil data lebih dari satu menggunakan method Find
	// dengan menggunakan find, maka akan melakukan select ke semua field dar tabel yang berkaitan (select *)
	err := db.Find(&users, "id in ?", []string{"1", "2", "3", "4", "5"}).Error

	// mengecek dengan assert, pastikan tidak ada error
	assert.Nil(t, err)
	assert.Equal(t, 5, len(users))
}

// implementasi advanced query - query condition
func TestQueryCondition(t *testing.T) {
	// membuat variabel untuk menyimpan data users
	var users []User

	// kondisi pertama (ambil data user dengan first_name nya mengandung kata User)
	// kondisi kedua
	// titik yang digunakan sebagai pemisah where, adalah query 'AND' di mysql
	err := db.Where("first_name like ?", "%User%").Where("password = ?", "rahasia").Find(&users).Error
	
	// mengecek dengan assert, pastikan tidak ada error
	assert.Nil(t, err)
	assert.Equal(t, 16, len(users))
}

// implementasi advanced query - OR operator
func TestOROperator(t *testing.T) {
	// membuat variabel untuk menyimpan data users
	var users []User

	// kondisi pertama (ambil data user dengan first_name nya mengandung kata User)
	// kondisi kedua menggunakan operator OR
	err := db.Where("first_name like ?", "%User%").Or("password = ?", "rahasia").Find(&users).Error
	
	// mengecek dengan assert, pastikan tidak ada error
	assert.Nil(t, err)
	assert.Equal(t, 17, len(users))
}

// implementasi advanced query - NOT operator
func TestNOTOperator(t *testing.T) {
	// membuat variabel untuk menyimpan data users
	var users []User

	// kondisi pertama (ambil data user dengan first_name nya mengandung kata User)
	// kondisi pertama menggunakan operator NOT
	err := db.Not("first_name like ?", "%User%").Where("password = ?", "rahasia").Find(&users).Error
	
	// mengecek dengan assert, pastikan tidak ada error
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
}

// implementasi advanced query - Select Fields
func TestSelectFields(t *testing.T) {
	// membuat variabel untuk menyimpan data users
	var users []User

	// mengambil data dengan menyeleksi kolom dengan select
	err := db.Select("id", "first_name").Find(&users).Error
	
	// mengecek dengan assert, pastikan tidak ada error
	assert.Nil(t, err)

	// mengecek data dengan perulangan
	for _, user := range users {
		// pastikan kolom user id tidak kosong
		assert.NotNil(t, user.ID)
		assert.NotEqual(t, "", user.Name.FirstName)
	}

	// memastikan total data sesuai
	assert.Equal(t, 17, len(users))
}

// implementasi advanced query - struct condition
func TestStructCondition(t *testing.T) {
	// membuat kondisi user menggunakan struct
	userCondition := User{
		Password: "rahasia",
		Name: Name{
			FirstName: "User 5",
			LastName: "", // akan diabaikan karena di anggap default value oleh struct
		},
	}

	// membuat variabel untuk menyimpan users
	var users []User

	// menggunakan user condition sebagai where, dimana-
	// struct field sebagai field dan value sebagai nilai yang dicari
	err := db.Where(userCondition).Find(&users).Error

	// mengecek error dengan assert
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
}

// implementasi advanced query - map condition
func TestMapCondition(t *testing.T) {
	// membuat map baru sebagai condition
	mapCondition := map[string]interface{} {
		"middle_name": "", // akan termasuk ke dalam kondisi pada query nantinya
	}

	// membuat variabel untuk menyimpan users
	var users []User

	// menggunakan map condition sebagai where, dimana-
	// map key sebagai field dan interface sebagai nilai yang dicari
	err := db.Where(mapCondition).Find(&users).Error

	// mengecek error dengan assert
	assert.Nil(t, err)
	assert.Equal(t, 16, len(users))
}

// implementasi advanced query - order, limit dan offset
func TestOrderLimitOffset(t *testing.T) {
	// membuat variabel untuk menyimpan users
	var users []User

	// melakukan query
	// method order bisa melakukan sorting lebih dari satu kolom
	// limit mengambil data 5
	// sedangkan offset menentukan cek data dari data ke berapa. contoh-
	// setelah di urutkan ke asc, akan skip data 1-4 dan dimulai mengambil data dari data urutan ke 5
	err := db.Order("id asc, first_name asc").Limit(5).Offset(5).Find(&users).Error

	// mengecek error dengan assert
	assert.Nil(t, err)
	assert.Equal(t, 5, len(users))
}

// implementasi query non model
// membuat struct user response untuk menyimpan hasil query di sini
type UserResponse struct {
	ID string
	FirstName string
	LastName string
}

func TestQueryNonModel(t *testing.T) {
	// membuat objek untuk menyimpan query ke model users
	var users []UserResponse

	// menggunakan method model untuk menunjukkan model utama yang ingin digunakan sebagai query
	// sedangkan struct UserResponse adalah tempat untuk menyimpan hasil query (dengan select kolom tertentu saja)
	err := db.Model(&User{}).Select("id", "first_name", "last_name").Find(&users).Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err)

	// memastikan jumlah data yang diambil sesuai
	assert.Equal(t, 17, len(users))
}

// implementasi update
func TestUpdate(t *testing.T) {
	// membuat variabel untuk menyimpan data user hasil query
	user := User{}
	err := db.Take(&user, "id = ?", "1").Error

	// memastikan tidak ada error pada query take
	assert.Nil(t, err)

	// mengubah data yang sudah diambil
	user.Name.FirstName = "Muhammad"
	user.Name.MiddleName = "Ilham"
	user.Name.LastName = "Nurizky"
	user.Password = "rahasia123"

	// menyimpan hasil perubahan data dengan method save
	err = db.Save(user).Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err)
}

// implementasi update lebih dari satu kolom (tidak mencakup semua kolom yang di update)
func TestUpdateSelectedColumns(t *testing.T) {
	// # Cara 1 - Updates pendekatan map
	// melakukan query data lebih dari satu kolom menggunakan updates (pendekatan map)
	err := db.Model(&User{}).Where("id = ?", "2").Updates(map[string]interface{}{
		// lakukan update ke kolom yang ingin di update, bisa lebih dari satu
		"first_name": "Taufik",
		"middle_name": "",
		"last_name": "Hidayat",
	}).Error

	// memastikan tidak ada error pada querys
	assert.Nil(t, err)

	// melakukan query data hanya satu kolom saja menggunakan update
	// method update memiliki parameter key, dan value
	err = db.Model(&User{}).Where("id = ?", "2").Update("password", "rahasia456").Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err)

	// # Cara 2 -  Updates pendekatan Struct
	// jika struct nya sudah sesuai (User), maka boleh langsung panggil method Updates() nya saja-
	// tanpa perlu mendefinisikan Model() User
	err = db.Where("id = ?", "3").Updates(User{
		// lakukan update ke kolom yang ingin di update, bisa lebih dari satu
		Name: Name{
			FirstName: "Dimas",
			// MiddleName: "", // kalau tidak ditambahkan akan memberikan nilai sama (tidak berpengaruh)
			LastName: "Prayoga",
		},
		Password: "rahasia789",
	}).Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err)
}

// implementasi auto_increment
func TestAutoIncrement(t *testing.T) {
	// melakukan perulangan untuk menambahkan data lebih dari satu-
	// dimana hanya mengisikan beberapa kolom saja, untuk menguji apakah kolom id (primary key) auto increment
	for i := 0; i < 10; i++ {
		// membuat data untuk user log
		userLog := UserLog{
			UserId: "1",
			Action: "Test Action",
		}

		// melakukan insert data ke database
		err := db.Create(&userLog).Error

		// memastikan tidak ada error pada query
		assert.Nil(t, err)

		// memastikan auto increment, dengan mengecek id dari setiap data yang berhasil di buat bukan 0
		assert.NotEqual(t, 0, userLog.ID)
		fmt.Println(userLog.ID) // menampilkan id yang berhasil dibuat ke database
	}
}

// implementasi  - auto increment
func TestSaveOrUpdate(t *testing.T) {
	// membuat data struct user log untuk ditambahkan dan di ubah ke database
	userLog := UserLog{
		UserId: "1",
		Action: "Action Test",
	}

	// melakukan query dengan save
	// ketika data stuct tidak mendefinisikan ID (primary key), maka akan dianggap operasi dibawah ini adalah create / insert
	err := db.Save(&userLog).Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err)

	// ketika berhasil dibuat, berarti data userLog sudah mempunyai ID
	// sehingga ketika diquery kembali dengan Save()
	userLog.UserId = "2" // maka akan menjadi update ketika di eksekusi
	err = db.Save(&userLog).Error
	
	// memastikan tidak ada error pada query
	assert.Nil(t, err)

	// dengan demikian method Save(), bisa melakukan upsert (update/insert)
	// hal ini cocok untuk data yang primary key nya auto increment
}

// implementasi upsert - non auto increment
func TestSaveOrUpdateNonAutoIncrement(t *testing.T) {
	// membuat data struct user log untuk ditambahkan dan di ubah ke database
	user := User{
		ID: "99",
		Name: Name{
			FirstName: "User 99",
		},
	}

	// melakukan save untuk pertama kali
	// maka method save tidak langsung melakukan create, karena ID ada di data struct diatas,-
	// maka method save akan melakukan update terlebih dahulu, jika datanya ada di database,-
	// maka akan di update. Jika tidak ada data/baris yang ditemukan = 0, maka lakukan create/insert.-
	// atau jika di update tidak ada data yang berubah, maka lakukan inser
	err := db.Save(&user).Error // create
	
	// memastikan tidak ada error pada query
	assert.Nil(t, err)

	// kemudian kita coba ubah salah satu data struct diatas setelah ditambahkan,
	// maka dakan melakukan update, karena data user sudah ada isinya/di create sebelumnya
	user.Name.FirstName = "User 99 Updated" 
	err = db.Save(&user).Error // update

	// memastikan tidak ada error pada query
	assert.Nil(t, err)
}

// implementasi upsert - Conflict (data duplikat)
func TestConflict(t *testing.T) {
	// membuat data struct user log untuk ditambahkan dan di ubah ke database
	user := User{
		ID: "88",
		Name: Name{
			FirstName: "User 88",
		},
	}

	// melakukan create dengan menambahkan clauses
	// jika pada saat insert duplikat (data yang primary key dari struct sudah ada),-
	// maka akan melakukan update secara otomatis
	err := db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&user).Error 

	// memastikan tidak ada error pada query
	assert.Nil(t, err)
}

// implementasi delete
func TestDelete(t *testing.T) {
	// # Cara ke 1
	// mengambil data user terlebih dahulu
	var user User
	err := db.Take(&user, "id = ?", "88").Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err)

	// setelah di dapatkan, kemudian lakukan delete data
	err = db.Delete(&user).Error
	
	// memastikan tidak ada error pada query
	assert.Nil(t, err)

	// Cara ke 2
	// langsung delete data tanpa diambil terlebih dahulu
	err = db.Delete(&User{}, "id = ?", "99").Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err)

	// Cara ke 3 
	// menghapus data menggunakan method where, dan method delete dengan memberikan entity nya
	err = db.Where("id = ?", "9").Delete(&User{}).Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err)
}

// implementasi soft delete
func TestSoftDelete(t *testing.T) {
	// membuat data struct todo
	todo := Todo {
		UserId: "1",
		Title: "Todo 1",
		Description: "Description 1",
	}

	// melakukan create / insert data todo ke database
	err := db.Create(&todo).Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 

	// setelah di tambahkan, coba kita hapus data todo yang baru saja ditambahkan
	err = db.Delete(&todo).Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
	assert.NotNil(t, todo.DeletedAt) // memastikan bahwa data kolom deleted_at telah di isi (soft delete) 

	// ketika selesai dihapus sebenanrnya data nya masih ada di database, sehingga proses hapus sebelumnya-
	// adalah proses mengisikan kolom deleted_at dengan waktu hapus. sebagai penanda bahwa data itu telah dihapus (soft delete)
	// sehingga ketika data sudah di dihapus (dalam soft delete), maka ketika dilakukan query, data tersebut tidak ditemukan-
	// karena GORM sudah mengerti bahwa data itu sudah termasuk ke dalam soft delete

	// mengambil data todo yang sudah di delete sebelumnya
	var todos []Todo
	err = db.Find(&todos).Error // mengambil semua data todos
	
	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
	assert.Equal(t, 0, len(todos))
}

// implementasi soft delete - unscoped
func TestUncscoped(t *testing.T) {
	// membuat variabel todo untuk menyimpan hasil query
	var todo Todo

	// mengambil data todo yang sudah termasuk ke dalam soft delete
	// err := db.First(&todo, "id = ?", 1).Error

	// memastikan tidak ada error pada query
	// assert.Nil(t, err) // akan error 

	// menampilkan data todo ke output
	// fmt.Println(todo) // akan error, karena data dianggap sudah di delete (soft delete)

	// namun jika kita ingin mengambil data yang sudah di delete dengan metode soft delete
	// kita bisa gunakan method Unscoped()

	// mengambil data todo yang sudah termasuk ke dalam soft delete dengan method Unscoped
	err := db.Unscoped().First(&todo, "id = ?", 2).Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
	fmt.Println(todo) // sukses

	// jadi meski kita hapus ulang dengan db.Delete(), maka tidak menghapus data, namun update kolom deleted_at
	// jika kita ingin benar benar menghapus data dari database bisa menggunakan unscoped lagi
	// err = db.Delete(&todo).Error // akan update deleted_at, bukan menghapus data
	err = db.Unscoped().Delete(&todo).Error // sukses, delete data secara permanen

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
	fmt.Println(todo) // sukses
}

// implementasi lock
func TestLock(t *testing.T) {
	// penggunaan locking cocok dilakukan pada transaction
	err := db.Transaction(func(tx *gorm.DB) error {
		// membuat variabel untuk menyimpan data user
		var user User

		// menggunakan transaction untuk lock
		// strength update hanya untuk data update, jika ingin mengambil data gunakan 'SHARE
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Take(&user, "id = ?", "1").Error

		// mengecek error
		if err != nil {
			return err
		}

		// update data user
		user.Name.FirstName = "Dimas"
		user.Name.MiddleName = ""
		user.Name.LastName = "Prayoga"

		// menyimpan perubahan
		err = tx.Save(&user).Error

		// return err (kalau error nil maupun ada, maka langsung keluar)
		return err
	})

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
}

// implementasi one to one (has one)
func TestCreateWallet(t *testing.T) {
	// menyiapkan data wallet yang akan di insert ke database
	wallet := Wallet{
		ID: "1",
		UserId: "1",
		Balance: 1000000,
	}

	// menambahkan data ke database
	err := db.Create(&wallet).Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
}

// implementasi one to one - preload
func TestRetrieveRelation(t *testing.T) {
	// menyimpan data user
	var user User

	// mengambil data user dan relasi nya ke model wallet
	// "Wallet", dalam preload adalah field relation yang di dapat dari model user
	err := db.Model(&User{}).Preload("Wallet").Take(&user, "id = ?", "1").Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
	assert.Equal(t, "1", user.ID) 
	assert.Equal(t, "1", user.Wallet.ID) 
}

// implementasi one to one (has one) - join
func TestRetrieveRelationJoin(t *testing.T) {
	// menyimpan data user
	var user User

	// mengambil data user dan relasi nya ke model wallet
	// menggunakan join, karena lebih cepat untuk data table yang relasinya has 
	// users.id adalah merujuk data utama yang ingin diambil, karena menggunakan join,-
	// pasti melakukan seleksi kedua tabel, sehingga diberikan nama tabelnya
	err := db.Model(&User{}).Joins("Wallet").Take(&user, "users.id = ?", "1").Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
	assert.Equal(t, "1", user.ID) 
	assert.Equal(t, "1", user.Wallet.ID) 
}

// implementasi upsert relation
func TestAutoCreateUpdate(t *testing.T) {
	// menyiapkan data user yang ingin ditambahkan
	user := User {
		ID: "20",
		Password: "rahasia",
		Name: Name{
			FirstName: "User 20",
		},

		// dan juga menambahkan sekalian untuk data wallet,-
		// nah jika data wallet belum ada, berhubung ini relasi. 
		// maka GORM akan otomatis menambahkan data wallet, jika data wallet belum ada
		Wallet: Wallet{
			ID: "20",
			UserId: "20",
			Balance: 1000000,
		},
	}

	// melakukan insert data user sekaligus wallet (jika belum ada di database)
	// jikalau data sudah ada di database, maka otomatis akan melakukan update berdasarkan ID
	err := db.Create(&user).Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
}

func TestSkipAutoCreateUpdate(t *testing.T) {
	// menyiapkan data user yang ingin ditambahkan
	user := User {
		ID: "21",
		Password: "rahasia",
		Name: Name{
			FirstName: "User 21",
		},

		// tetap melampirkan data wallet
		Wallet: Wallet{
			ID: "21",
			UserId: "21",
			Balance: 1000000,
		},
	}

	// melakukan insert data user tanpa meng-create atau update data wallet
	// dengan enggunakan db.Omit(clauses.Associations)
	// sehingga akan insert data user saja, data wallet akan di skip
	err := db.Omit(clause.Associations).Create(&user).Error 

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
}

// implementasi one to many
func TestUserAndAddresses(t *testing.T) {
	// menyiapkan data user baru
	user := User {
		ID: "50",
		Password: "rahasia",
		Name: Name{
			FirstName: "User 50",
		},

		// sekaligus menambahkan data wallet
		// untuk relasi one to one, jika ingin melakukan select bisa menggunakan JOIN
		Wallet: Wallet{
			ID: "50",
			UserId: "50",
			Balance: 500000,
		}, 

		// relasi one to many juga memiliki upsert data relation
		// sehingga ketika data tidak ada pada saat menambahkan user baru, dan data tersebut-
		// di jabarkan dalam bentuk struct, maka data akan secara otomatis di tambahkan ke database oleh GORM
		// sedangkan untuk one to many, jika ingin menggunakan select gunakan PRELOAD
		Addresses: []Address{
			{
				// karena id data address bersifat auto increment, maka upsert data relation akan-
				// otomatis menambahkan data yang sesuai
				UserId: "50",
				Address: "Indonesia",
			},
			{
				UserId: "50",
				Address: "Banyuwangi",
			},
		},
	}

	// menambahkan data user dan data foreign key nya (wallet dan address)
	// karena bersifat upsert data relation
	err := db.Create(&user).Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
}

func TestPreloadJoinOneToMany(t *testing.T) {
	// menyimpan hasil select data dengan preload dan join
	var users []User

	// melakukan query untuk select data user
	// preload digunakan untuk select data untuk relasi one to many
	// sedangkan join digunakan untuk select data untuk relasi one to one
	// mengambil data lebih dari satu
	err := db.Model(&User{}).Preload("Addresses").Joins("Wallet").Find(&users).Error 
	
	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
}

func TestTakePreloadJoinOneToMany(t *testing.T) {
	// menyimpan hasil select data dengan preload dan join
	var user User

	// melakukan query untuk select data user
	// preload digunakan untuk select data untuk relasi one to many
	// sedangkan join digunakan untuk select data untuk relasi one to one
	// hanya mengambil satu data
	err := db.Model(&User{}).Preload("Addresses").Joins("Wallet").Take(&user, "users.id = ?", "50").Error 
	
	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
}

// implementasi belongs to
func TestBelongsToAddress(t *testing.T) {
	// jika ingin mengambil data Address dengan include data user, maka bisa menggunakan preload dan join-
	// karena sifatnya belongsto (User memiliki banyak address)
	// namun jika kita ingin mengambil data user dengan include data Address, maka harus menggunakan Preload-
	// karena 1 user kemungkinan bisa memiliki lebih dari 1 address, sehingga tidak bisa menggunakan Join

	// 1. menggunakan preload
	fmt.Println("Preload")
	
	// menyipakan data address untuk select dengan preload
	var addresses []Address
	
	// melakukan query dengan preload dan include data user
	err := db.Model(&Address{}).Preload("User").Find(&addresses).Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 

	// 2. menggunakan join
	fmt.Println("Join")

	// menyipakan data address untuk select dengan preload
	var address Address
	
	// melakukan query dengan preload dan include data user
	err = db.Model(&Address{}).Joins("User").Find(&address).Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
}

func TestBelongsToWallet(t *testing.T) {
	// jika ingin mengambil data Wallet dengan include data user, maka bisa menggunakan preload dan join-
	// karena sifatnya belongsto (User memiliki satu Wallet)
	// namun jika kita ingin mengambil data user dengan include data Wallet, maka harus menggunakan Preload-
	// karena 1 user kemungkinan bisa memiliki lebih dari 1 Wallet, sehingga tidak bisa menggunakan Join

	// 1. menggunakan preload
	fmt.Println("Preload")
	
	// menyipakan data wallet untuk select dengan preload
	var wallets []Wallet
	
	// melakukan query dengan preload dan include data user
	err := db.Model(&Wallet{}).Preload("User").Find(&wallets).Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 

	// 2. menggunakan join
	fmt.Println("Join")

	// menyipakan data wallet untuk select dengan preload
	var wallet Wallet
	
	// melakukan query dengan preload dan include data user
	err = db.Model(&Wallet{}).Joins("User").Find(&wallet).Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
}

// implementasi many to many
func TestCreateManyToMany(t *testing.T) {
	// menyiapkan data product, untuk melakukan simulasi pengujian
	product := Product{
		ID: "P001",
		Name: "Contoh Product",
		Price: 200000,
	}

	// menambahkan data product terlebih dahulu di database
	err := db.Create(&product).Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 

	// untuk melakukan create, update dan delete di tabel penghubung (konsep many to many)
	// maka bisa langsung gunakan method Table(), karena tabel penhubung tidak menerapkan struct
	err = db.Table("user_like_product").Create(map[string]interface{} {
		"user_id": "1",
		"product_id": product.ID,
	}).Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
	
	err = db.Table("user_like_product").Create(map[string]interface{} {
		"user_id": "2",
		"product_id": product.ID,
	}).Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
}

func TestPreloadManyToManyProduct(t *testing.T) {
	// menyimpan data product
	var product Product

	// melakukan query dengan preload, untuk mengambil siapa saja (user) yang menyukai sebuah product
	// preload liked by users, adalah field relasi pada tabel Product
	err := db.Preload("LikedByUsers").Take(&product, "id = ?", "P001").Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
	assert.Equal(t, 2, len(product.LikedByUsers)) // karena yang like product tadi ada 2 user pada pengujian sebelumnya
}

func TestPreloadManyToManyUser(t *testing.T) {
	// menyimpan data user
	var user User

	// melakukan query dengan preload, untuk mengambil product apa saja yang disukai oleh user
	// preload like products, adalah field relasi pada tabel User
	err := db.Preload("LikeProducts").Take(&user, "id = ?", "1").Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
	assert.Equal(t, 1, len(user.LikeProducts))
}

// implementasi association mode
func TestAssociation(t *testing.T) {
	// menyiapkan data product
	var product Product

	// mengambil data produk berdasarkan id tertentu
	err := db.Take(&product, "id = ?", "P001").Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 

	// menyiapkan data users
	var users []User

	// mengambil data user yang menyukai product ini, berdasarkan nilai tertentu pada suatu kolom
	// dengan mengunakan method Association(), kita bisa melakukan query untuk filtering data yang berelasi (user)
	// juga dengan kita menambahkan relasi pada tabel Product ke method Association(), maka sebenarnya-
	// kita melakukan query di table user saat itu juga
	err = db.Model(&product).Where("users.first_name LIKE ?", "%User%").Association("LikedByUsers").Find(&users)

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
	assert.Equal(t, 1, len(users))
}

func TestAssociationAppend(t *testing.T) {
	// menyiapkan data user
	var user User

	// mengambil sebuah data user
	err := db.Take(&user, "id = ?", "3").Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 

	// menyiapkan data product
	var product Product

	// mengambil sebuah data product
	err = db.Take(&product, "id = ?", "P001").Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 

	// method append digunakan untuk menambahkan data sepertihalnya pada baris 1036 (menambahkan data ke database)
	// bedanya, kalau yang sebelumnya, kita menyebutkan table. sedangkan ini tidak menggunakan tabel
	// cukup dengan mengambil model dasar nya yaitu Product (karena ingin menambahkan data berdasarkan relasi LikedByUsers)-
	// dan method Append adalah user, maka artinya user akan like product
	err = db.Model(&product).Association("LikedByUsers").Append(&user)

	// jika menggunakan method Append, meskipun data sudah ada, dan ketika mau di insert dengan data yang sama,-
	// maka akan melakukan update

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
}

func TestAssociationReplace(t *testing.T) {
	// method Replace(), cocok untuk relasi one to one atau belongs to
	// menggunakan db transaction, karena terdapat beberapa operasi berulang (menghapus relasi lama, dan menginsert relasi yang baru)
	// sehingga disarankan untuk replace menggunakan transaction
	err := db.Transaction(func(tx *gorm.DB) error {
		// menyiapkan data user
		var user User

		// mengambil sebuah data user
		err := tx.Take(&user, "id = ?", "1").Error
		
		// memastikan tidak ada error pada query
		assert.Nil(t, err) 

		// menyiapkan data Wallet baru
		wallet := Wallet{
			ID: "01",
			UserId: user.ID,
			Balance: 800000,
		}

		// melakukan replace (pergantian data yang sudah ada), dengan method Replace
		// jadi pada kode di bawah ini 
		err = tx.Model(&user).Association("Wallet").Replace(&wallet)

		return err
	})

	// ! pengujian ini ketika dijalankan akan menyebabkan error, penjelasan di note (baris 333)

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
}

func TestAssociationDelete(t *testing.T) {
	// menyiapkan data user
	var user User

	// mengambil sebuah data user
	err := db.Take(&user, "id = ?", "3").Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 

	// menyiapkan data product
	var product Product

	// mengambil sebuah data product
	err = db.Take(&product, "id = ?", "P001").Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 

	// user 3 sudah tidak liked ke product P001, kita bisa hapus relasi menggunakan Delete()
	err = db.Model(&product).Association("LikedByUsers").Delete(&user)

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
}

func TestAssociationClear(t *testing.T) {
	// menyiapkan data product
	var product Product

	// mengambil sebuah data product
	err := db.Take(&product, "id = ?", "P001").Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 

	// menghapus semua data relasi dari tabel Product ke tabel User, melalui tabel user like product-
	// dengan menggunakan method Clear()
	// menjadikan product P001, tidak di sukai oleh user manapun lagi
	err = db.Model(&product).Association("LikedByUsers").Clear()

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
}

// implementasi preloading
func TestPreloading(t *testing.T) {
	// menyiapkan data user
	var user User

	// melakukan query untuk mengambil data user (id = "1") yang memiliki nilai wallet sekian
	// untuk melakukan query ke tabel Wallet, kita bisa menggunakan relation pada tabel User yaitu "Wallet"-
	// yang mana akan kita masukkan ke dalam method Preload
	// dan menambahkan inline condition, untuk kondisi yang kita butuhkan
	err := db.Preload("Wallet", "balance >= ?", 1000000).Take(&user, "id = ?", "1").Error

	// jika user dengan id 1 balance nya tidak memenuhi kondisi, maka data wallet tidak akan diambil
	fmt.Println(user)

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
}

// implementasi nested preloading
func TestPreloadNested(t *testing.T) {
	// menyiapkan data wallet
	var wallet Wallet

	// melakukan query ke data wallet, dan mengambil data relasinya yaitu User, dan Addresses milik User
	err := db.Preload("User.Addresses").Take(&wallet, "id = ?", "50").Error
	
	// memastikan tidak ada error pada query
	assert.Nil(t, err) 

	// menampilkan data wallet, user dan addresses milik user
	fmt.Println(wallet)
	fmt.Println(wallet.User)
	fmt.Println(wallet.User.Addresses)
}

// implementasi preload all
func TestPreloadAll(t *testing.T) {
	// menyiapkan data user
	var user User

	// melakukan preload ke seluruh relasi di tabel User (mencakup Wallet, Addresses, dan LikeProducts)
	// namun perlu diperhatikan untuk preload all di sini, karena banyak query yang dilakukan
	err := db.Preload(clause.Associations).Take(&user, "id = ?", "1").Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 

	// menampilkan data user
	fmt.Println(user)
}

// implementasi joins
func TestJoinQuery(t *testing.T) {
	// menyiapkan data users
	var users []User

	// melakukan joins secara manual
	// mirip seperti inner join, wajib ada untuk kedua data user dan wallet nya
	err := db.Joins("join wallets on wallets.user_id = users.id").Find(&users).Error 

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
	assert.Equal(t, 3, len(users))

	// mengosongkan data users
	users = []User{}

	// melakukan joins dengan field relation pada model User
	err = db.Joins("Wallet").Find(&users).Error // defaultnya menggunakan query left join

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
	assert.Equal(t, 19, len(users))
}

// implementasi join dengan pengkondisian
func TestJoinWithCondition(t *testing.T) {
	// menyiapkan data users
	var users []User

	// menggunakan join dengan kondisi (manual)
	// jika join secara manual, wajib menyebutkan nama tabel untuk pengkondisiannya
	// akan menghasilkan data yang tanpa melampirkan wallet
	err := db.Joins("join wallets on wallets.user_id = users.id AND wallets.balance > ?", 500000).Find(&users).Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
	assert.Equal(t, 2, len(users))

	// mengkosongkan data users
	users = []User{}

	// menggunakan join dengan kondisi (menggunakan field / alias)
	// jika kita menambahkan kondisi pada join, dengan cara ini cukup memanggil alias (field relationnya)
	// menghasilkan data dengan melampirkan datanya
	err = db.Joins("Wallet").Where("Wallet.balance > ?", 500000).Find(&users).Error // menggunakan field relation (Wallet)

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
	assert.Equal(t, 2, len(users))
}

// implementasi query aggregation
func TestCount(t *testing.T) {
	// menyiapkan data count
	var count int64

	// menghitung data User yang memiliki balance diatas sekian (dengan kondisi tertentu)
	err := db.Model(&User{}).Joins("Wallet").Where("Wallet.balance > ?", 500000).Count(&count).Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
	assert.Equal(t, int64(2), count)
}

// implementasi query aggregation yang lain (manual)
type AggregationResult struct {
	TotalBalance int64
	MinBalance int64
	MaxBalance int64
	AvgBalance float64
}

func TestAggregation(t *testing.T) {
	// menyiapkan hasil aggregation
	var result AggregationResult

	// melakukan aggregation secara manual dengan select (untuk selain count)
	err := db.Model(&Wallet{}).Select("sum(balance) as total_balance", "min(balance) as min_balance",
	"max(balance) as max_balance", "avg(balance) as avg_balance").Take(&result).Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
	assert.Equal(t, int64(3000000), result.TotalBalance)
	assert.Equal(t, int64(1000000), result.MinBalance)
	assert.Equal(t, int64(1000000), result.MaxBalance)
	assert.Equal(t, float64(1000000), result.AvgBalance)
}


func TestAggregationGroupByHaving(t *testing.T) {
	// menyiapkan hasil aggregation
	var results []AggregationResult

	// melakukan aggregation secara manual dengan select (untuk selain count)
	// menambahkan joins untuk group by dan having - untuk mengelompokkkan user dengan balance diatas sekian
	err := db.Model(&Wallet{}).Select("sum(balance) as total_balance", "min(balance) as min_balance",
	"max(balance) as max_balance", "avg(balance) as avg_balance").
	Joins("User").Group("User.id").Having("sum(balance) > ?", 1000000).Find(&results).Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
	assert.Equal(t, 0, len(results)) // karena tidak ada balanace yang diatas 1 juta
}

// implementasi context
func TestContext(t *testing.T) {
	// membuat context baru
	ctx := context.Background()

	// menyiapkan data users
	var users []User

	// melakukan query dengan context
	err := db.WithContext(ctx).Find(&users).Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
	assert.Equal(t, 19, len(users)) 
}

// implementasi scopes
// dengan scopes memungkinkan kita untuk melakukan kustomisasi logic pada query database
func BrokeWalletBalance(db *gorm.DB) *gorm.DB {
	// di dalam function scopes, kita bisa tambahkan kustomisasi query yang kita inginkan
	return db.Where("balance = ?", 0)
}

func SultanWalletBalance(db *gorm.DB) *gorm.DB {
	// di dalam function scopes, kita bisa tambahkan kustomisasi query yang kita inginkan
	return db.Where("balance >= ?", 1000000)
}

func TestScopes(t *testing.T) {
	// menyiapkan data wallets
	var wallets []Wallet

	// memanggil function scopes yang sudah dibuat sebelumnya,-
	// dengan menggunakan scopes query yang kita berikan lebih ringkas dan bisa digunakan kembali (reusable)
	err := db.Scopes(BrokeWalletBalance).Find(&wallets).Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
	assert.Equal(t, 0, len(wallets)) 

	// mengosongkan data wallets
	wallets = []Wallet{}

	// memanggil function scopes kedua
	err = db.Scopes(SultanWalletBalance).Find(&wallets).Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
	assert.Equal(t, 3, len(wallets)) 
}

// implementasi Migrator
func TestMigrator(t *testing.T) {
	// melakukan migrasi dari struct ke table database secara otomatis dengan method AutoMigrate
	err := db.Migrator().AutoMigrate(&GuestBook{})

	// disaranakan menggunakan library migration golang, jangan gunakan migrator
	// migrator hanya untuk pengetesan di komputer lokal saja

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
}

// implementasi hook
func TestHook(t *testing.T) {
	// menyiapkan data user
	user := User{
		ID: "", // kan mentrigger method BeforeCreate() pada model User
		Password: "rahasia",
		Name: Name{
			FirstName: "User Saya",
		},
	}

	// melakukan create
	err := db.Create(&user).Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
	assert.NotEqual(t, "", user.ID)

	// menampilkan id user
	fmt.Println(user.ID)
}