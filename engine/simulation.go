package engine

import (
	"context"
	"fmt"
	"sync"
	"time"
)



type Simulation struct {
	Interval    time.Duration
	Marketplace *Marketplace
	Agents      []Agent
	DB          *DBManager
	TickCount   int
}

func NewSimulation(interval time.Duration, m *Marketplace, dbMgr *DBManager) *Simulation {
	return &Simulation{
		Interval:    interval,
		Marketplace: m,
		Agents:      make([]Agent, 0),
		DB:          dbMgr,
		TickCount:   0,
	}
}

func (s *Simulation) AddAgent(agent Agent) {
	s.Agents = append(s.Agents, agent)
}

func (s *Simulation) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(s.Interval)
	defer ticker.Stop()

	fmt.Printf("[SİSTEM] Simülasyon motoru başlatıldı. Tick aralığı: %v\n", s.Interval)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("[SİSTEM] Simülasyon döngüsü durduruluyor...")
			return
		case <-ticker.C:
			s.TickCount++
			fmt.Printf("\n--- TICK #%d Başladı ---\n", s.TickCount)

			// 1. Ajanların eylemleri (Alış/Satış)
			var agentWg sync.WaitGroup
			for _, agent := range s.Agents {
				agentWg.Add(1)
				go func(a Agent) {
					defer agentWg.Done()
					a.Act(s.Marketplace, s.TickCount, s.DB)
				}(agent)
			}
			agentWg.Wait()

			// 2. Piyasa fiyatlarını ve enflasyonu güncelle (Tick logları DB'ye burada yazılır)
			s.Marketplace.UpdateDynamicPrices(s.TickCount)

			fmt.Printf("--- TICK #%d Tamamlandı ---\n", s.TickCount)
		}
	}
}