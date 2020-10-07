package maildir

import (
	"bytes"
	"fmt"

	"github.com/emersion/go-maildir"

	"github.com/bigkevmcd/john/pkg/mailet"
)

// TODO: support for identifying a path from the to email.

// MaildirMailet is a handler that stores the received mails in a Maildir
// directory.
type MaildirMailet struct {
	path string
}

// NewMaildirMailet creates and returns a new MaildirMailet that writes mails to
// the provided directory.
func NewMaildirMailet(dir string) *MaildirMailet {
	return &MaildirMailet{path: dir}
}

func (mm *MaildirMailet) Handle(m *mailet.Mail) error {
	d, err := maildir.NewDelivery(mm.path)
	if err != nil {
		return fmt.Errorf("failed to create a new delivery in %q: %w", mm.path, err)
	}
	var b bytes.Buffer

	if _, err := b.Write(m.Data); err != nil {
		return fmt.Errorf("failed to write the message body to the buffer: %w", err)
	}

	if _, err := d.Write(b.Bytes()); err != nil {
		return fmt.Errorf("failed to write the message body to the delivery: %w", err)
	}
	return d.Close()
}
