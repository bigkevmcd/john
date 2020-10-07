package mailet

import "fmt"

type Processor struct {
	mailets []Mailet
}

func (p Processor) Handle(m *Mail) error {
	for _, v := range p.mailets {
		if err := v.Handle(m); err != nil {
			return fmt.Errorf("failed to process mail: %w", err)
		}
	}
	return nil
}

func NewProcessor(m []Mailet) *Processor {
	return &Processor{mailets: m}
}
