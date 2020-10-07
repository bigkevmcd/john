package handler

import (
	"io"
	"io/ioutil"
	"net"
	"net/mail"
	"strings"
	"testing"

	"github.com/bigkevmcd/john/pkg/mailet"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

const data = "From: test@example.com\nTo: test1@example.com\nSubject: testing\n\nTesting"

var testAddress = &net.TCPAddr{
	IP:   net.ParseIP("192.168.0.252"),
	Port: 32001,
}

func TestHandler(t *testing.T) {
	m := &stubMailet{handled: []*mailet.Mail{}}
	h := MakeHandler(m)
	e := makeTestMail()

	h(e.RemoteAddr, e.From, e.To, []byte(data))

	want := []*mailet.Mail{e}
	// Can't compare the bodies directly.
	if diff := cmp.Diff(want, m.handled, cmpopts.IgnoreFields(mail.Message{}, "Body")); diff != "" {
		t.Fatalf("handled mails:\n%s", diff)
	}
	if diff := cmp.Diff(mustReadAll(t, e.Message.Body), mustReadAll(t, m.handled[0].Message.Body)); diff != "" {
		t.Fatalf("handled mail bodies:\n%s", diff)
	}
}

type stubMailet struct {
	handled []*mailet.Mail
}

func (s *stubMailet) Handle(m *mailet.Mail) error {
	s.handled = append(s.handled, m)
	return nil
}

func mustReadAll(t *testing.T, i io.Reader) string {
	t.Helper()
	b, err := ioutil.ReadAll(i)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func makeTestMail() *mailet.Mail {
	return &mailet.Mail{
		RemoteAddr: testAddress,
		From:       "test@example.com",
		To:         []string{"test1@example.com"},
		Message: mail.Message{
			Header: mail.Header{
				"From":    {"test@example.com"},
				"Subject": {"testing"},
				"To":      {"test1@example.com"},
			},
			Body: strings.NewReader("Testing"),
		},
	}
}
