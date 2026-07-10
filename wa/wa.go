package wa

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
	_ "modernc.org/sqlite"

	"gorm.io/gorm"
	"sopingi.com/fikom/ai"
	"sopingi.com/fikom/models"
)

// var untuk client whatssap
var clientWa *whatsmeow.Client

// array id wa
var recipientIDs []string

// var untuk database
var DB *gorm.DB

func connectWhatsApp(ctx context.Context, container *sqlstore.Container, clientLog waLog.Logger) error {
	for {
		deviceStore, err := container.GetFirstDevice(ctx)
		if err != nil {
			return err
		}

		if deviceStore != nil {
			deviceStore.Platform = "macOS"
		}

		clientWa = whatsmeow.NewClient(deviceStore, clientLog)
		clientWa.AddEventHandler(eventHandler)

		if clientWa.Store.ID == nil {
			qrChan, _ := clientWa.GetQRChannel(ctx)
			err = clientWa.Connect()
			if err != nil {
				if strings.Contains(strings.ToLower(err.Error()), "logged out from another device") || strings.Contains(err.Error(), "401") {
					fmt.Println("Session lama sudah logout. Scan QR code untuk login ulang.")
					_ = deviceStore.Delete(ctx)
					continue
				}
				return err
			}
			for evt := range qrChan {
				if evt.Event == "code" {
					fmt.Println("Scan this QR code with WhatsApp:")
					qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				} else {
					fmt.Println("Login event:", evt.Event)
				}
			}
			if err = clientWa.SendPresence(ctx, types.PresenceAvailable); err != nil {
				fmt.Println("Peringatan: Gagal mengirim presence:", err)
			}
			return nil
		}

		err = clientWa.Connect()
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "logged out from another device") || strings.Contains(err.Error(), "401") {
				fmt.Println("Session lama sudah logout. Scan QR code untuk login ulang.")
				_ = deviceStore.Delete(ctx)
				continue
			}
			return err
		}
		if err = clientWa.SendPresence(ctx, types.PresenceAvailable); err != nil {
			fmt.Println("Peringatan: Gagal mengirim presence:", err)
		}
		return nil
	}
}

func cleanRecipientJID(recipient string) (types.JID, error) {
	recipient = strings.TrimSpace(recipient)
	if recipient == "" {
		return types.JID{}, fmt.Errorf("nomor penerima tidak boleh kosong")
	}

	parts := strings.Split(recipient, "@")
	phone := parts[0]
	server := "s.whatsapp.net"
	if len(parts) > 1 {
		server = parts[1]
	}

	var cleanPhone strings.Builder
	for _, r := range phone {
		if r >= '0' && r <= '9' {
			cleanPhone.WriteRune(r)
		}
	}
	phoneStr := cleanPhone.String()

	if strings.HasPrefix(phoneStr, "0") {
		phoneStr = "62" + phoneStr[1:]
	} else if strings.HasPrefix(phoneStr, "8") {
		phoneStr = "62" + phoneStr
	}

	fullJID := phoneStr + "@" + server
	return types.ParseJID(fullJID)
}

func sendTextMessage(ctx context.Context, recipient string, messageText string) error {
	jid, err := cleanRecipientJID(recipient)
	if err != nil {
		return err
	}

	_, err = clientWa.SendMessage(ctx, jid, &waE2E.Message{
		Conversation: proto.String(messageText),
	})
	return err
}

