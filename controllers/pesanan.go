package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sopingi.com/fikom/models"
)

// Binding dari POST/PUT JSON
type StrukturPesanan struct {
	Id            uint    `json:"id"`
	NamaProduk    string  `json:"nama_produk" binding:"required"`
	Harga         float64 `json:"harga" binding:"required"`
	Jumlah        int     `json:"jumlah" binding:"required"`
	Diskon        float64 `json:"diskon"`
	TotalHarga    float64 `json:"total_harga"`
	TipePelanggan string  `json:"tipe_pelanggan"`
}

func hitungPesanan(harga float64, jumlah int) (float64, float64, float64, string) {
	totalKotor := harga * float64(jumlah)
	var diskon float64
	tipePelanggan := "Regular"

	if totalKotor > 500000 {
		diskon = totalKotor * 0.10
		tipePelanggan = "Gold"
	}

	totalBayar := totalKotor - diskon
	return totalKotor, diskon, totalBayar, tipePelanggan
}

func PesananTampil(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var modelPesanan []models.Pesanan
	hasil := db.Find(&modelPesanan)
	kesalahan := hasil.Error

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil tampil data pesanan",
			"kesalahan": nil,
			"data":      modelPesanan,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal tampil data pesanan",
			"kesalahan": kesalahan.Error(),
			"data":      nil,
		})
	}
}

func PesananTambah(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dataPesanan StrukturPesanan

	if err := c.ShouldBindJSON(&dataPesanan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca data",
			"kesalahan": err.Error(),
		})
		return
	}

	totalKotor, diskon, totalBayar, tipePelanggan := hitungPesanan(dataPesanan.Harga, dataPesanan.Jumlah)

	modelPesanan := models.Pesanan{
		NamaProduk:    dataPesanan.NamaProduk,
		Harga:         dataPesanan.Harga,
		Jumlah:        dataPesanan.Jumlah,
		Diskon:        diskon,
		TotalHarga:    totalBayar,
		TipePelanggan: tipePelanggan,
	}

	hasil := db.Create(&modelPesanan)
	kesalahan := hasil.Error
	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil tambah data pesanan",
			"kesalahan": nil,
			"data": gin.H{
				"id":             modelPesanan.ID,
				"nama_produk":    modelPesanan.NamaProduk,
				"harga":          modelPesanan.Harga,
				"jumlah":         modelPesanan.Jumlah,
				"total_kotor":    totalKotor,
				"diskon":         modelPesanan.Diskon,
				"total_harga":    modelPesanan.TotalHarga,
				"tipe_pelanggan": modelPesanan.TipePelanggan,
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal tambah data pesanan",
			"kesalahan": kesalahan.Error(),
			"data":      modelPesanan,
		})
	}
}

func PesananUbah(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dataPesanan StrukturPesanan

	if err := c.ShouldBindJSON(&dataPesanan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca data",
			"kesalahan": err.Error(),
		})
		return
	}

	if dataPesanan.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Id wajib diisi",
			"kesalahan": "id tidak boleh 0",
		})
		return
	}

	var modelPesanan models.Pesanan
	hasilCari := db.First(&modelPesanan, dataPesanan.Id)
	if hasilCari.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Data tidak ditemukan",
			"kesalahan": hasilCari.Error.Error(),
			"data":      nil,
		})
		return
	}

	_, diskon, totalBayar, tipePelanggan := hitungPesanan(dataPesanan.Harga, dataPesanan.Jumlah)

	modelPesanan.NamaProduk = dataPesanan.NamaProduk
	modelPesanan.Harga = dataPesanan.Harga
	modelPesanan.Jumlah = dataPesanan.Jumlah
	modelPesanan.Diskon = diskon
	modelPesanan.TotalHarga = totalBayar
	modelPesanan.TipePelanggan = tipePelanggan

	hasil := db.Save(&modelPesanan)
	kesalahan := hasil.Error
	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil ubah data pesanan",
			"kesalahan": nil,
			"data":      modelPesanan,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal ubah data pesanan",
			"kesalahan": kesalahan.Error(),
			"data":      modelPesanan,
		})
	}
}

func PesananHapus(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dataPesanan StrukturPesanan

	if err := c.ShouldBindJSON(&dataPesanan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca data",
			"kesalahan": err.Error(),
		})
		return
	}

	if dataPesanan.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Id wajib diisi",
			"kesalahan": "id tidak boleh 0",
		})
		return
	}

	var modelPesanan models.Pesanan
	hasil := db.Delete(&modelPesanan, dataPesanan.Id)
	kesalahan := hasil.Error

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil hapus data pesanan",
			"kesalahan": nil,
			"data":      dataPesanan,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal hapus data pesanan",
			"kesalahan": kesalahan.Error(),
			"data":      dataPesanan,
		})
	}
}
