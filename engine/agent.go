package engine

import (
	"fmt"
	"math/rand"
)

type Agent interface {
	Act(m *Marketplace, tickNo int, db *DBManager)
	GetID() string
}

type BaseAgent struct {
	ID        string
	Balance   float64
	Inventory []MarketItem
	Interests []string
}

func (a *BaseAgent) GetID() string {
	return a.ID
}

func (a *BaseAgent) Act(m *Marketplace, tickNo int, db *DBManager) {
	m.mu.RLock()
	if len(m.Items) == 0 {
		m.mu.RUnlock()
		return
	}

	// Rastgele bir ilan seç
	randomIndex := rand.Intn(len(m.Items))
	targetItem := m.Items[randomIndex]
	m.mu.RUnlock()

	// Kategori fiyatını dinamik olarak al
	currentPrice := m.GetPrice(targetItem.Category)

	// Eğer ajanın bakiye gücü yetiyorsa satın al
	if a.Balance >= currentPrice {
		m.RecordDemand(targetItem.Category)

		a.Balance -= currentPrice
		purchasedItem := MarketItem{
			ID:       targetItem.ID,
			SellerID: targetItem.SellerID,
			Urun:     targetItem.Urun,
			Category: targetItem.Category,
			Price:    currentPrice,
		}
		a.Inventory = append(a.Inventory, purchasedItem)

		fmt.Printf("🛒 [%s] satın aldı: %s (%s) - Fiyat: %.2f TL (Kalan: %.2f TL)\n",
			a.ID, targetItem.Urun, targetItem.Category, currentPrice, a.Balance)

		// İşlemi SQLite veritabanına asenkron olarak kaydet
		if db != nil {
			db.LogTransaction(AgentTransactionLog{
				TickNo:     tickNo,
				AgentName:  a.ID,
				ItemName:   targetItem.Urun,
				Category:   targetItem.Category,
				Price:      currentPrice,
				NewBalance: a.Balance,
			})
		}
	} else {
		// Paran yetmiyorsa pazara kendi uzmanlık alanından daha ucuz bir ürün koy (Arz yarat)
		if len(a.Interests) > 0 {
			favCategory := a.Interests[rand.Intn(len(a.Interests))]
			basePrice := m.GetPrice(favCategory)
			sellPrice := basePrice * 0.95 // Rekabetçi fiyat
			productName := fmt.Sprintf("%s-Urunu-%d", favCategory, rand.Intn(100))

			m.AddItem(0, productName, favCategory, sellPrice)
		}
	}
}