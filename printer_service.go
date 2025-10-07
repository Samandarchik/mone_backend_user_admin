package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// Telegram message structure
type TelegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

// Send formatted message to Telegram
func sendToTelegram(order *Order, printerItems map[uint][]PrinterItem, printerSuccess bool) error {
	user := findUserByID(order.UserID)
	if user == nil {
		return fmt.Errorf("user topilmadi")
	}

	// Format the message
	var message strings.Builder
	message.WriteString("🧾 *НОВЫЙ ПОРЯДОК*\n\n")
	message.WriteString(fmt.Sprintf("📋 *Заказ ID:* `%s`\n", order.OrderID))
	message.WriteString(fmt.Sprintf("👤 *Клиент:* %s\n", order.Username))
	message.WriteString(fmt.Sprintf("🏢 *Ветвь:* %s\n", order.FilialName))
	loc := time.FixedZone("UTC+5", 5*60*60)
	tashkentTime := time.Now().In(loc)
	message.WriteString(fmt.Sprintf("⏰ *Время:* %s\n\n", tashkentTime.Format("2006-01-02 15:04:05")))

	// Printer status at the top
	printerStatusText := "❌ *Невозможно отправить на принтер* @Baxtiyor0055"
	if printerSuccess {
		printerStatusText = "✅ *Отправлено в типографию*"
	}
	message.WriteString(fmt.Sprintf("🖨️ *СТАТУС:* %s\n\n", printerStatusText))

	message.WriteString("📦 *ТОВАРЫ:*\n")

	// Display items grouped by printer
	for printerID, items := range printerItems {
		message.WriteString(fmt.Sprintf("\n🖨️ *Printer %d:*\n", printerID))
		for _, item := range items {
			message.WriteString(fmt.Sprintf("   • %s", item.Product))
			if float32(int(item.Count)) == item.Count {
				message.WriteString(fmt.Sprintf(" - %.0f %s\n", item.Count, item.Type))
			} else {
				message.WriteString(fmt.Sprintf(" - %.2f %s\n", item.Count, item.Type))
			}
		}
	}

	telegramMsg := TelegramMessage{
		ChatID:    "-4985547344",
		Text:      message.String(),
		ParseMode: "Markdown",
	}

	jsonData, err := json.Marshal(telegramMsg)
	if err != nil {
		log.Printf("Telegram JSON marshal xato: %v", err)
		return err
	}

	telegramURL := "https://api.telegram.org/bot8157743798:AAELzxyyFLSMxbT-XL4l-3ZVmxVBXYOY0Ro/sendMessage"

	resp, err := http.Post(telegramURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Telegram ga yuborishda xato: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		log.Printf("✅ Telegram ga yuborildi: %s", order.OrderID)
	} else {
		log.Printf("❌ Telegram ga yuborishda xato - Status: %d", resp.StatusCode)
		return fmt.Errorf("telegram yuborishda xato: status %d", resp.StatusCode)
	}

	return nil
}

// Send order to printer (grouped by category.Printer)
func sendToPrinter(order *Order) error {
	user := findUserByID(order.UserID)
	if user == nil {
		return fmt.Errorf("user topilmadi")
	}

	filial := findFilialByID(order.FilialID)
	if filial == nil {
		return fmt.Errorf("filial topilmadi")
	}

	// Mahsulotlarni category.Printer bo'yicha guruhlash
	printerItems := make(map[uint][]PrinterItem)
	printerSuccess := false

	for _, item := range order.Items {
		product := findProductByID(item.ProductID)
		if product != nil {
			category := GetCategoryByID(product.CategoryID)
			if category != nil {
				// category.Printer bu printer ID
				printerID := uint(category.Printer)
				printerItems[printerID] = append(printerItems[printerID], PrinterItem{
					Product: item.Name,
					Count:   item.Count,
					Type:    item.Type,
				})
			}
		}
	}

	allSuccess := true

	// Har bir printer uchun alohida chek yuborish
	for printerID, items := range printerItems {
		printRequest := PrinterRequest{
			Printer:  fmt.Sprintf("%d", printerID),
			OrderID:  order.OrderID,
			Category: fmt.Sprintf("Printer %d", printerID),
			Username: order.Username,
			Filial:   order.FilialName,
			Items:    items,
		}

		jsonData, err := json.Marshal(printRequest)
		if err != nil {
			log.Printf("JSON marshal xato: %v", err)
			allSuccess = false
			continue
		}

		resp, err := http.Post("http://localhost:8080/print", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Printer ga yuborishda xato (Printer %d): %v", printerID, err)
			allSuccess = false
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			log.Printf("✅ Chek yuborildi: Printer %d - %s (%s)", printerID, order.Username, order.FilialName)
		} else {
			log.Printf("❌ Chek yuborishda xato: Printer %d - Status: %d", printerID, resp.StatusCode)
			allSuccess = false
		}
	}

	printerSuccess = allSuccess

	// Send to Telegram after processing all printers
	if len(printerItems) > 0 {
		err := sendToTelegram(order, printerItems, printerSuccess)
		if err != nil {
			log.Printf("Telegram ga yuborishda xato: %v", err)
		}
	}

	if !allSuccess {
		return fmt.Errorf("ba'zi cheklar yuborilmadi")
	}

	return nil
}
