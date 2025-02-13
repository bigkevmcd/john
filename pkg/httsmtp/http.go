package httsmtp

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"github.com/bigkevmcd/john/pkg/mailet"
)

const (
	EnvelopeFromHeader   = "John-Envelope-From"
	EnvelopeToHeader     = "John-Envelope-To"
	BodyMailHeaderPrefix = "John-Mail-"
)

// SMTPHandler is an http.Handler that converts HTTP requests to SMTP mail
// bodies and passes them off to a Handler.
type SMTPHandler struct {
	mailet.Mailet
}

// MakeHandler creates and returns an HTTP handler that parses HTTP requests and
// converts them to mails before processing them with the Handler.
func MakeHandler(h mailet.Mailet) *SMTPHandler {
	return &SMTPHandler{
		Mailet: h,
	}
}

// ServeHTTP implements the http.Handler interface.
func (s *SMTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mail, err := mailFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.Mailet.Handle(mail); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func mailFromRequest(r *http.Request) (*mailet.Mail, error) {
	from := r.Header.Get(EnvelopeFromHeader)
	if from == "" {
		return nil, fmt.Errorf("header %s must be provided", EnvelopeFromHeader)
	}

	to := r.Header.Values(EnvelopeToHeader)
	if len(to) == 0 {
		return nil, fmt.Errorf("header %s must be provided", EnvelopeToHeader)
	}

	rawHost, rawPort, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(rawPort)
	if err != nil {
		return nil, fmt.Errorf("parsing RemoteAddr port: %w", err)
	}

	trimPrefix := func(s []string) []string {
		var trimmed []string

		for _, v := range s {
			trimmed = append(trimmed, strings.TrimPrefix(v, BodyMailHeaderPrefix))
		}

		return trimmed
	}

	headers := mail.Header{}
	for header, value := range r.Header {
		if strings.HasPrefix(header, BodyMailHeaderPrefix) {
			headers[strings.TrimPrefix(header, BodyMailHeaderPrefix)] = trimPrefix(value)
		}
	}

	buf, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("parsing RemoteAddr port: %w", err)
	}

	m := &mailet.Mail{
		From: from,
		To:   to,
		RemoteAddr: &net.TCPAddr{
			IP:   net.ParseIP(rawHost),
			Port: port,
		},
		Message: mail.Message{
			Header: headers,
			Body:   bytes.NewReader(buf),
		},
	}

	return m, nil
}
