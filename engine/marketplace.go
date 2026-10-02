package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// Product JSON dosyasındaki gerçek alan adlarıyla birebir eşleştirildi
type Product struct {
	ID       int     `json:"UrunId"`
	Urun     string  `json:"Urun"`
	Category string  `json:"KtgAdi"`
	SubCat   string  `json:"Ktg2Adi"`
}

type MarketItem struct {
	ID       int     `json:"id"`
	SellerID int     `json:"seller_id"`
	Urun     string  `json:"urun"`
	Category string  `json:"kategori"`
	Price    float64 `json:"fiyat"`
}

type CategoryStat struct {
	CurrentPrice float64
	BasePrice    float64
	DemandCount  int
	SupplyCount  int
}

type Marketplace struct {
	mu         sync.RWMutex
	Items      []MarketItem
	Categories map[string]*CategoryStat
	nextItemID int
}

func NewMarketplace(filePath string) (*Marketplace, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("ürün dosyası açılamadı: %v", err)
	}
	defer file.Close()

	var products []Product
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&products); err != nil {
		return nil, fmt.Errorf("JSON parse edilemedi: %v", err)
	}

	m := &Marketplace{
		Items:      make([]MarketItem, 0),
		Categories: make(map[string]*CategoryStat),
		nextItemID: 1,
	}

	categoryTotals := make(map[string]float64)
	categoryCounts := make(map[string]float64)

	for _, p := range products {
		cat := p.Category
		if cat == "" {
			cat = "Genel"
		}

		// JSON'da fiyat olmadığı için UrunId'ye göre dinamik ama tutarlı bir fiyat üretiyoruz (50 TL - 450 TL arası)
		// Düzeltilmiş satır (p.Uunn -> p.Urun)
		generatedPrice := 50.0 + (float64(p.ID%400) + float64(len(p.Urun)%50))

		categoryTotals[cat] += generatedPrice
		categoryCounts[cat]++

		m.Items = append(m.Items, MarketItem{
			ID:       m.nextItemID,
			SellerID: 0,
			Urun:     p.Urun,
			Category: cat,
			Price:    generatedPrice,
		})
		m.nextItemID++
	}

	for cat, total := range categoryTotals {
		avgPrice := total / categoryCounts[cat]
		m.Categories[cat] = &CategoryStat{
			CurrentPrice: avgPrice,
			BasePrice:    avgPrice,
			DemandCount:  0,
			SupplyCount:  int(categoryCounts[cat]),
		}
	}

	return m, nil
}

func (m *Marketplace) GetPrice(category string) float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if stat, exists := m.Categories[category]; exists {
		return stat.CurrentPrice
	}
	return 100.0
}

func (m *Marketplace) RecordDemand(category string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if stat, exists := m.Categories[category]; exists {
		stat.DemandCount++
	}
}

func (m *Marketplace) AddItem(sellerID int, urun string, category string, price float64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	item := MarketItem{
		ID:       m.nextItemID,
		SellerID: sellerID,
		Urun:     urun,
		Category: category,
		Price:    price,
	}
	m.nextItemID++
	m.Items = append(m.Items, item)

	if stat, exists := m.Categories[category]; exists {
		stat.SupplyCount++
	}
}

func (m *Marketplace) UpdateDynamicPrices() {
	m.mu.Lock()
	defer m.mu.Unlock()

	elasticity := 0.08

	for _, stat := range m.Categories {
		totalActivity := float64(stat.DemandCount + stat.SupplyCount)
		if totalActivity == 0 {
			continue
		}

		imbalance := float64(stat.DemandCount - stat.SupplyCount) / totalActivity
		priceDelta := stat.CurrentPrice * elasticity * imbalance

		stat.CurrentPrice += priceDelta

		minPrice := stat.BasePrice * 0.2
		if stat.CurrentPrice < minPrice {
			stat.CurrentPrice = minPrice
		}

		stat.DemandCount = 0
		stat.SupplyCount = 0
	}
}

func (m *Marketplace) PrintMarketSummary() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	fmt.Println("\n --- PİYASA FİYAT & ENFLASYON RAPORU ---")
	for cat, stat := range m.Categories {
		var inflation float64
		if stat.BasePrice > 0 {
			inflation = ((stat.CurrentPrice - stat.BasePrice) / stat.BasePrice) * 100.0
		}
		
		arrow := "🟢"
		if inflation > 0 {
			arrow = "🔴" // Enflasyon
		} else if inflation < 0 {
			arrow = "🔵" // Deflasyon
		}
		fmt.Printf("%s %-12s: %.2f TL (Başlangıç: %.2f TL | Değişim: %+.1f%%)\n",
			arrow, cat, stat.CurrentPrice, stat.BasePrice, inflation)
	}
	fmt.Println("----------------------------------------")
}