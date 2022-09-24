package storage

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type Storage struct {
	storage app.BrowserStorage
}

func From(s app.BrowserStorage) *Storage {
	return &Storage{
		storage: s,
	}
}

type Example = Pair[string, string]

func (s *Storage) GetExamples() ([]Example, error) {
	var list []Example
	if err := s.storage.Get("examples", &list); err !=nil {
		return nil, err
	}
	return list, nil
}

func (s *Storage) AddExample(ex Example) error {
	examples, _ := s.GetExamples()
	examples = append(examples, ex)
	return s.storage.Set("examples", examples)
}

func (s *Storage) DelExample(index int) error {
	examples, _ := s.GetExamples()
	examples = remove(examples, index)
	return s.storage.Set("examples", examples)
}

type Pair[T any, U any] struct {
	First  T
	Second U
}

func remove[T any](slice []T, s int) []T {
	if s >= len(slice) {
		return slice
	}
	return append(slice[:s], slice[s+1:]...)
}
