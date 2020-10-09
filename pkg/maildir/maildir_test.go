package maildir

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/emersion/go-maildir"
	"github.com/google/go-cmp/cmp"
)

func TestPollMails(t *testing.T) {
	body := "From: test@example.com\nTo: user@example.com\nSubject: Testing\n\nThis is the body"
	dir := tempMaildir(t)
	del, err := maildir.NewDelivery(dir)
	assertNoError(t, err)
	_, err = del.Write([]byte(body))
	assertNoError(t, err)
	assertNoError(t, del.Close())

	md := maildir.Dir(dir)
	unseen, err := md.Unseen()
	assertNoError(t, err)
	if u := len(unseen); u != 1 {
		t.Fatalf("unseen count = %d, want %d", u, 1)
	}

	r, err := md.Open(unseen[0])
	assertNoError(t, err)
	b, err := ioutil.ReadAll(r)
	assertNoError(t, err)

	if diff := cmp.Diff(body, string(b)); diff != "" {
		t.Fatalf("failed:\n%s", diff)
	}

	unseen, err = md.Unseen()
	assertNoError(t, err)
	if u := len(unseen); u != 0 {
		t.Fatalf("unseen count = %d, want %d", u, 0)
	}

}

func tempMaildir(t *testing.T) string {
	t.Helper()
	dir, err := ioutil.TempDir(os.TempDir(), "john")
	assertNoError(t, err)

	err = maildir.Dir(dir).Init()
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
