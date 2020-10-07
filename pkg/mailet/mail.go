package mailet

import (
	"net"
	"net/mail"
)

// Mail represents incoming emails received by the SMTP server.
type Mail struct {
	RemoteAddr net.Addr
	From       string
	To         []string
	Message    mail.Message
}
