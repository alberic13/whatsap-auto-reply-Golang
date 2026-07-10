package ai

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/google/generative-ai-go/client"
	"github.com/google/generative-ai-go/option"
)

var geminiClient *client.Client
var userHistories = make(map[string]interface{})
var mu sync.Mutex

func InitAi() {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	
	if apiKey == "" {
		fmt.Println("Warning: GOOGLE_API_KEY not set. AI features will not work.")
		return
	}
	
	ctx := context.Background()
	var err error
	geminiClient, err = client.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		fmt.Printf("Error initializing Gemini client: %v\n", err)
		return
	}
	
	fmt.Println("AI Engine berhasil diinisialisasi.")
}

func TanyaAi(userID string, userInput string) string {
	if geminiClient == nil {
		return "Mohon maaf, AI belum dikonfigurasi. Silakan set GOOGLE_API_KEY environment variable."
	}
	
	ctx := context.Background()
	model := geminiClient.GenerativeModel("gemini-1.5-flash")
	
	// System instruction
	systemInstruction := "Anda adalah FikomBot, asisten virtual resmi Fakultas Ilmu Komputer UDB Surakarta. Jawab pertanyaan dengan singkat, ramah, dan informatif dalam Bahasa Indonesia."
	model.SystemInstruction = systemInstruction
	
	// Ambil histori milik user ini dengan aman (Lock)
	mu.Lock()
	userHistory, exists := userHistories[userID]
	var session interface{}
	if exists {
		session = userHistory
	}
	mu.Unlock()
	
	// Mulai session atau lanjutkan
	var cs *client.ChatSession
	if session != nil {
		cs = session.(*client.ChatSession)
	} else {
		cs = model.StartChat()
	}
	
	// Kirim pesan user ke Gemini
	resp, err := cs.SendMessage(ctx, cs.History[len(cs.History):len(cs.History)]...)
	if err != nil {
		// Jika error dengan history, coba tanpa history
		cs = model.StartChat()
		resp, err = cs.SendMessage(ctx, cs.History[len(cs.History):len(cs.History)]...)
		if err != nil {
			return "Mohon maaf, terjadi gangguan saat memproses jawaban."
		}
	}
	
	// Extract response text
	var jawabanAi string
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		jawabanAi = fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
	} else {
		jawabanAi = "Mohon maaf, AI tidak memberikan respon"
	}
	
	// Simpan session dan history
	mu.Lock()
	userHistories[userID] = cs
	// Limit history size to last 20 messages
	if len(cs.History) > 20 {
		cs.History = cs.History[len(cs.History)-20:]
	}
	mu.Unlock()
	
	return jawabanAi
}

