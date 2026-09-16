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
	fmt.Println(" Otonom Ajan Simülasyonu Başlatılıyor...")

	// 1. İşletim sisteminden gelen kapanma sinyallerini yönetmek için Context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Sinyalleri dinleyen kanal (CTRL+C gibi)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 2. Goroutine'lerin işini bitirmesini beklemek için WaitGroup
	var wg sync.WaitGroup

	// 3. Motoru başlat (Her 1 saniyede bir tetiklenecek şekilde)
	sim := engine.NewSimulation(1 * time.Second)
	// Simülasyonu oluştur
	// Simülasyona otonom ajanlar ekleyelim
	sim.AddAgent(&engine.BaseAgent{ID: "Trader-Alpha", Balance: 1000.0})
	sim.AddAgent(&engine.BaseAgent{ID: "Trader-Beta", Balance: 1500.0})
	sim.AddAgent(&engine.BaseAgent{ID: "Merchant-Gamma", Balance: 500.0})

	wg.Add(1)
	go sim.Start(ctx, &wg)

	// 4. Kullanıcı CTRL+C basana kadar ana programı burada beklet
	<-sigChan
	fmt.Println("\n[SİSTEM] Kapatma sinyali alındı, motor durduruluyor...")
	
	// Motorun içindeki ctx.Done() kanalını tetikler
	cancel()

	// Motorun güvenli bir şekilde durmasını bekle
	wg.Wait()
	fmt.Println(" Simülasyon güvenli bir şekilde kapatıldı.")
}