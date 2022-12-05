package main

import (
	"errors"
	"golang.org/x/net/html"
)

type stack struct {
	stack []html.Token
	// TODO: they say it is needed. But I do not understand why, so remove it for now
	// lock  sync.RWMutex
}

func (s *stack) Push(el html.Token) {
	s.stack = append(s.stack, el)
}

func (s *stack) Pop() error {
	l := len(s.stack)
	if l > 0 {
		s.stack = s.stack[:l-1]
		return nil
	}

	return errors.New("stack is empty")
}

func (s *stack) Head() (html.Token, error) {
	l := len(s.stack)
	if l > 0 {
		return s.stack[l-1], nil
	}

	return html.Token{}, errors.New("stack is empty")
}

func (s *stack) Empty() bool {
	return len(s.stack) == 0
}
