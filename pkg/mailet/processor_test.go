package mailet

import (
	"errors"
	"net"
	"net/mail"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestProcessMailets(t *testing.T) {
	mailets := []Mailet{
		&stubMailet{name: "stub1"},
		&stubMailet{name: "stub2"},
	}

	p := NewProcessor(mailets)
	m := makeTestMail()

	assertNoError(t, p.Handle(m))

	want := []string{"stub1", "stub2"}
	if diff := cmp.Diff(want, m.Message.Header["stubs"]); diff != "" {
		t.Fatalf("processing stubs failed:\n%s", diff)
	}
}

func TestProcessMailetsWithAnError(t *testing.T) {
	mailets := []Mailet{
		&stubMailet{name: "stub1", err: errors.New("this is a test")},
		&stubMailet{name: "stub2"},
	}

	p := NewProcessor(mailets)
	m := makeTestMail()

	if err := p.Handle(m); err == nil {
		t.Fatal("failed to get an error handling mail")
	}

	want := []string{"stub1"}
	if diff := cmp.Diff(want, m.Message.Header["stubs"]); diff != "" {
		t.Fatalf("processing stubs failed:\n%s", diff)
	}
}

func TestProcessMailetsWithATerminationError(t *testing.T) {
	mailets := []Mailet{
		&stubMailet{name: "stub1", err: TerminateProcessing},
		&stubMailet{name: "stub2"},
	}

	p := NewProcessor(mailets)
	m := makeTestMail()
	assertNoError(t, p.Handle(m))

	want := []string{"stub1"}
	if diff := cmp.Diff(want, m.Message.Header["stubs"]); diff != "" {
		t.Fatalf("processing stubs failed:\n%s", diff)
	}
}

type stubMailet struct {
	name string
	err  error
}

func (s *stubMailet) Handle(m *Mail) error {
	m.Message.Header["stubs"] = append(m.Message.Header["stubs"], s.name)
	return s.err
}

func makeTestMail() *Mail {
	return &Mail{
		RemoteAddr: &net.TCPAddr{
			IP:   net.ParseIP("192.168.0.252"),
			Port: 32001,
		},
		From: "test@example.com",
		To:   []string{"test1@example.com"},
		Message: mail.Message{
			Header: mail.Header{},
			Body:   strings.NewReader("From: test@example.com\nTo: test1@example.com\nSubject: testing\n\nTesting"),
		},
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
