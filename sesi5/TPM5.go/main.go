package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

// Product represents the model for a product
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// ProductInput represents the input structure for creating products
type ProductInput struct {
	Status string     `json:"status"`
	Data   []*Product `json:"data"`
}

// ProductUpdateData represents the data structure for updating a product
type ProductUpdateData struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var (
	products   = make(map[int]*Product) // Menyimpan produk dalam map dengan key ID produk
	idCounter  = 1                      // Inisialisasi counter ID produk
	productsMu sync.Mutex               // Mutex untuk menghindari race condition
)

func main() {
	http.HandleFunc("/products", productsHandler) // Menangani request ke endpoint /products
	http.HandleFunc("/products/", productHandler) // Menangani request ke endpoint /products/{id}

	http.ListenAndServe(":8080", nil) // Menjalankan server di port 8080
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getProducts(w, r) // Menangani GET request untuk mendapatkan semua produk
	case http.MethodPost:
		createProduct(w, r) // Menangani POST request untuk membuat produk baru
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed) // Menangani metode yang tidak diizinkan
	}
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/products/"):]) // Mengambil ID dari URL
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getProduct(w, r, id) // Menangani GET request untuk mendapatkan produk berdasarkan ID
	case http.MethodPut:
		updateProduct(w, r, id) // Menangani PUT request untuk memperbarui produk berdasarkan ID
	case http.MethodDelete:
		deleteProduct(w, r, id) // Menangani DELETE request untuk menghapus produk berdasarkan ID
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed) // Menangani metode yang tidak diizinkan
	}
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	productsMu.Lock()         // Mengunci mutex untuk menghindari race condition
	defer productsMu.Unlock() // Membuka kunci mutex setelah fungsi selesai

	productList := []*Product{}
	for _, product := range products {
		productList = append(productList, product)
	}

	w.Header().Set("Content-Type", "application/json") // Mengatur header respons ke JSON
	json.NewEncoder(w).Encode(productList)             // Mengencode dan mengirim data produk dalam format JSON
}

func getProduct(w http.ResponseWriter, r *http.Request, id int) {
	productsMu.Lock()         // Mengunci mutex untuk menghindari race condition
	defer productsMu.Unlock() // Membuka kunci mutex setelah fungsi selesai

	product, ok := products[id]
	if !ok {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json") // Mengatur header respons ke JSON
	json.NewEncoder(w).Encode(product)                 // Mengencode dan mengirim data produk dalam format JSON
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	productsMu.Lock()         // Mengunci mutex untuk menghindari race condition
	defer productsMu.Unlock() // Membuka kunci mutex setelah fungsi selesai

	var input ProductInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	for _, product := range input.Data {
		product.ID = idCounter
		idCounter++
		products[product.ID] = product
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(input)
}

func updateProduct(w http.ResponseWriter, r *http.Request, id int) {
	productsMu.Lock()         // Mengunci mutex untuk menghindari race condition
	defer productsMu.Unlock() // Membuka kunci mutex setelah fungsi selesai

	product, ok := products[id]
	if !ok {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	var input struct {
		Status string              `json:"status"`
		Data   []ProductUpdateData `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if len(input.Data) > 0 {
		updateData := input.Data[0]
		product.Name = updateData.Name
		product.Price = updateData.Price
	}

	products[id] = product // Memperbarui produk di map

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product) // Mengencode dan mengirim data produk dalam format JSON
}

func deleteProduct(w http.ResponseWriter, r *http.Request, id int) {
	productsMu.Lock()         // Mengunci mutex untuk menghindari race condition
	defer productsMu.Unlock() // Membuka kunci mutex setelah fungsi selesai

	if _, ok := products[id]; !ok {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	delete(products, id) // Menghapus produk dari map

	w.WriteHeader(http.StatusNoContent) // Mengatur status code respons ke 204 No Content
}
