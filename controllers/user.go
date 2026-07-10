package controllers

import (
	"crypto/sha1"
	"fmt"
	"net/http"

	jwtV3 "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sopingi.com/fikom/models"
)

// Binding dari POST JSON
type StrukturUserTambah struct {
	Nama     string `binding:"required"`
	Username string `binding:"required"`
	Password string `binding:"required"`
}

// Binding dari PUT JSON
type StrukturUserUbah struct {
	Id       uint   `binding:"required"`
	Nama     string `binding:"required"`
	Username string `binding:"required"`
	Password string `binding:"required"`
}

// Binding dari DELETE JSON
type StrukturUserHapus struct {
	Id uint `binding:"required"`
}

// Binding dari POST JSON
type StrukturLogin struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
}

func hashPassword(password string) string {
	checksum := sha1.Sum([]byte(password))
	return fmt.Sprintf("%x", checksum)
}

func setupJWT() *jwtV3.GinJWTMiddleware {
	return &jwtV3.GinJWTMiddleware{}
}

func UserLogin(c *gin.Context) (any, error) {
	//ambil koneksi variabel db dari main
	db := c.MustGet("db").(*gorm.DB)
	// membuat variabel data User dengan struktur user dan menangkap data dari request
	var dataUser StrukturLogin
	if err := c.ShouldBindJSON(&dataUser); err != nil {
		//kembalikan data kosong dan eror input login
		return nil, jwtV3.ErrMissingLoginValues
	}

	//enkripsi password dengan sha1
	var sha = sha1.New()
	sha.Write([]byte(dataUser.Password))
	var encrypted = sha.Sum(nil)
	var encryptedString = fmt.Sprintf("%x", encrypted)

	//membuat variabel model user
	var modelUser models.User
	//mencari data user berdasarkan username dan password
	cekUser := db.Where("username = ?", dataUser.Username).Where("password = ?", encryptedString).First(&modelUser)
	if cekUser.Error == nil {
		//kembalikan data user dan eror=nil
		return modelUser, nil
	} else {
		//kembalikan data kosong dan eror gagal login
		return nil, jwtV3.ErrFailedAuthentication
	}
}

func UserTampil(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var modelUser []models.User
	hasil := db.Find(&modelUser)
	kesalahan := hasil.Error

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil Tampil data user",
			"kesalahan": nil,
			"data":      modelUser,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal Tampil data user",
			"kesalahan": kesalahan.Error(),
			"data":      nil,
		})
	}
}

func UserTambah(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dataUser StrukturUserTambah

	if err := c.ShouldBindJSON(&dataUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca Data",
			"kesalahan": err.Error(),
		})
		return
	}

	//enkripsi password dengan sha1
	var sha = sha1.New()
	sha.Write([]byte(dataUser.Password))
	var encrypted = sha.Sum(nil)
	var encryptedString = fmt.Sprintf("%x", encrypted)

	//membuat data baru dengan model user
	modelUser := models.User{
		Nama:     dataUser.Nama,
		Username: dataUser.Username,
		Password: encryptedString,
	}

	hasil := db.Create(&modelUser)
	kesalahan := hasil.Error
	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil tambah data user",
			"kesalahan": nil,
			"data":      modelUser,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal tambah data user",
			"kesalahan": kesalahan.Error(),
			"data":      modelUser,
		})
	}
}

func UserUbah(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dataUser StrukturUserUbah

	if err := c.ShouldBindJSON(&dataUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca Data",
			"kesalahan": err.Error(),
		})
		return
	}

	//membuat variabel model user
	var modelUser models.User
	//mencari data user dan merubah datanya
	cekUser := db.First(&modelUser, dataUser.Id)
	if cekUser.Error == nil {
		//enkripsi password dengan sha1
		var sha = sha1.New()
		sha.Write([]byte(dataUser.Password))
		var encrypted = sha.Sum(nil)
		var encryptedString = fmt.Sprintf("%x", encrypted)

		modelUser.Nama = dataUser.Nama
		modelUser.Username = dataUser.Username
		modelUser.Password = encryptedString

		hasil := db.Save(&modelUser)
		kesalahan := hasil.Error
		if hasil.Error == nil {
			c.JSON(http.StatusOK, gin.H{
				"status":    true,
				"pesan":     "Berhasil ubah data",
				"kesalahan": nil,
				"data":      modelUser,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status":    false,
				"pesan":     "Gagal ubah Data",
				"kesalahan": kesalahan.Error(),
				"data":      modelUser,
			})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Data Tidak ditemukan",
			"kesalahan": cekUser.Error.Error(),
			"data":      modelUser,
		})
	}
}

func UserHapus(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dataUser StrukturUserHapus

	if err := c.ShouldBindJSON(&dataUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca Data",
			"kesalahan": err.Error(),
		})
		return
	}

	//membuat variabel model user
	var modelUser models.User
	//mencari data user dan menghapus datanya
	cekUser := db.Delete(&modelUser, dataUser.Id)
	if cekUser.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil hapus data",
			"kesalahan": nil,
			"data":      dataUser,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Data Tidak ditemukan",
			"kesalahan": cekUser.Error.Error(),
			"data":      modelUser,
		})
	}
}
