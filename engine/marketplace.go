package engine

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

// Product, JSON dosyasındaki ürün yapısını temsil eder
type Product struct {
	UrunId         int64   `json:"UrunId"`
	Urun           string  `json:"Urun"`
	KtgAdi         string  `json:"KtgAdi"`
	MarkaAd        string  `json:"MarkaAd"`
	UrunAciklamasi string  `json:"UrunAciklamasi"`
	Price          float64 `json:"Price"` // Ürün fiyatı
}

// Marketplace, pazar yerindeki ürün havuzunu yönetir
type Marketplace struct {
	Products []Product
}

// NewMarketplace, JSON dosyasından ürünleri belleğe yükler ve fiyatlandırır
func NewMarketplace(jsonPath string) (*Marketplace, error) {
	fmt.Println("📦 Pazar yeri verileri yükleniyor...")
	file, err := os.Open(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("veri dosyası açılamadı: %v", err)
	}
	defer file.Close()

	var products []Product
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&products); err != nil {
		return nil, fmt.Errorf("JSON çözümlenemedi: %v", err)
	}

	// Ürünlere rastgele fiyatlar atayalım (50 TL - 1050 TL arası)
	for i := range products {
		products[i].Price = float64(rand.Intn(1000)+50) + 0.99
	}

	fmt.Printf("✅ Pazar yeri hazır! Toplam Ürün Sayısı: %d (Fiyatlar belirlendi)\n", len(products))
	return &Marketplace{Products: products}, nil
}