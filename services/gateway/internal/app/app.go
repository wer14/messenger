package app

import (
	"fmt"

	"github.com/wer14/messenger/services/gateway/internal/app/strategy"
)

type App struct {
	strategies []strategy.Strategy
}

func NewApp(strategies ...strategy.Strategy) *App {
	return &App{
		strategies: strategies,
	}
}

func (s *App) Run() error {
	for _, strategy := range s.strategies {
		if err := strategy.Start(); err != nil {
			return fmt.Errorf("strategy start failed: %w", err)
		}
	}

	return nil
}

func (s *App) Stop() error {
	for _, strategy := range s.strategies {
		if err := strategy.Stop(); err != nil {
			return fmt.Errorf("strategy stop failed: %w", err)
		}
	}

	return nil
}
