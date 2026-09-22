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
	if len(marketplace.Products) == 0 {
		return nil
	}

	// Pazar yerinden rastgele bir ürün seç
	randomIndex := rand.Intn(len(marketplace.Products))
	selectedProduct := marketplace.Products[randomIndex]

	// Bakiye kontrolü ve satın alma mantığı
	if a.Balance >= selectedProduct.Price {
		a.Balance -= selectedProduct.Price
		a.Inventory = append(a.Inventory, selectedProduct)

		fmt.Printf(" [Ajan: %s] SATIN ALDI! | Ürün: '%s' | Fiyat: %.2f TL | Kalan Bakiye: %.2f TL\n",
			a.ID, selectedProduct.Urun, selectedProduct.Price, a.Balance)
	} else {
		fmt.Printf(" [Ajan: %s] (Bakiye: %.2f TL) - Tick %d: '%s' ürününü inceledi, bakiye yetersiz (Fiyat: %.2f TL).\n",
			a.ID, a.Balance, tick, selectedProduct.Urun, selectedProduct.Price)
	}

	return nil
}