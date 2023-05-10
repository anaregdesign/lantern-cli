package query

import (
	"errors"
	"golang.org/x/exp/constraints"
	"strconv"
	"time"
)

var (
	ErrInvalidVerb      = errors.New("invalid verb")
	ErrInvalidObjective = errors.New("invalid objective")

	ErrInvalidArg = errors.New("invalid arg")
)

type Number interface {
	constraints.Integer | constraints.Float
}

func IntArg(s string) (int, error) {
	return strconv.Atoi(s)
}

func StringArg(s string) (string, error) {
	return s, nil
}

func FloatArg(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func BoolArg(s string) (bool, error) {
	return strconv.ParseBool(s)
}

func DatetimeArg(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

func ValueArg(s string) (interface{}, error) {
	if v, err := IntArg(s); err == nil {
		return v, nil
	}
	if v, err := FloatArg(s); err == nil {
		return v, nil
	}
	if v, err := BoolArg(s); err == nil {
		return v, nil
	}
	if v, err := DatetimeArg(s); err == nil {
		return v, nil
	}
	return StringArg(s)
}

func GetVerb(s string) (string, error) {
	if s == "get" {
		return s, nil
	}
	return "", ErrInvalidVerb
}

func PutVerb(s string) (string, error) {
	if s == "put" {
		return s, nil
	}
	return "", ErrInvalidVerb
}

func DeleteVerb(s string) (string, error) {
	if s == "delete" {
		return s, nil
	}
	return "", ErrInvalidVerb
}

func AddVerb(s string) (string, error) {
	if s == "add" {
		return s, nil
	}
	return "", ErrInvalidVerb
}

func VertexObjective(s string) (string, error) {
	if s == "vertex" {
		return s, nil
	}
	return "", ErrInvalidObjective
}

func EdgeObjective(s string) (string, error) {
	if s == "edge" {
		return s, nil
	}
	return "", ErrInvalidObjective
}

func NeighborObjective(s string) (string, error) {
	if s == "neighbor" {
		return s, nil
	}
	return "", ErrInvalidObjective
}

func SPTRelevanceObjective(s string) (string, error) {
	if s == "spt_relevance" {
		return s, nil
	}
	return "", ErrInvalidObjective
}

func SPTCostObjective(s string) (string, error) {
	if s == "spt_cost" {
		return s, nil
	}
	return "", ErrInvalidObjective
}

func MSTRelevanceObjective(s string) (string, error) {
	if s == "mst_relevance" {
		return s, nil
	}
	return "", ErrInvalidObjective
}

func MSTCostObjective(s string) (string, error) {
	if s == "mst_cost" {
		return s, nil
	}
	return "", ErrInvalidObjective
}

func GraphObjective(s string) (string, error) {
	if v, err := NeighborObjective(s); err == nil {
		return v, nil
	}
	if v, err := SPTRelevanceObjective(s); err == nil {
		return v, nil
	}
	if v, err := SPTCostObjective(s); err == nil {
		return v, nil
	}
	if v, err := MSTRelevanceObjective(s); err == nil {
		return v, nil
	}
	if v, err := MSTCostObjective(s); err == nil {
		return v, nil
	}
	return "", ErrInvalidObjective
}
