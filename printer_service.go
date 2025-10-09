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
func sendToTelegram(order *Order, allItems []PrinterItem, printerSuccess bool) error {
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

	// Display all items in one list (no category grouping)
	for _, item := range allItems {
		message.WriteString(fmt.Sprintf("   • %s - %d %s\n", item.Product, item.Count, item.Type))
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

// Send order to printer (updated version)
func sendToPrinter(order *Order) error {
	user := findUserByID(order.UserID)
	if user == nil {
		return fmt.Errorf("user topilmadi")
	}

	filial := findFilialByID(order.FilialID)
	if filial == nil {
		return fmt.Errorf("filial topilmadi")
	}

	// Barcha mahsulotlarni bitta ro'yxatga yig'ish (kategoriyasiz)
	var allItems []PrinterItem
	categoryItems := make(map[uint][]PrinterItem) // Printer uchun kategoriya bo'yicha kerak bo'lsa
	printerSuccess := false

	for _, item := range order.Items {
		product := findProductByID(item.ProductID)
		printerItem := PrinterItem{
			Product: item.Name,
			Count:   item.Count,
			Type:    item.Type,
		}

		// Telegram uchun barcha itemlar bitta ro'yxatda
		allItems = append(allItems, printerItem)

		// Printer uchun kategoriya bo'yicha (agar kerak bo'lsa)
		if product != nil {
			categoryItems[product.CategoryID] = append(categoryItems[product.CategoryID], printerItem)
		}
	}

	allSuccess := true

	// Har bir kategoriya uchun alohida chek yuborish (printer uchun)
	for categoryID, items := range categoryItems {
		category := GetCategoryByID(categoryID)
		if category == nil {
			log.Printf("Kategoriya topilmadi: %d", categoryID)
			allSuccess = false
			continue
		}

		printRequest := PrinterRequest{
			Printer:  "p1",
			OrderID:  order.OrderID,
			Category: category.Name,
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

		resp, err := http.Post("https://marxabo1.javohir-jasmina.uz/print", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Printer ga yuborishda xato (%s): %v", "https://marxabo1.javohir-jasmina.uz/print", err)
			allSuccess = false
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			log.Printf("✅ Chek yuborildi: %s (Kategoriya %d) - %s (%s)", "https://marxabo1.javohir-jasmina.uz/print", categoryID, order.Username, order.FilialName)
		} else {
			log.Printf("❌ Chek yuborishda xato: %s - Status: %d", "https://marxabo1.javohir-jasmina.uz/print", resp.StatusCode)
			allSuccess = false
		}
	}

	// Set overall printer success status
	printerSuccess = allSuccess

	// Send to Telegram with all items in one list
	if len(allItems) > 0 {
		err := sendToTelegram(order, allItems, printerSuccess)
		if err != nil {
			log.Printf("Telegram ga yuborishda xato: %v", err)
		}
	}

	if !allSuccess {
		return fmt.Errorf("ba'zi cheklar yuborilmadi")
	}

	return nil
}
