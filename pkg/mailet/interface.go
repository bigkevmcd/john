package mailet

//go:generate mockgen -destination=./mock_mailet.go -package=mailet github.com/bigkevmcd/john/pkg/mailet Mailet
type Mailet interface {
	Handle(*Mail) error
}
