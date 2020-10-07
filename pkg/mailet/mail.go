package mailet

import "net"

// Mail represents incoming emails received by the SMTP server.
type Mail struct {
	RemoteAddr net.Addr
	From       string
	To         []string
	Data       []byte
}
