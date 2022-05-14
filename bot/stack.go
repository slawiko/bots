package main

import (
	"errors"
)

type stack struct {
	stack []string
	// TODO: they say it is needed. But I do not understand why, so remove it for now
	// lock  sync.RWMutex
}

func (s *stack) Push(el string) {
	s.stack = append(s.stack, el)
}

func (s *stack) Pop() error {
	l := len(s.stack)
	if l > 0 {
		s.stack = s.stack[:l-1]
		return nil
	}

	return errors.New("Stack is empty")
}

func (s *stack) Head() (string, error) {
	l := len(s.stack)
	if l > 0 {
		return s.stack[l-1], nil
	}

	return "", errors.New("Stack is empty")
}
