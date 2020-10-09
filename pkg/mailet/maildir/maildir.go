package maildir

import (
	"bytes"
	"fmt"
	"io"
	"net/mail"
	"sort"

	"github.com/emersion/go-maildir"

	"github.com/bigkevmcd/john/pkg/mailet"
)

// TODO: support for identifying a path from the to email.

// MaildirMailet is a handler that stores the received mails in a Maildir
// directory.
type MaildirMailet struct {
	path string
}

// New creates and returns a new MaildirMailet that writes mails to
// the provided directory.
func New(dir string) *MaildirMailet {
	return &MaildirMailet{path: dir}
}

func (mm *MaildirMailet) Handle(m *mailet.Mail) error {
	d, err := maildir.NewDelivery(mm.path)
	if err != nil {
		return fmt.Errorf("failed to create a new delivery in %q: %w", mm.path, err)
	}
	var b bytes.Buffer

	for _, k := range headerKeys(m.Message.Header) {
		for _, mv := range m.Message.Header[k] {
			fmt.Fprintf(&b, "%s: %s\n", k, mv)
		}
	}
	fmt.Fprintln(&b) // Headers are separated from the body by a '\n'.

	if _, err := io.Copy(&b, m.Message.Body); err != nil {
		return fmt.Errorf("failed to write the message body to the buffer: %w", err)
	}

	if _, err := d.Write(b.Bytes()); err != nil {
		return fmt.Errorf("failed to write the message body to the delivery: %w", err)
	}
	return d.Close()
}

func headerKeys(h mail.Header) []string {
	keys := []string{}
	for k := range h {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
