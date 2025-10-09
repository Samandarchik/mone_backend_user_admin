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
func sendToTelegram(order *Order, printerItems map[uint][]PrinterItem, printerCategories map[uint]*Category, printerSuccess bool) error {
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
		category := printerCategories[printerID]
		if category != nil {
			message.WriteString(fmt.Sprintf("\n🔸 *%s (Printer: %d):*\n", category.Name, printerID))
			for _, item := range items {
				message.WriteString(fmt.Sprintf("   • %s", item.Product))
				message.WriteString(fmt.Sprintf(" - %.3f %s\n", float64(item.Count), item.Type))
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

// Send order to printer (PRINTER ID bilan guruhlab yuborish)
func sendToPrinter(order *Order) error {
	user := findUserByID(order.UserID)
	if user == nil {
		return fmt.Errorf("user topilmadi")
	}

	filial := findFilialByID(order.FilialID)
	if filial == nil {
		return fmt.Errorf("filial topilmadi")
	}

	// Mahsulotlarni PRINTER ID bilan guruhlash
	printerItems := make(map[uint][]PrinterItem)  // map[PrinterID][]Items
	printerCategories := make(map[uint]*Category) // map[PrinterID]Category
	allSuccess := true

	// 1-QADAM: Har bir order item uchun
	for _, item := range order.Items {
		// Product ni ID bilan toping
		product := findProductByID(item.ProductID)
		if product == nil {
			log.Printf("Mahsulot topilmadi: %d", item.ProductID)
			continue
		}

		// Product ichidagi CategoryID bilan Category ni toping
		category := GetCategoryByID(product.CategoryID)
		if category == nil {
			log.Printf("Kategoriya topilmadi: %d", product.CategoryID)
			allSuccess = false
			continue
		}

		// Category ichidagi Printer field ni olib printerID qiling
		printerID := category.Printer

		// Mahsulotni Printer ID bilan guruhlab qo'ying
		printerItems[printerID] = append(printerItems[printerID], PrinterItem{
			Product: item.Name,
			Count:   item.Count,
			Type:    item.Type,
		})

		// Printer kategoriyasini saqlang (keyinroq Telegram uchun)
		printerCategories[printerID] = category
	}

	// 2-QADAM: Har bir PRINTER uchun alohida chek yuborish
	for printerID, items := range printerItems {
		category := printerCategories[printerID]

		// Printer uchun request tayyorlash
		printRequest := PrinterRequest{
			Printer:  printerID,
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

		// Printer serveriga yuborish
		resp, err := http.Post("https://marxabo1.javohir-jasmina.uz/print", "application/json", bytes.NewBuffer(jsonData))
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

	// 3-QADAM: Telegram ga yuborish
	if len(printerItems) > 0 {
		err := sendToTelegram(order, printerItems, printerCategories, allSuccess)
		if err != nil {
			log.Printf("Telegram ga yuborishda xato: %v", err)
			// Printer yuborilgan bo'lsa, xato qayd qilamiz lekin fail qilmaymiz
		}
	}

	if !allSuccess {
		return fmt.Errorf("ba'zi cheklar yuborilmadi")
	}

	return nil
}

// Helper funksiya - Category ni ID bilan topish
func GetCategoryByID(id uint) *Category {
	return findCategoryByID(id)
}

func findCategoryByID(id uint) *Category {
	for i, c := range categories {
		if c.ID == id {
			return &categories[i]
		}
	}
	return nil
}
