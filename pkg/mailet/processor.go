package mailet

import (
	"errors"
	"fmt"
)

// TerminateProcessing indicates that this mailet wants to terminate the
// procesing of the chain.
var TerminateProcessing = errors.New("Terminate Processing")

type MailetFunc func(m *Mail)

// Processor applies a mail to each of the set of configured Mailets.
type Processor struct {
	mailets []Mailet
}

// Handle iterates through all the configured Mailets and asks them to handle an
// incoming Mail.
//
// A Mailet can indicate that the processing should terminate by returning the
// TerminateProcessing error.
func (p Processor) Handle(m *Mail) error {
	for _, v := range p.mailets {
		if err := v.Handle(m); err != nil {
			if err == TerminateProcessing {
				return nil
			}
			return fmt.Errorf("failed to process mail: %w", err)
		}
	}
	return nil
}

// NewProcessor creates and returns a new Processor.
func NewProcessor(m []Mailet) *Processor {
	return &Processor{mailets: m}
}
