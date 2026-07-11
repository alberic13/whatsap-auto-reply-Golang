package ai

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestGemini(t *testing.T) {
	// Memuat file .env dari root folder proyek (naik satu direktori)
	err := godotenv.Load("../.env")
	if err != nil {
		t.Log("Peringatan: file .env tidak ditemukan, menggunakan environment variable bawaan OS.")
	}

	// Mengambil API Key dari env
	apiKey := os.Getenv("GOOGLE_API_KEY")
	t.Logf("Menggunakan API Key: %s", apiKey)

	// Inisialisasi AI
	InitAi()

	if geminiClient == nil {
		t.Fatal("Gagal menginisialisasi Gemini Client. Pastikan GOOGLE_API_KEY sudah diset di .env")
	}

	// Mencoba melakukan tanya jawab dengan Gemini menggunakan fungsi TanyaAi asli
	pertanyaan := "Halo, siapa kamu? Jawab dengan sangat singkat."
	t.Logf("Mengirim pertanyaan ke Gemini melalui TanyaAi: %s", pertanyaan)
	
	jawaban := TanyaAi("user_test_123", pertanyaan)
	t.Logf("Jawaban dari Gemini: %s", jawaban)

	// Memeriksa jika ada kegagalan respon
	if jawaban == "Mohon maaf, AI belum dikonfigurasi. Silakan set GOOGLE_API_KEY di file .env." ||
		jawaban == "Mohon maaf, terjadi gangguan saat memproses jawaban." ||
		jawaban == "Mohon maaf, AI tidak memberikan respon" {
		t.Errorf("Gagal mendapatkan respon yang benar dari Gemini. Error/Pesan: %s", jawaban)
	}
}
