package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sopingi.com/fikom/models"
)

type StrukturPesan struct {
	Kode    string `binding:"required"`
	Balasan string `binding:"required"`
}

func PesanTampil(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var modelPesan []models.Pesan
	hasil := db.Find(&modelPesan)
	kesalahan := hasil.Error

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil Tampil data pesan",
			"kesalahan": nil,
			"data":      modelPesan,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal Tampil data pesan",
			"kesalahan": kesalahan.Error(),
			"data":      nil,
		})
	}
}

func PesanTambah(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dataPesan StrukturPesan

	if err := c.ShouldBindJSON(&dataPesan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca Data",
			"kesalahan": err.Error(),
		})
		return
	}

	modelPesan := models.Pesan{
		Kode:      dataPesan.Kode,
		Balasan:   dataPesan.Balasan,
		CreatedAt: time.Now(),
	}

	hasil := db.Create(&modelPesan)
	kesalahan := hasil.Error
	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil tambah data pesan",
			"kesalahan": nil,
			"data":      modelPesan,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal tambah data pesan",
			"kesalahan": kesalahan.Error(),
			"data":      modelPesan,
		})
	}
}

func PesanUbah(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dataPesan StrukturPesan

	if err := c.ShouldBindJSON(&dataPesan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca Data",
			"kesalahan": err.Error(),
		})
		return
	}

	var modelPesan models.Pesan
	hasilCari := db.First(&modelPesan, "kode = ?", dataPesan.Kode)
	if hasilCari.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Data tidak ditemukan",
			"kesalahan": hasilCari.Error.Error(),
			"data":      nil,
		})
		return
	}

	modelPesan.Balasan = dataPesan.Balasan
	modelPesan.UpdatedAt = time.Now()

	hasil := db.Save(&modelPesan)
	kesalahan := hasil.Error
	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil ubah data pesan",
			"kesalahan": nil,
			"data":      modelPesan,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal ubah data pesan",
			"kesalahan": kesalahan.Error(),
			"data":      modelPesan,
		})
	}
}

func PesanHapus(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dataPesan struct {
		Kode string `json:"kode" binding:"required"`
	}

	if err := c.ShouldBindJSON(&dataPesan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca Data",
			"kesalahan": err.Error(),
		})
		return
	}

	var modelPesan models.Pesan
	hasil := db.Delete(&modelPesan, "kode = ?", dataPesan.Kode)
	kesalahan := hasil.Error

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil hapus data pesan",
			"kesalahan": nil,
			"data":      dataPesan,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal hapus data pesan",
			"kesalahan": kesalahan.Error(),
			"data":      dataPesan,
		})
	}
}
