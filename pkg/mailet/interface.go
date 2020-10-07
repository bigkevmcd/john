package mailet

// Mailet is a handler for incoming SMTP mails.
type Mailet interface {
	Handle(*Mail) error
}
