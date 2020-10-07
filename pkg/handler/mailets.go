package handler

import (
	"bytes"
	"log"
	"net"
	"net/mail"

	"github.com/bigkevmcd/john/pkg/mailet"
	"github.com/mhale/smtpd"
)

// MakeHandler creates a new handler configured appropriately.
func MakeHandler(handlers ...mailet.Mailet) smtpd.Handler {
	p := mailet.NewProcessor(handlers)
	return func(origin net.Addr, from string, to []string, data []byte) {
		msg, err := mail.ReadMessage(bytes.NewReader(data))
		if err != nil {
			log.Printf("failed to ReadMessage: %s", err)
			return
		}
		m := &mailet.Mail{
			RemoteAddr: origin,
			From:       from,
			To:         to,
			Message:    *msg,
		}
		if err := p.Handle(m); err != nil {
			log.Printf("failed to Handle message: %s", err)
		}
	}
}
