package maildir

import (
	"io/ioutil"
	"net"
	"net/mail"
	"os"
	"strings"
	"testing"

	"github.com/bigkevmcd/john/pkg/mailet"
	"github.com/emersion/go-maildir"
)

var _ mailet.Mailet = (*MaildirMailet)(nil)

func TestHandle(t *testing.T) {
	base := tempDir(t)
	mm, err := New(base)
	assertNoError(t, err)
	data := "From: test@example.com\nSubject: Testing\nTo: user@example.com\n\nTesting email\n"
	m, err := mail.ReadMessage(strings.NewReader(data))
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
		Message: *m,
	}

	assertNoError(t, mm.Handle(mail))

	dir := maildir.Dir(base)
	unseen, err := dir.Unseen()
	assertNoError(t, err)
	if c := len(unseen); c != 1 {
		t.Fatalf("expected %d mails to be received, got %d", 1, c)
	}

	o, err := dir.Open(unseen[0])
	assertNoError(t, err)
	b, err := ioutil.ReadAll(o)
	assertNoError(t, err)
	if c := string(b); c != data {
		t.Fatalf("failed to write body, got %q, want %q", c, data)
	}
}

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := ioutil.TempDir(os.TempDir(), "john")
	assertNoError(t, err)

	t.Cleanup(func() {
		err := os.RemoveAll(dir)
		assertNoError(t, err)
	})
	return dir
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
