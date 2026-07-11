package ai

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var geminiClient *genai.Client
var userHistories = make(map[string]*genai.ChatSession)
var mu sync.Mutex

func InitAi() {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	
	if apiKey == "" {
		fmt.Println("Warning: GOOGLE_API_KEY not set. AI features will not work.")
		return
	}
	
	ctx := context.Background()
	var err error
	geminiClient, err = genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		fmt.Printf("Error initializing Gemini client: %v\n", err)
		return
	}
	
	fmt.Println("AI Engine berhasil diinisialisasi.")
}

func TanyaAi(userID string, userInput string) string {
	if geminiClient == nil {
		return "Mohon maaf, AI belum dikonfigurasi. Silakan set GOOGLE_API_KEY di file .env."
	}
	
	ctx := context.Background()
	model := geminiClient.GenerativeModel("gemini-3.5-flash")
	
	// System instruction
	systemInstruction := "Anda adalah ZaldeBot, asisten virtual zalde digital solution IT tech. Jawab pertanyaan dengan singkat, ramah, dan informatif dalam Bahasa Indonesia."
	model.SystemInstruction = genai.NewUserContent(genai.Text(systemInstruction))
	
	// Ambil histori milik user ini dengan aman (Lock)
	mu.Lock()
	cs, exists := userHistories[userID]
	if !exists {
		cs = model.StartChat()
		userHistories[userID] = cs
	}
	mu.Unlock()
	
	// Kirim pesan user ke Gemini
	resp, err := cs.SendMessage(ctx, genai.Text(userInput))
	if err != nil {
		// Jika error dengan history, coba tanpa history
		mu.Lock()
		cs = model.StartChat()
		userHistories[userID] = cs
		mu.Unlock()
		
		resp, err = cs.SendMessage(ctx, genai.Text(userInput))
		if err != nil {
			return "Mohon maaf, terjadi gangguan saat memproses jawaban."
		}
	}
	
	// Extract response text
	var jawabanAi string
	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil && len(resp.Candidates[0].Content.Parts) > 0 {
		if txt, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
			jawabanAi = string(txt)
		} else {
			jawabanAi = fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
		}
	} else {
		jawabanAi = "Mohon maaf, AI tidak memberikan respon"
	}
	
	// Limit history size to last 20 messages
	mu.Lock()
	if len(cs.History) > 20 {
		cs.History = cs.History[len(cs.History)-20:]
	}
	mu.Unlock()
	
	return jawabanAi
}
