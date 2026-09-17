package engine

import (
	"context" // Bu satırı ekliyoruz
	"fmt"
	"sync"
	"time"
)

type Simulation struct {
	tickInterval time.Duration
	agents       []Agent
	marketplace  *Marketplace // Marketplace doğrudan motorun içine eklendi
	mu           sync.Mutex
}

// NewSimulation artık marketplace'i de parametre olarak alıyor
func NewSimulation(interval time.Duration, mp *Marketplace) *Simulation {
	return &Simulation{
		tickInterval: interval,
		marketplace:  mp,
	}
}

func (s *Simulation) AddAgent(agent Agent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.agents = append(s.agents, agent)
}

func (s *Simulation) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(s.tickInterval)
	defer ticker.Stop()

	var currentTick uint64 = 0
	fmt.Println(" Simülasyon motoru başlatıldı...")

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("\n⏳ Simülasyon durduruluyor... Toplam Tick: %d\n", currentTick)
			return
		case <-ticker.C:
			currentTick++
			fmt.Printf("[TICK %d] Zaman aktı. Ajanlar tetikleniyor...\n", currentTick)
			s.triggerAgents(currentTick) // Artık dışarıdan parametre vermeye gerek yok
		}
	}
}

// triggerAgents motorun içindeki s.marketplace'i doğrudan kullanır
func (s *Simulation) triggerAgents(tick uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, agent := range s.agents {
		go func(ag Agent) {
			if err := ag.Act(tick, s.marketplace); err != nil {
				fmt.Printf(" Ajan hatası [%s]: %v\n", ag.GetID(), err)
			}
		}(agent)
	}
}