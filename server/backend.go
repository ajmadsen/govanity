package server

type Backend interface {
	IsRepo(name string) (bool, error)
	Canonical(name string) string
}
