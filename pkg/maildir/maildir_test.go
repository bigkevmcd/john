package maildir

import (
	"testing"

	"github.com/bigkevmcd/john/pkg/mailet"
)

var _ mailet.Mailet = (*MaildirMailet)(nil)

func TestHelloWorld(t *testing.T) {
	// t.Fatal("not implemented")
}
