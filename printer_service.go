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
		message.WriteString(fmt.Sprintf("   • %s - %2f %s\n", item.Product, item.Count, item.Type))
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

// Send order to printer (kategoriya ajratmasdan bitta ro'yxat qilib)
func sendToPrinter(order *Order) error {
	user := findUserByID(order.UserID)
	if user == nil {
		return fmt.Errorf("user topilmadi")
	}

	filial := findFilialByID(order.FilialID)
	if filial == nil {
		return fmt.Errorf("filial topilmadi")
	}

	// Barcha mahsulotlarni bitta ro'yxatga yig'ish
	var allItems []PrinterItem
	for _, item := range order.Items {
		printerItem := PrinterItem{
			Product: item.Name,
			Count:   item.Count,
			Type:    item.Type,
		}
		allItems = append(allItems, printerItem)
	}

	// Endi butun ro'yxatni bitta printer so'rovi sifatida yuboramiz
	printRequest := PrinterRequest{
		Printer:  "p1",
		OrderID:  order.OrderID,
		Category: "Barcha mahsulotlar", // umumiy nom berish mumkin
		Username: order.Username,
		Filial:   order.FilialName,
		Items:    allItems,
	}

	jsonData, err := json.Marshal(printRequest)
	if err != nil {
		log.Printf("JSON marshal xato: %v", err)
		return err
	}

	resp, err := http.Post("http://localhost:8080/print", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Printer ga yuborishda xato: %v", err)
		return err
	}
	defer resp.Body.Close()

	printerSuccess := false
	if resp.StatusCode == 200 {
		printerSuccess = true
		log.Printf("✅ Chek yuborildi: %s (%s - %s)", order.OrderID, order.Username, order.FilialName)
	} else {
		log.Printf("❌ Chek yuborishda xato - Status: %d", resp.StatusCode)
	}

	// Telegramga hamma mahsulotlar bilan bitta xabar yuboramiz
	err = sendToTelegram(order, allItems, printerSuccess)
	if err != nil {
		log.Printf("Telegramga yuborishda xato: %v", err)
	}

	if !printerSuccess {
		return fmt.Errorf("chek yuborilmadi, printer xato status qaytardi")
	}

	return nil
}

//https://marxabo1.javohir-jasmina.uz/print
//http://localhost:8080/print
