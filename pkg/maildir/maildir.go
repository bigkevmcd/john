package maildir

import (
	"github.com/bigkevmcd/john/pkg/mailet"
)

// MaildirMailet is a handler that stores the received mails in a Maildir
// directory.
type MaildirMailet struct {
}

func (mm *MaildirMailet) Handle(m mailet.Mail) error {
	return nil
}
