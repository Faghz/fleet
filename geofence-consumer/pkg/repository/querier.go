package repository

//go:generate mockgen -destination=../mocks/querier-mock.go -package=mocks github.com/elzestia/fleet/pkg/repository Querier
type Querier interface {
}

var _ Querier = (*Queries)(nil)
