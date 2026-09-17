package engine

import (
	"fmt"
	"math/rand"
)

// Agent arayüzü
type Agent interface {
	GetID() string
	Act(tick uint64, marketplace *Marketplace) error
}

// BaseAgent, temel özellikleri barındırır
type BaseAgent struct {
	ID        string
	Balance   float64
	Inventory []Product // Ajanın sahip olduğu ürünler
}

func (a *BaseAgent) GetID() string {
	return a.ID
}

// Act metoduna pazar yerini de dahil ediyoruz ki ürün seçebilsinler
func (a *BaseAgent) Act(tick uint64, marketplace *Marketplace) error {
	// Pazar yerinden rastgele bir ürün seçelim
	if len(marketplace.Products) > 0 {
		randomIndex := rand.Intn(len(marketplace.Products))
		selectedProduct := marketplace.Products[randomIndex]

		fmt.Printf("🤖 [Ajan: %s] (Bakiye: %.2f TL) - Tick %d: Pazar yerinden '%s' (%s) ürününü inceledi.\n", 
			a.ID, a.Balance, tick, selectedProduct.Urun, selectedProduct.MarkaAd)
	}
	return nil
}