package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("🚀 Server ishga tushmoqda...")
	initData()

	r := mux.NewRouter()

	// CORS middleware ni eng birinchi qo'shish
	r.Use(corsMiddleware)

	// Health check
	r.HandleFunc("/", healthCheck).Methods("GET", "OPTIONS")
	r.HandleFunc("/health", healthCheck).Methods("GET", "OPTIONS")

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// ================= AUTH ROUTES =================
	api.HandleFunc("/login", login).Methods("POST", "OPTIONS")
	api.HandleFunc("/register", register).Methods("POST", "OPTIONS")

	// ================= USER ROUTES =================
	// User uchun (token kerak)
	api.HandleFunc("/products", authenticateJWT(getProductsHandler)).Methods("GET", "OPTIONS")
	api.HandleFunc("/orders", authenticateJWT(createOrderHandler)).Methods("POST", "OPTIONS")
	api.HandleFunc("/orders", authenticateJWT(getOrdersHandler)).Methods("GET", "OPTIONS")
	api.HandleFunc("/orders/{id:[0-9]+}", authenticateJWT(getOrderHandler)).Methods("GET", "OPTIONS")

	// Public endpoints (token kerak lekin admin emas)
	api.HandleFunc("/filials", authenticateJWT(getFilialsHandler)).Methods("GET", "OPTIONS")
	api.HandleFunc("/categories", authenticateJWT(getCategoriesHandler)).Methods("GET", "OPTIONS")

	// ================= ADMIN ROUTES =================
	// Filials management
	api.HandleFunc("/filials", requireAdmin(addFilialHandler)).Methods("POST", "OPTIONS")
	api.HandleFunc("/filials/{id:[0-9]+}", requireAdmin(getFilialHandler)).Methods("GET", "OPTIONS")
	api.HandleFunc("/filials/{id:[0-9]+}", requireAdmin(updateFilialHandler)).Methods("PUT", "OPTIONS")
	api.HandleFunc("/filials/{id:[0-9]+}", requireAdmin(deleteFilialHandler)).Methods("DELETE", "OPTIONS")

	// Categories management
	api.HandleFunc("/categories", requireAdmin(addCategoryHandler)).Methods("POST", "OPTIONS")
	api.HandleFunc("/categories/{id:[0-9]+}", requireAdmin(getCategoryHandler)).Methods("GET", "OPTIONS")
	api.HandleFunc("/categories/{id:[0-9]+}", requireAdmin(updateCategoryHandler)).Methods("PUT", "OPTIONS")
	api.HandleFunc("/categories/{id:[0-9]+}", requireAdmin(deleteCategoryHandler)).Methods("DELETE", "OPTIONS")

	// Products management - MUHIM: /products/all ni /products dan oldin qo'yish
	api.HandleFunc("/products/all", requireAdmin(getAllProductsHandler)).Methods("GET", "OPTIONS")
	api.HandleFunc("/products", requireAdmin(addProductHandler)).Methods("POST", "OPTIONS")
	api.HandleFunc("/products/{id:[0-9]+}", requireAdmin(getProductHandler)).Methods("GET", "OPTIONS")
	api.HandleFunc("/products/{id:[0-9]+}", requireAdmin(updateProductHandler)).Methods("PUT", "OPTIONS")
	api.HandleFunc("/products/{id:[0-9]+}", requireAdmin(deleteProductHandler)).Methods("DELETE", "OPTIONS")

	// Users management
	api.HandleFunc("/users", requireAdmin(getUsersHandler)).Methods("GET", "OPTIONS")
	api.HandleFunc("/users/{id:[0-9]+}", requireAdmin(getUserHandler)).Methods("GET", "OPTIONS")
	api.HandleFunc("/users/{id:[0-9]+}", requireAdmin(updateUserHandler)).Methods("PUT", "OPTIONS")
	api.HandleFunc("/users/{id:[0-9]+}", requireAdmin(deleteUserHandler)).Methods("DELETE", "OPTIONS")
	api.HandleFunc("/users/{id:[0-9]+}/assign-filial", requireAdmin(assignFilialHandler)).Methods("PUT", "OPTIONS")

	// Orders management
	api.HandleFunc("/orderslist", requireAdmin(getOrdersListHandler)).Methods("GET", "OPTIONS")
	api.HandleFunc("/orders/{id:[0-9]+}", requireAdmin(updateOrderHandler)).Methods("PUT", "OPTIONS")
	api.HandleFunc("/orders/{id:[0-9]+}", requireAdmin(deleteOrderHandler)).Methods("DELETE", "OPTIONS")

	fmt.Println("✅ Server ishga tushdi!")
	fmt.Println("📍 URL: http://localhost:1010")
	fmt.Println("📋 Health check: http://localhost:1010/health")
	fmt.Println("🔐 API Base URL: http://localhost:1010/api")
	fmt.Println("📱 CORS: Barcha domenlar uchun ochiq")
	fmt.Println("")
	fmt.Println("🔑 Auth Endpoints:")
	fmt.Println("   POST /api/login")
	fmt.Println("   POST /api/register")
	fmt.Println("")
	fmt.Println("👤 User Endpoints:")
	fmt.Println("   GET  /api/products")
	fmt.Println("   GET  /api/orders")
	fmt.Println("   POST /api/orders")
	fmt.Println("   GET  /api/orders/{id}")
	fmt.Println("   GET  /api/filials")
	fmt.Println("   GET  /api/categories")
	fmt.Println("")
	fmt.Println("⚡ Admin Endpoints:")
	fmt.Println("   📍 Filials: GET, POST, PUT, DELETE /api/filials")
	fmt.Println("   📂 Categories: GET, POST, PUT, DELETE /api/categories")
	fmt.Println("   📦 Products: GET, POST, PUT, DELETE /api/products")
	fmt.Println("   👥 Users: GET, POST, PUT, DELETE /api/users")
	fmt.Println("   📋 Orders: GET, PUT, DELETE /api/orders")
	fmt.Println("   📊 Orders List: GET /api/orderslist")

	log.Fatal(http.ListenAndServe(":1010", r))
}
