package server

import (
	"errors"
	"log"
	"strings"
)

type TextBackend map[string]string

func (b TextBackend) IsRepo(name string) (bool, error) {
	_, ok := b[name]
	return ok, nil
}

func (b TextBackend) Canonical(name string) string {
	return b[name]
}

func (b TextBackend) Add(redir string) error {
	parts := strings.SplitN(redir, ":", 2)
	if len(parts) != 2 {
		return errors.New("malformed; must be of form \"repo-name:canonical-path\"")
	}
	if _, ok := b[parts[0]]; ok {
		return errors.New("conflict; redirection already exists")
	}
	b[parts[0]] = parts[1]
	log.Printf("handling redirection from %q -> %q", parts[0], parts[1])
	return nil
}
