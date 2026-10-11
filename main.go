package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"ai-agent-economy/engine"
)

func main() {
	fmt.Println("Otonom Ajan Simülasyonu Başlatılıyor...")

	// 1. Veritabanı Yöneticisini Başlat (Doğrudan engine paketinden)
	dbMgr, err := engine.NewDBManager("simulation.db")
	if err != nil {
		log.Fatalf(" Kritik Hata - Veritabanı başlatılamadı: %v", err)
	}
	defer dbMgr.Close()

	// 2. Pazar yeri verilerini yükle ve veritabanı referansını ver
	marketplace, err := engine.NewMarketplace("data/products.json", dbMgr)
	if err != nil {
		fmt.Printf(" Kritik Hata: %v", err)
		return
	}

	// 3. İşletim sisteminden gelen kapanma sinyallerini yönetmek için Context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	// 4. Motoru başlat (DB yöneticisi simülasyona da veriliyor)
	sim := engine.NewSimulation(1*time.Second, marketplace, dbMgr)

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

	// 5. Simülasyon Raporu: Ajanların ve Piyasanın son durumu
	fmt.Println("\n === SİMÜLASYON FİNAL RAPORU ===")
	agents := []engine.Agent{agentAlpha, agentBeta, agentGamma}
	for _, agent := range agents {
		if bAgent, ok := agent.(*engine.BaseAgent); ok {
			fmt.Printf(" Ajan: %-15s | Kalan Bakiye: %8.2f TL | Envanter Ürün Sayısı: %d\n",
				bAgent.ID, bAgent.Balance, len(bAgent.Inventory))
			
			if len(bAgent.Inventory) > 0 {
				fmt.Println("    Aldığı Bazı Ürünler:")
				for i, prod := range bAgent.Inventory {
					if i >= 3 {
						fmt.Printf("      ... ve %d ürün daha\n", len(bAgent.Inventory)-3)
						break
					}
					fmt.Printf("      - %s (%.2f TL)\n", prod.Urun, prod.Price)
				}
			}
		}
	}

	// Piyasa Fiyat ve Enflasyon Raporu
	marketplace.PrintMarketSummary()
}