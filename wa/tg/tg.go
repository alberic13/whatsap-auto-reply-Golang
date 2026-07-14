package tg

import (
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"sopingi.com/fikom/ai"
	"sopingi.com/fikom/models"
)

var DB *gorm.DB

func KonekTelegram(db *gorm.DB) {
	DB = db
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Println("Peringatan: TELEGRAM_BOT_TOKEN tidak diset. Bot Telegram tidak aktif.")
		return
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic("Gagal menginisialisasi Telegram Bot:", err)
	}

	bot.Debug = false
	log.Printf("Telegram Bot terhubung sebagai @%s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // Abaikan update non-message
			continue
		}

		go handleMessage(bot, update.Message)
	}
}

func handleMessage(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	pesanRaw := msg.Text
	pesanClean := strings.TrimSpace(pesanRaw)
	pesanLower := strings.ToLower(pesanClean)
	// Membuat user ID unik untuk Telegram agar riwayat chat AI terpisah per user
	userID := fmt.Sprintf("tg_%d", msg.From.ID) 

	var balasan string

	// Kirim status "typing" ke Telegram
	action := tgbotapi.NewChatAction(chatID, tgbotapi.ChatTyping)
	bot.Send(action)

	// Logika routing pesan mirip dengan WhatsApp
	if strings.HasPrefix(pesanLower, "/ai") {
		// Hapus command "/ai" di depan
		pertanyaan := strings.TrimSpace(pesanClean[3:])
		if pertanyaan != "" {
			balasan = ai.TanyaAi(userID, pertanyaan)
		} else {
			balasan = "Masukkan pertanyaan setelah command /ai. Contoh: `/ai apa kabar?`"
		}
	} else if pesanLower == "tes" {
		balasan = "ada yang bisa di bantu ?"
	} else {
		// Cari dari database berdasarkan kode autoreply
		var pesanDB models.Pesan
		result := DB.Where("kode = ?", pesanLower).First(&pesanDB)
		if result.Error == nil {
			balasan = pesanDB.Balasan
		}
	}

	// Kirim pesan balasan jika ada isi
	if balasan != "" {
		reply := tgbotapi.NewMessage(chatID, balasan)
		_, err := bot.Send(reply)
		if err != nil {
			log.Println("Gagal mengirim balasan Telegram:", err)
		}
	}
}
