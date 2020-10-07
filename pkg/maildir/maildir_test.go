package maildir

import (
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/bigkevmcd/john/pkg/mailet"
	"github.com/emersion/go-maildir"
)

var _ mailet.Mailet = (*MaildirMailet)(nil)

func TestHandle(t *testing.T) {
	base := tempDir(t)
	mm := NewMaildirMailet(base)
	data := "From: test@example.com\nTo: user@example.com\nSubject: Testing\n\nTesting email\n"

	mail := mailet.Mail{
		RemoteAddr: &net.TCPAddr{
			IP:   net.ParseIP("192.168.0.252"),
			Port: 32001,
		},
		From: "test@example.com",
		To:   []string{"user@example.com"},
		Data: []byte(data),
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

	// TODO: Should this move?
	for _, v := range []string{"tmp", "new", "cur"} {
		assertNoError(t, os.MkdirAll(filepath.Join(dir, v), 0700))
	}

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
