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
	Interests []string // Yeni: Ajanın ilgi duyduğu/uzmanlaştığı kategoriler (Örn: "Kitap", "Kırtasiye")
}

func (a *BaseAgent) GetID() string {
	return a.ID
}

func (a *BaseAgent) Act(tick uint64, marketplace *Marketplace) error {
	// 1. Gelir Döngüsü: Her ajan her tick'te sistemden taban gelir elde eder (150 TL)
	earnedIncome := 150.00
	a.Balance += earnedIncome

	if len(marketplace.Products) == 0 {
		return nil
	}

	// 2. Akıllı Ürün Seçimi: Uzmanlık alanına uygun ürün bulmaya çalışalım
	var selectedProduct Product
	found := false

	// Eğer ajanın ilgi alanları varsa, önce o kategorilerden ürün arayalım (10 deneme hakkı verelim)
	if len(a.Interests) > 0 {
		for i := 0; i < 10; i++ {
			p := marketplace.Products[rand.Intn(len(marketplace.Products))]
			// Ürünün kategorisi ajanın ilgi alanlarından biriyle eşleşiyor mu?
			for _, interest := range a.Interests {
				if p.KtgAdi == interest {
					selectedProduct = p
					found = true
					break
				}
			}
			if found {
				break
			}
		}
	}

	// Eğer ilgi alanına uygun bulunamadıysa veya ajanın ilgi alanı yoksa rastgele bir ürün seç
	if !found {
		selectedProduct = marketplace.Products[rand.Intn(len(marketplace.Products))]
	}

	// 3. Satın Alma Mantığı
	if a.Balance >= selectedProduct.Price {
		a.Balance -= selectedProduct.Price
		a.Inventory = append(a.Inventory, selectedProduct)

		fmt.Printf("💰 [Ajan: %s] SATIN ALDI! [%s] | Ürün: '%s' | Fiyat: %.2f TL | Kalan Bakiye: %.2f TL\n",
			a.ID, selectedProduct.KtgAdi, selectedProduct.Urun, selectedProduct.Price, a.Balance)
	} else {
		fmt.Printf("🤖 [Ajan: %s] (Bakiye: %.2f TL) - Tick %d: '%s' (%s) inceledi, bakiye yetersiz (Fiyat: %.2f TL).\n",
			a.ID, a.Balance, tick, selectedProduct.Urun, selectedProduct.KtgAdi, selectedProduct.Price)
	}

	return nil
}