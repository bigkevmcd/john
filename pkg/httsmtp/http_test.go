package httsmtp

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/mail"
	"strings"
	"testing"

	"github.com/bigkevmcd/john/pkg/mailet"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestServeHTTPMissingHeaders(t *testing.T) {
	fake := &fakeMailet{}

	srv := httptest.NewServer(MakeHandler(fake))
	t.Cleanup(srv.Close)

	testMail := &mailet.Mail{
		RemoteAddr: &net.TCPAddr{
			IP:   net.ParseIP("192.168.0.252"),
			Port: 32001,
		},
		Message: mail.Message{
			Body: strings.NewReader("Testing\n\n"),
		},
	}

	req := newHTTPRequestFromMail(t, srv, testMail)

	resp, err := srv.Client().Do(req)
	assertNoError(t, err)

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("got StatusCode %v, want %v", resp.StatusCode, http.StatusBadRequest)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if diff := cmp.Diff("header John-Envelope-From must be provided\n", string(b)); diff != "" {
		t.Errorf("invalid error response: diff -want +got\n%s", diff)
	}
}

func TestServeHTTP(t *testing.T) {
	fake := &fakeMailet{}

	srv := httptest.NewServer(MakeHandler(fake))
	t.Cleanup(srv.Close)

	req := newHTTPRequestFromMail(t, srv, newTestMail())

	_, err := srv.Client().Do(req)
	assertNoError(t, err)

	want := []*mailet.Mail{
		{
			RemoteAddr: &net.TCPAddr{
				IP:   net.ParseIP("192.168.0.252"),
				Port: 32001,
			},
			From: "test@example.com",
			To:   []string{"test1@example.com"},
			Message: mail.Message{
				Header: mail.Header{
					"From":    []string{"test@example.com"},
					"To":      []string{"test1@example.com"},
					"Subject": []string{"testing"},
				},
				Body: strings.NewReader("Testing\n"),
			},
		},
	}
	// Ignore the RemoteAddr because it's from the request client and the Body
	// because it's an io.Reader.
	if diff := cmp.Diff(want, fake.captured, cmpopts.IgnoreFields(mailet.Mail{}, "RemoteAddr"), cmpopts.IgnoreFields(mail.Message{}, "Body")); diff != "" {
		t.Errorf("sending mail failed: diff -want +got\n%s", diff)
	}

	b, err := io.ReadAll(fake.captured[0].Message.Body)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff("Testing\n\n", string(b)); diff != "" {
		t.Errorf("sending mail .Body failed: diff -want +got\n%s", diff)
	}
}

func newHTTPRequestFromMail(t *testing.T, ts *httptest.Server, mail *mailet.Mail) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, ts.URL, mail.Message.Body)
	if err != nil {
		t.Fatal(err)
	}

	if mail.From != "" {
		req.Header.Set(EnvelopeFromHeader, mail.From)
	}

	for _, to := range mail.To {
		req.Header.Add(EnvelopeToHeader, to)
	}

	for header, values := range mail.Message.Header {
		for _, value := range values {
			req.Header.Add(BodyMailHeaderPrefix+header, value)
		}
	}

	return req
}

type fakeMailet struct {
	captured []*mailet.Mail
	err      error
}

func (s *fakeMailet) Handle(m *mailet.Mail) error {
	s.captured = append(s.captured, m)
	return s.err
}

func newTestMail() *mailet.Mail {
	return &mailet.Mail{
		RemoteAddr: &net.TCPAddr{
			IP:   net.ParseIP("192.168.0.252"),
			Port: 32001,
		},
		From: "test@example.com",
		To:   []string{"test1@example.com"},
		Message: mail.Message{
			Header: mail.Header{
				"From":    []string{"test@example.com"},
				"To":      []string{"test1@example.com"},
				"Subject": []string{"testing"},
			},
			Body: strings.NewReader("Testing\n\n"),
		},
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
