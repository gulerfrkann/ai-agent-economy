package engine

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Simulation, tüm dünyayı yöneten ana motordur.
type Simulation struct {
	tickInterval time.Duration
	tickCount    uint64
	agents       []Agent // Ajanları burada tutacağız
	mu           sync.RWMutex // Eşzamanlı okuma/yazma için kilit
}

// Yeni bir simülasyon başlatmak için yapıcı (constructor) fonksiyon
func NewSimulation(interval time.Duration) *Simulation {
	return &Simulation{
		tickInterval: interval,
		tickCount:    0,
		agents:       make([]Agent, 0),
	}
}

// Start, zamanı akıtan ana döngüdür (Event Loop)
func (s *Simulation) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(s.tickInterval)
	defer ticker.Stop()

	fmt.Println(" Simülasyon motoru başlatıldı...")

	for {
		select {
		case <-ctx.Done(): // Sistemden kapanma sinyali (CTRL+C) gelirse
			fmt.Printf(" Simülasyon durduruluyor... Toplam Tick: %d\n", s.tickCount)
			return
		case <-ticker.C:
			s.mu.Lock()
			s.tickCount++
			currentTick := s.tickCount
			s.mu.Unlock()

			fmt.Printf("[TICK %d] Zaman aktı. Ajanlar tetikleniyor...\n", currentTick)
			
			// Burada her tick'te ajanların "Act" (harekete geç) metodunu asenkron çağıracağız
			s.triggerAgents(currentTick)
		}
	}
}

// Ajan eklemek için metod
func (s *Simulation) AddAgent(agent Agent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.agents = append(s.agents, agent)
}

// Ajanları tetikleyen fonksiyon
func (s *Simulation) triggerAgents(tick uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, agent := range s.agents {
		// Her ajanı eşzamanlı (goroutine ile) tetikle ki birbirlerini bekletmesinler
		go func(ag Agent) {
			if err := ag.Act(tick); err != nil {
				fmt.Printf("❌ Ajan hatası [%s]: %v\n", ag.GetID(), err)
			}
		}(agent)
	}
}