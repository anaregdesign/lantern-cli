package parser

import (
	"errors"
	"strconv"
	"time"
)

var (
	Verbs = []string{
		"get",
		"put",
		"delete",
		"illuminate",
	}

	Objectives = []string{
		"vertex",
		"edge",
	}
	IlluminateObjectives = []string{
		"neighbor",
		"spt_relevance",
		"spt_cost",
		"mst_relevance",
		"mst_cost",
	}

	ErrInvalidVerb      = errors.New("invalid verb")
	ErrInvalidObjective = errors.New("invalid objective")

	ErrOutOfChoice = errors.New("out of choice")
	ErrInvalidArg  = errors.New("invalid arg")
)

func String(s *Source) (string, error) {
	defer s.Next()
	str, err := s.Peek()
	if err != nil {
		return "", err
	}

	return str, nil
}

func Integer(s *Source) (int, error) {
	defer s.Next()
	str, err := s.Peek()
	if err != nil {
		return 0, err
	}
	v, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func Float(s *Source) (float64, error) {
	defer s.Next()
	str, err := s.Peek()
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func Bool(s *Source) (bool, error) {
	defer s.Next()
	str, err := s.Peek()
	if err != nil {
		return false, err
	}
	v, err := strconv.ParseBool(str)
	if err != nil {
		return false, err
	}
	return v, nil
}

func Datetime(s *Source) (time.Time, error) {
	defer s.Next()
	str, err := s.Peek()
	if err != nil {
		return time.Time{}, err
	}
	v, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return time.Time{}, err
	}
	return v, nil
}

func Value(s *Source) (interface{}, error) {
	defer s.Next()
	str, err := s.Peek()
	if err != nil {
		return struct{}{}, err
	}
	if v, err := strconv.Atoi(str); err == nil {
		return v, nil
	}
	if v, err := strconv.ParseFloat(str, 64); err == nil {
		return v, nil
	}
	if v, err := strconv.ParseBool(str); err == nil {
		return v, nil
	}
	if v, err := time.Parse(time.RFC3339, str); err == nil {
		return v, nil
	}
	return str, nil
}

func AnyOf(s *Source, choices []string) (string, error) {
	defer s.Next()
	str, err := s.Peek()
	if err != nil {
		return "", err
	}
	for _, v := range choices {
		if str == v {
			return str, nil
		}
	}
	return "", ErrOutOfChoice
}

func Verb(s *Source) (string, error) {
	return AnyOf(s, Verbs)
}

func Objective(s *Source) (string, error) {
	return AnyOf(s, Objectives)
}

func IlluminateObjective(s *Source) (string, error) {
	return AnyOf(s, IlluminateObjectives)
}
