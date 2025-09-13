package main

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWT Claims
type Claims struct {
	UserID  uint   `json:"user_id"`
	Phone   string `json:"phone"`
	IsAdmin bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// Core Models
type Filial struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

type Category struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
	FilialID uint   `json:"filial_id"`
}

type Product struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	CategoryID uint   `json:"category_id"`
	Filials    []uint `json:"filials"`
}

type Order struct {
	ID         uint        `json:"id"`
	OrderID    string      `json:"order_id"`
	UserID     uint        `json:"user_id"`
	Username   string      `json:"username"`
	FilialID   uint        `json:"filial_id"`
	FilialName string      `json:"filial_name"`
	Items      []OrderItem `json:"items"`
	Total      float64     `json:"total"`
	Status     string      `json:"status"`
	Created    time.Time   `json:"created"`
	Updated    time.Time   `json:"updated"`
}

type OrderItem struct {
	ProductID uint    `json:"product_id"`
	Name      string  `json:"name"`
	Count     int     `json:"count"`
	Subtotal  float64 `json:"subtotal"`
}

// Request structs
type LoginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type RegisterUserRequest struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type AddFilialRequest struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

type UpdateFilialRequest struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

type AddCategoryRequest struct {
	Name string `json:"name"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name"`
}

type AddProductRequest struct {
	Name       string `json:"name"`
	CategoryID uint   `json:"category_id"`
	Filials    []uint `json:"filials"`
}

type UpdateProductRequest struct {
	Name       string `json:"name"`
	CategoryID uint   `json:"category_id"`
	Filials    []uint `json:"filials"`
}

type AssignFilialRequest struct {
	FilialID uint `json:"filial_id"`
}

type UpdateUserRequest struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	IsAdmin  bool   `json:"is_admin"`
	FilialID uint   `json:"filial_id"`
}

type CreateOrderRequest struct {
	// Filial   string            `json:"filial"`
	Items []CreateOrderItem `json:"items"`
}

type CreateOrderItem struct {
	ProductID uint `json:"product_id"`
	Count     int  `json:"count"`
}

type UpdateOrderRequest struct {
	Status string `json:"status"`
}

type PrinterRequest struct {
	Printer  string        `json:"printer"`
	Category string        `json:"category"`
	Username string        `json:"username"`
	OrderID  string        `json:"order_id"`
	Filial   string        `json:"filial"`
	Items    []PrinterItem `json:"items"`
}
type PrinterItem struct {
	Product string `json:"product"`
	Count   int    `json:"count"`
}

// Response structs
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  UserProfile `json:"user"`
}

type UserProfile struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	IsAdmin bool   `json:"is_admin"`
	Filial  Filial `json:"filial,omitempty"`
}

type GroupedProductsResponse struct {
	Success bool                       `json:"success"`
	Message string                     `json:"message"`
	Data    map[string][]ProductSimple `json:"data"`
}

type ProductSimple struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ProductDetails struct {
	ID           uint     `json:"id"`
	Name         string   `json:"name"`
	CategoryID   uint     `json:"category_id"`
	CategoryName string   `json:"category_name"`
	Filials      []uint   `json:"filials"`
	FilialNames  []string `json:"filial_names"`
}
