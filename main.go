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
	fmt.Println("🌍 Otonom Ajan Simülasyonu Başlatılıyor...")

	// 1. Pazar yeri verilerini yükle (Python ile dönüştürdüğümüz JSON)
	marketplace, err := engine.NewMarketplace("data/products.json")
	if err != nil {
		fmt.Printf("❌ Kritik Hata: %v\n", err)
		return
	}
	// İleride ajanlar bu marketplace üzerinden rastgele ürünler seçecek/satacak
	_ = marketplace

	// 2. İşletim sisteminden gelen kapanma sinyallerini yönetmek için Context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	// 3. Motoru başlat
	sim := engine.NewSimulation(1 * time.Second)

	sim.AddAgent(&engine.BaseAgent{ID: "Trader-Alpha", Balance: 1000.0})
	sim.AddAgent(&engine.BaseAgent{ID: "Trader-Beta", Balance: 1500.0})
	sim.AddAgent(&engine.BaseAgent{ID: "Merchant-Gamma", Balance: 500.0})

	wg.Add(1)
	go sim.Start(ctx, &wg)

	<-sigChan
	fmt.Println("\n[SİSTEM] Kapatma sinyali alındı, motor durduruluyor...")
	
	cancel()
	wg.Wait()
	fmt.Println("✅ Simülasyon güvenli bir şekilde kapatıldı.")
}