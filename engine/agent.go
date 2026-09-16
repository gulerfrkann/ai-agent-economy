package engine
import (
	"fmt"
)

// Agent, simülasyon içerisindeki her bir yapay zeka aktörünü temsil eder.
type Agent interface {
	GetID() string
	Act(tick uint64) error
}

// BaseAgent, temel özellikleri barındıran basit bir ajan yapısıdır.
type BaseAgent struct {
	ID      string
	Balance float64
}

func (a *BaseAgent) GetID() string {
	return a.ID
}

func (a *BaseAgent) Act(tick uint64) error {
	// Şimdilik her ajan tick başına konsola durumunu yazdırsın
	fmt.Printf(" [Ajan: %s] (Bakiye: %.2f TL) - Tick %d işleniyor...\n", a.ID, a.Balance, tick)
	return nil
}