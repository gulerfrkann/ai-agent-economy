package engine

import (
	"encoding/json"
	"fmt"
	"os"
)

// Product, JSON dosyasındaki ürün yapısını temsil eder
type Product struct {
	UrunId           int64  `json:"UrunId"`
	Urun             string `json:"Urun"`
	KtgAdi           string `json:"KtgAdi"`
	MarkaAd          string `json:"MarkaAd"`
	UrunAciklamasi   string `json:"UrunAciklamasi"`
}

// Marketplace, pazar yerindeki ürün havuzunu yönetir
type Marketplace struct {
	Products []Product
}

// NewMarketplace, JSON dosyasından ürünleri belleğe yükler
func NewMarketplace(jsonPath string) (*Marketplace, error) {
	fmt.Println(" Pazar yeri verileri yükleniyor...")
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

	fmt.Printf(" Pazar yeri hazır! Toplam Ürün Sayısı: %d\n", len(products))
	return &Marketplace{Products: products}, nil
}