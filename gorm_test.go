package belajar_go_lang_gorm

import (
	"fmt"
	"strconv"
	"testing"

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
	})

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
	err := db.Unscoped().First(&todo, "id = ?", 1).Error

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
	fmt.Println(todo) // sukses

	// jadi meski kita hapus ulang dengan db.Delete(), maka tidak menghapus data, namun update kolom deleted_at
	// jika kita ingin benar benar menghapus data dari database bisa menggunakan unscoped lagi
	// err = db.Delete(&todo).Error // akan update deleted_at, bukan menghapus data
	err = db.Unscoped().Delete(&todo).Error // sukses, delete data secara permanen

	// memastikan tidak ada error pada query
	assert.Nil(t, err) 
	fmt.Println(todo) // s
}