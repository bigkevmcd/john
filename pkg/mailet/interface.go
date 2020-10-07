package mailet

type Mailet interface {
	Handle(*Mail) error
}
