package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // Driver MySQL
	"github.com/joho/godotenv"        // Pustaka pemuat .env
)

// InitDB bertugas membuka koneksi ke MySQL
func InitDB() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Println("Info: File .env tidak ditemukan, menggunakan environment sistem.")
	}

	// 2. Mengambil data dari environment
	host := os.Getenv("DB_MYSQL_HOST")
	port := os.Getenv("DB_MYSQL_PORT")
	user := os.Getenv("DB_MYSQL_USER")
	password := os.Getenv("DB_MYSQL_PASSWORD")
	dbname := os.Getenv("DB_MYSQL_NAME")

	// 3. Merakit DSN (Data Source Name)
	// Menambahkan charset utf8mb4 agar support emoji/karakter khusus
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local&charset=utf8mb4", 
		user, password, host, port, dbname)

	// 4. Membuka koneksi
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Gagal membuka koneksi database: %v", err)
	}

	// 5. Verifikasi koneksi
	err = db.Ping()
	if err != nil {
		log.Fatalf("Database MySQL tidak merespon: %v", err)
	}

	fmt.Println("🚀 Koneksi MySQL Berhasil!")
	return db
}