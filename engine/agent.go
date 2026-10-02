package engine

import (
	"errors"
	"fmt"
	"math/rand"
)

type Agent interface {
	GetID() string
	Act(tick uint64, marketplace *Marketplace) error
}

type BaseAgent struct {
	ID        string
	Balance   float64
	Inventory []MarketItem // Artık MarketItem kullanıyoruz
	Interests []string     // Kategori uzmanlıkları
}

func (a *BaseAgent) GetID() string {
	return a.ID
}

func (a *BaseAgent) Act(tick uint64, marketplace *Marketplace) error {
	// 1. Pazar tahtasından rastgele bir ürün seç
	marketplace.mu.RLock()
	if len(marketplace.Items) == 0 {
		marketplace.mu.RUnlock()
		return errors.New("pazarda ürün kalmadı")
	}

	// Rastgele bir ilan seç
	randomIndex := rand.Intn(len(marketplace.Items))
	targetItem := marketplace.Items[randomIndex]
	marketplace.mu.RUnlock()

	// 2. Kategori fiyatını dinamik olarak al
	currentPrice := marketplace.GetPrice(targetItem.Category)

	// Eğer ajanın bakiye gücü yetiyorsa ve ilgilendiği kategorideyse satın almayı dene
	if a.Balance >= currentPrice {
		// Talep kaydını düş (Arz-talep dengesi ve enflasyon için)
		marketplace.RecordDemand(targetItem.Category)

		a.Balance -= currentPrice
		a.Inventory = append(a.Inventory, MarketItem{
			ID:       targetItem.ID,
			SellerID: targetItem.SellerID,
			Urun:     targetItem.Urun, // <-- Buraya ürün adını ekliyoruz
			Category: targetItem.Category,
			Price:    currentPrice,
		})

		fmt.Printf("🛒 [%s] ajan satın aldı: %s kategorisinden ürün (Fiyat: %.2f TL, Kalan Bakiye: %.2f TL)\n",
			a.ID, targetItem.Category, currentPrice, a.Balance)
	} else {
		// Paran yetmiyorsa pazara kendi uzmanlık alanından daha ucuz bir ürün koy (Arz yarat)
		// Paran yetmiyorsa pazara kendi uzmanlık alanından daha ucuz bir ürün koy (Arz yarat)
		if len(a.Interests) > 0 {
			favCategory := a.Interests[rand.Intn(len(a.Interests))]
			basePrice := marketplace.GetPrice(favCategory)
			sellPrice := basePrice * 0.95 // Rekabetçi fiyat
			productName := fmt.Sprintf("%s-Urunu-%d", favCategory, rand.Intn(100))

			marketplace.AddItem(0, productName, favCategory, sellPrice)
		}
	}

	return nil
}