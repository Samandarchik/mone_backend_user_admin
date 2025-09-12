package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Send order to printer
func sendToPrinter(order *Order) error {
	user := findUserByID(order.UserID)
	if user == nil {
		return fmt.Errorf("user topilmadi")
	}

	filial := findFilialByID(order.FilialID)
	if filial == nil {
		return fmt.Errorf("filial topilmadi")
	}

	// Mahsulotlarni kategoriya bo'yicha guruhlash
	categoryItems := make(map[uint][]PrinterItem)

	for _, item := range order.Items {
		product := findProductByID(item.ProductID)
		if product != nil {
			categoryItems[product.CategoryID] = append(categoryItems[product.CategoryID], PrinterItem{
				Product: item.Name,
				Count:   item.Count,
			})
		}
	}

	allSuccess := true

	// Har bir kategoriya uchun alohida chek yuborish
	for categoryID, items := range categoryItems {

		printRequest := PrinterRequest{
			Printer:  "p1",
			OrderID: order.OrderID,
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

		resp, err := http.Post("http://192.168.1.12:5000/print", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Printer ga yuborishda xato (%s): %v", "http://192.168.1.12:5000/print", err)
			allSuccess = false
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			log.Printf("✅ Chek yuborildi: %s (Kategoriya %d) - %s (%s)", "http://192.168.1.12:5000/print", categoryID, order.Username, order.FilialName)
		} else {
			log.Printf("❌ Chek yuborishda xato: %s - Status: %d", "http://192.168.1.12:5000/print", resp.StatusCode)
			allSuccess = false
		}
	}

	if !allSuccess {
		return fmt.Errorf("ba'zi cheklar yuborilmadi")
	}

	return nil
}
