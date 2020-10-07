package handler

import (
	"bytes"
	"log"
	"net"
	"net/mail"
)

// MailestHandler is a "mhale/smtpd" Handler implementation that processes the
// the email through the configured Mailets.
func MailetsHandler(origin net.Addr, from string, to []string, data []byte) {
	msg, _ := mail.ReadMessage(bytes.NewReader(data))
	subject := msg.Header.Get("Subject")
	log.Printf("Received mail from %s for %s with subject %s", from, to[0], subject)
}
