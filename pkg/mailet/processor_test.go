package mailet

import (
	"errors"
	"fmt"
	"net"
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

	if diff := cmp.Diff("stub1\nstub2\n", string(m.Data)); diff != "" {
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

	if diff := cmp.Diff("stub1\n", string(m.Data)); diff != "" {
		t.Fatalf("processing stubs failed:\n%s", diff)
	}
}

type stubMailet struct {
	name string
	err  error
}

func (s *stubMailet) Handle(m *Mail) error {
	m.Data = append(m.Data, []byte(fmt.Sprintf("%s\n", s.name))...)
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
		Data: []byte(""),
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
