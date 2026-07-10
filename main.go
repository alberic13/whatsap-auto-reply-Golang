package main

import (
	"log"
	"os"
	"time"

	jwtV3 "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"

	"github.com/joho/godotenv"
	"sopingi.com/fikom/controllers"
	"sopingi.com/fikom/fungsi"
	"sopingi.com/fikom/models"
	"sopingi.com/fikom/wa"
	"sopingi.com/fikom/ai"
)

func main() {

	//membaca file .env
	godotenv.Load()

	//panggil koneksi
	db := koneksi()

	//auto migrasi model ke db
	db.AutoMigrate(&models.Suhu{})
	db.AutoMigrate(&models.Informasi{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Pesanan{})
	db.AutoMigrate(&models.Dokumen{})
	db.AutoMigrate(&models.Pesan{})

	// Seed data awal jika tabel pesans masih kosong
	var count int64
	db.Model(&models.Pesan{}).Count(&count)
	if count == 0 {
		db.Create(&models.Pesan{
			Kode:      "info",
			Balasan:   "[PESAN OTOMATIS] ini konten informasi",
			CreatedAt: time.Now(),
		})
		db.Create(&models.Pesan{
			Kode:      "prodi",
			Balasan:   "[PESAN OTOMATIS] ini konten program studi",
			CreatedAt: time.Now(),
		})
		db.Create(&models.Pesan{
			Kode:      "sosmed",
			Balasan:   "[PESAN OTOMATIS] ini konten sosial media",
			CreatedAt: time.Now(),
		})
		log.Println("Tabel pesans berhasil di-seed dengan data awal (info, prodi, sosmed).")
	}

	// Jalankan WhatsApp Bot Auto-Reply secara asynchronous
	// go wa.KonekWa(db)

	r := gin.Default()

	//menambahkan middleware JWT
	key_jwt := os.Getenv("KEY_JWT")
	authMiddleware, err := jwtV3.New(&jwtV3.GinJWTMiddleware{
		Realm:       "fikom UDB",
		Key:         []byte(key_jwt),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour * 24,
		IdentityKey: "id",

		PayloadFunc: func(data any) jwt.MapClaims {
			value, ok := data.(models.User)
			if ok {
				return jwt.MapClaims{
					"id":   value.ID,
					"nama": value.Nama,
				}
			}
			return jwt.MapClaims{}
		},

		Authenticator: controllers.UserLogin,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	errInit := authMiddleware.MiddlewareInit()
	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	//membuat variabel db untuk membawa koneksi
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	//route tanpa middleware
	r.POST("/login", authMiddleware.LoginHandler)

	//route group dengan middleware jwt
	auth := r.Group("/backend", authMiddleware.MiddlewareFunc())

	auth.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": true,
			"pesan":  "Berhasil tampil",
		})
	})

	auth.POST("/programstudi", fungsi.BacaDataProdi)

	//route suhu
	auth.GET("/suhu", controllers.Tampil)
	auth.POST("/suhu", controllers.Tambah)
	auth.PUT("/suhu", controllers.Ubah)
	auth.DELETE("/suhu", controllers.Hapus)

	//route informasi
	auth.GET("/informasi", controllers.InformasiTampil)
	auth.POST("/informasi", controllers.InformasiTambah)
	auth.PUT("/informasi", controllers.InformasiUbah)
	auth.DELETE("/informasi", controllers.InformasiHapus)

	//route user
	auth.GET("/user", controllers.UserTampil)
	auth.POST("/user", controllers.UserTambah)
	auth.PUT("/user", controllers.UserUbah)
	auth.DELETE("/user", controllers.UserHapus)

	//route pesanan
	auth.GET("/pesanan", controllers.PesananTampil)
	auth.POST("/pesanan", controllers.PesananTambah)
	auth.PUT("/pesanan", controllers.PesananUbah)
	auth.DELETE("/pesanan", controllers.PesananHapus)
	auth.POST("/drive", controllers.DriveUpload)
	auth.GET("/drive", controllers.DriveTampil)
	auth.GET("/drive/:id", controllers.DriveUnduh)

	//route pesan
	auth.GET("/pesan", controllers.PesanTampil)
	auth.POST("/pesan", controllers.PesanTambah)
	auth.PUT("/pesan", controllers.PesanUbah)
	auth.DELETE("/pesan", controllers.PesanHapus)

	//jika route tidak ada
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"status": false,
			"pesan":  "Route Tidak Ditemukan",
		})
	})


	//membaca port dari file .env
	port := os.Getenv("PORT")
	go r.Run(":" + port)
	//ai.MulaiChatAi()
	ai.InitAi()
	wa.KonekWa(db)

}