func sendTextMessageToRecipients(ctx context.Context, recipients []string, messageText string) error {
	for _, recipient := range recipients {
		if err := sendTextMessage(ctx, recipient, messageText); err != nil {
			return err
		}
		fmt.Println("Pesan terkirim ke", recipient)
	}
	return nil
}

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		fmt.Println("Pesan diterima!", v.Message.GetConversation())

		//ketik disini
		fmt.Println(" => dari saya sendiri = ", v.Info.IsFromMe)
		fmt.Println(" => server = ", v.Info.MessageSource.Chat.Server)
		fmt.Println(" => apakah dari group = ", v.Info.IsGroup)
		fmt.Println(" => apakah dari broadcast = ", v.Info.IsIncomingBroadcast())

		if !v.Info.IsFromMe &&
			(v.Info.MessageSource.Chat.Server == "lid" || v.Info.MessageSource.Chat.Server == "s.whatsapp.net") &&
			!v.Info.IsGroup &&
			!v.Info.IsIncomingBroadcast() {
			go func(msg *events.Message) {
				fmt.Println("PENGIRIM = ", msg.Info.Sender.User)
				
				var pesan string
				//jika text biasa
				if msg.Message.GetConversation() != "" {
					pesan = msg.Message.GetConversation()
				//jika text ada format tertentu
				} else if msg.Message.ExtendedTextMessage != nil && msg.Message.ExtendedTextMessage.GetText() != "" {
					pesan = msg.Message.ExtendedTextMessage.GetText()
				}
				fmt.Println("PESAN = " + pesan)
				
				//membuat array id_pesan
				var id_pesan []types.MessageID
				id_pesan = append(id_pesan, msg.Info.ID)
				
				//status pesan dibaca
				clientWa.MarkRead(context.Background(), id_pesan, time.Now(), msg.Info.Chat, msg.Info.Sender)
				
				//pengirim berlangkanan untuk menerima notif
				clientWa.SubscribePresence(context.Background(), msg.Info.Sender)
				
				//notif online
				clientWa.SendPresence(context.Background(), types.PresenceAvailable)
				
				//jeda 2 detik, monggo boleh diubah
				time.Sleep(2 * time.Second)
				
				//notif mengetik
				clientWa.SendChatPresence(context.Background(), msg.Info.Sender, types.ChatPresenceComposing, types.ChatPresenceMediaText)
				
				//jeda 3 detik, monggo boleh diubah
				time.Sleep(3 * time.Second)
				
				//notif berhenti mengetik
				clientWa.SendChatPresence(context.Background(), msg.Info.Sender, types.ChatPresencePaused, types.ChatPresenceMediaText)
				
				//untuk uji coba balasan hanya utk pesan berisi "tes"
				pesanClean := strings.TrimSpace(pesan)
				pesanLower := strings.ToLower(pesanClean)
				if strings.HasPrefix(pesanLower, "[ai]") {
					pertanyaan := strings.TrimSpace(pesanClean[4:])
					if pertanyaan != "" {
						jawabanAi := ai.TanyaAi(msg.Info.Sender.User, pertanyaan)
						kirimPesanText(msg.Info.Sender, jawabanAi)
					} else {
						kirimPesanText(msg.Info.Sender, "Masukkan pertanyaan setelah prefiks [ai]. Contoh: *[ai]Selamat pagi*")
					}
				} else if pesanLower == "tes" {
					kirimPesan(msg.Info.Sender)
				} else {
					kirimPesanDatabase(msg.Info.Sender, pesanLower)
				}
			}(v)
		}

	case *events.Receipt:
		if v.Type == types.ReceiptTypeRead {
			fmt.Printf("Status pesan: dibaca oleh %s untuk message id %v\n", v.Chat.String(), v.MessageIDs)
		}
	}
}

func promptForList(prompt string) ([]string, error) {
	text, err := promptForText(prompt)
	if err != nil {
		return nil, err
	}

	rawItems := strings.Split(text, ",")
	items := make([]string, 0, len(rawItems))
	for _, item := range rawItems {
		item = strings.TrimSpace(item)
		if item != "" {
			items = append(items, item)
		}
	}
	return items, nil
}

func promptForText(prompt string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

func findDBPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "examplestore.db"
	}
	dir := cwd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return filepath.Join(dir, "examplestore.db")
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "examplestore.db"
}

func KonekWa(db *gorm.DB) {
	DB = db

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	ctx := context.Background()
	dbPath := findDBPath()
	dbURI := fmt.Sprintf("file:%s?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)", filepath.ToSlash(dbPath))
	container, err := sqlstore.New(ctx, "sqlite", dbURI, dbLog)
	if err != nil {
		panic(err)
	}
	defer container.Close()

	clientLog := waLog.Stdout("Client", "DEBUG", true)
	if err = connectWhatsApp(ctx, container, clientLog); err != nil {
		panic(err)
	}
	fmt.Println("WhatsApp Bot terhubung dan mendengarkan pesan...")

	client := clientWa
	clientWa = client

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}

func kirimPesan(IDPenerima types.JID) {
	clientWa.SendMessage(
		context.Background(),
		IDPenerima.ToNonAD(),
		&waE2E.Message{
			Conversation: proto.String("ada yang bisa di bantu ?"),
		},
	)
}

func kirimPesanText(IDPenerima types.JID, messageText string) {
	_, err := clientWa.SendMessage(context.Background(), IDPenerima.ToNonAD(), &waE2E.Message{
		Conversation: proto.String(messageText),
	})
	if err != nil {
		fmt.Println("Gagal mengirim balasan dari database:", err)
	}
}

func kirimPesanDatabase(IDPenerima types.JID, kode string) {
	var pesan models.Pesan
	// mencari berdasarkan primary key (Kode)
	result := DB.Where("kode = ?", kode).First(&pesan)
	if result.Error == nil {
		kirimPesanText(IDPenerima, pesan.Balasan)
	}
}
