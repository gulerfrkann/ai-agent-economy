package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gulerfrkann/ai-agent-economy/engine"
)

func main() {
	fmt.Println("Otonom Ajan Simülasyonu Başlatılıyor...")

	// 1. Pazar yeri verilerini yükle (Python ile dönüştürdüğümüz JSON)
	marketplace, err := engine.NewMarketplace("data/products.json")
	if err != nil {
		fmt.Printf(" Kritik Hata: %v\n", err)
		return
	}

	// 2. İşletim sisteminden gelen kapanma sinyallerini yönetmek için Context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	// 3. Motoru başlat
	sim := engine.NewSimulation(1 * time.Second, marketplace)

	// Ajanlarımızı uzmanlık alanlarıyla birlikte tanımlayalım
	agentAlpha := &engine.BaseAgent{
		ID:        "Trader-Alpha",
		Balance:   1000.0,
		Interests: []string{"Kitap", "Edebiyat", "Tarih"},
	}
	agentBeta := &engine.BaseAgent{
		ID:        "Trader-Beta",
		Balance:   1500.0,
		Interests: []string{"Oyuncak", "Bebek", "Kırtasiye"},
	}
	agentGamma := &engine.BaseAgent{
		ID:        "Merchant-Gamma",
		Balance:   500.0,
		Interests: []string{"Elektronik", "Film", "Müzik"},
	}

	sim.AddAgent(agentAlpha)
	sim.AddAgent(agentBeta)
	sim.AddAgent(agentGamma)

	wg.Add(1)
	go sim.Start(ctx, &wg)

	<-sigChan
	fmt.Println("\n[SİSTEM] Kapatma sinyali alındı, motor durduruluyor...")
	
	cancel()
	wg.Wait()
	fmt.Println(" Simülasyon güvenli bir şekilde kapatıldı.")

	// 4. Simülasyon Raporu: Ajanların son durumunu yazdıralım
	fmt.Println("\n === SİMÜLASYON FİNAL RAPORU ===")
	agents := []engine.Agent{agentAlpha, agentBeta, agentGamma}
	for _, agent := range agents {
		if bAgent, ok := agent.(*engine.BaseAgent); ok {
			fmt.Printf(" Ajan: %-15s | Kalan Bakiye: %8.2f TL | Envanter Ürün Sayısı: %d\n",
				bAgent.ID, bAgent.Balance, len(bAgent.Inventory))
			
			// İsteğe bağlı: Envanterindeki son 2-3 ürünü de detaylı gösterelim ki aldıklarını görebilelim
			if len(bAgent.Inventory) > 0 {
				fmt.Println("    Aldığı Bazı Ürünler:")
				for i, prod := range bAgent.Inventory {
					if i >= 3 { // Çok uzatmamak için ilk 3 ürünü gösterelim
						fmt.Printf("      ... ve %d ürün daha\n", len(bAgent.Inventory)-3)
						break
					}
					fmt.Printf("      - %s (%.2f TL)\n", prod.Urun, prod.Price)
				}
			}
			fmt.Println("--------------------------------------------------")
		}
	}
}