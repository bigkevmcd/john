package vacations

import (
	"net"
	"net/mail"
	"strings"
	"testing"

	"github.com/bigkevmcd/john/pkg/mailet"
)

var _ mailet.Mailet = (*VacationMailet)(nil)

func TestHandle(t *testing.T) {
	m, err := New()
	assertNoError(t, err)
	data := "From: test@example.com\nSubject: Testing\nTo: user@example.com\n\nTesting email\n"
	msg, err := mail.ReadMessage(strings.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	mail := &mailet.Mail{
		RemoteAddr: &net.TCPAddr{
			IP:   net.ParseIP("192.168.0.252"),
			Port: 32001,
		},
		From:    "test@example.com",
		To:      []string{"user@example.com"},
		Message: *msg,
	}

	assertNoError(t, m.Handle(mail))
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
