package vacations

import (
	"time"

	"github.com/bigkevmcd/john/pkg/mailet"
)

// Vacation represents a period during which to send vacation notifications to
// the "from" element of a processed email.
type Vacation struct {
	Start time.Time
	End   time.Time
	Email string
}

// VacationMailet is a handler that can trigger vacation responses
// automatically.
type VacationMailet struct {
}

// New creates and returns a new VacationMailet.
func New() (*VacationMailet, error) {
	return &VacationMailet{}, nil
}

func (mm *VacationMailet) Handle(m *mailet.Mail) error {
	return nil
}
