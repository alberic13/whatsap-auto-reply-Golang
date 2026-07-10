package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sopingi.com/fikom/models"
)

type StrukturInformasi struct {
	Id         uint
	Judul      string `binding:"required"`
	Konten     string `binding:"required"`
	UrlDokumen string
}

func InformasiTampil(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var modelInformasi []models.Informasi
	hasil := db.Find(&modelInformasi)
	kesalahan := hasil.Error

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil Tampil data informasi",
			"kesalahan": nil,
			"data":      modelInformasi,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal Tampil data informasi",
			"kesalahan": kesalahan.Error(),
			"data":      nil,
		})
	}
}

func InformasiTambah(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dataInformasi StrukturInformasi

	if err := c.ShouldBindJSON(&dataInformasi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca Data",
			"kesalahan": err.Error(),
		})
		return
	}

	modelInformasi := models.Informasi{
		Judul:      dataInformasi.Judul,
		Konten:     dataInformasi.Konten,
		UrlDokumen: dataInformasi.UrlDokumen,
	}

	hasil := db.Create(&modelInformasi)
	kesalahan := hasil.Error
	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil tambah data informasi",
			"kesalahan": nil,
			"data":      modelInformasi,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal tambah data informasi",
			"kesalahan": kesalahan.Error(),
			"data":      modelInformasi,
		})
	}
}

func InformasiUbah(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dataInformasi StrukturInformasi

	if err := c.ShouldBindJSON(&dataInformasi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca Data",
			"kesalahan": err.Error(),
		})
		return
	}

	if dataInformasi.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Id wajib diisi",
			"kesalahan": "id tidak boleh 0",
		})
		return
	}

	var modelInformasi models.Informasi
	hasilCari := db.First(&modelInformasi, dataInformasi.Id)
	if hasilCari.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Data tidak ditemukan",
			"kesalahan": hasilCari.Error.Error(),
			"data":      nil,
		})
		return
	}

	modelInformasi.Judul = dataInformasi.Judul
	modelInformasi.Konten = dataInformasi.Konten
	modelInformasi.UrlDokumen = dataInformasi.UrlDokumen

	hasil := db.Save(&modelInformasi)
	kesalahan := hasil.Error
	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil ubah data informasi",
			"kesalahan": nil,
			"data":      modelInformasi,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal ubah data informasi",
			"kesalahan": kesalahan.Error(),
			"data":      modelInformasi,
		})
	}
}

func InformasiHapus(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dataInformasi StrukturInformasi

	if err := c.ShouldBindJSON(&dataInformasi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca Data",
			"kesalahan": err.Error(),
		})
		return
	}

	if dataInformasi.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Id wajib diisi",
			"kesalahan": "id tidak boleh 0",
		})
		return
	}

	var modelInformasi models.Informasi
	hasil := db.Delete(&modelInformasi, dataInformasi.Id)
	kesalahan := hasil.Error

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil hapus data informasi",
			"kesalahan": nil,
			"data":      dataInformasi,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal hapus data informasi",
			"kesalahan": kesalahan.Error(),
			"data":      dataInformasi,
		})
	}
}
