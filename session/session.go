package session

import "fmt"

type Session struct {
	Repo string
	Base string
	Head string
}

func New(repo, base, head string) (*Session, error) {
	if repo == "" {
		return nil, fmt.Errorf("repo path (-C) is required")
	}
	if base == "" {
		return nil, fmt.Errorf("base ref (--base) is required")
	}
	if head == "" {
		return nil, fmt.Errorf("head ref (--head) is required")
	}
	return &Session{Repo: repo, Base: base, Head: head}, nil
}
