package engine

import (
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
	Inventory []Product
}

func (a *BaseAgent) GetID() string {
	return a.ID
}

func (a *BaseAgent) Act(tick uint64, marketplace *Marketplace) error {
	// 1. Gelir Döngüsü: Her ajan her tick'te sistemden taban gelir elde eder (Örn: 150 TL)
	earnedIncome := 150.00
	a.Balance += earnedIncome

	if len(marketplace.Products) == 0 {
		return nil
	}

	// Pazar yerinden rastgele bir ürün seç
	randomIndex := rand.Intn(len(marketplace.Products))
	selectedProduct := marketplace.Products[randomIndex]

	// 2. Satın Alma Mantığı
	if a.Balance >= selectedProduct.Price {
		a.Balance -= selectedProduct.Price
		a.Inventory = append(a.Inventory, selectedProduct)

		fmt.Printf(" [Ajan: %s] SATIN ALDI! (+%.2f TL Gelir) | Ürün: '%s' | Fiyat: %.2f TL | Kalan Bakiye: %.2f TL\n",
			a.ID, earnedIncome, selectedProduct.Urun, selectedProduct.Price, a.Balance)
	} else {
		fmt.Printf(" [Ajan: %s] (+%.2f TL Gelir | Bakiye: %.2f TL) - Tick %d: '%s' inceledi, bakiye yetersiz (Fiyat: %.2f TL).\n",
			a.ID, earnedIncome, a.Balance, tick, selectedProduct.Urun, selectedProduct.Price)
	}

	return nil
}