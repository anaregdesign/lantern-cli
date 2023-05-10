package parser

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anaregdesign/lantern/client"
	"strconv"
	"strings"
	"time"
)

var ErrInvalidQuery = errors.New("invalid parser")

var ErrGetVertex = errors.New("error, Usage: get vertex <key>")

var ErrGetEdge = errors.New("error, Usage: get edge <from> <to>")
var ErrPutVertex = errors.New("error, Usage: put vertex <key> <value> <ttl>")
var ErrAddEdge = errors.New("error, Usage: add edge <from> <to> <ttl>")
var ErrGetGraph = errors.New("error, Usage: get graph <neighbor|spt_cost|spt_relevance> <seed> <step> <k> <Tfidf>")

type Executor struct {
	Client *client.Lantern
}

func (e *Executor) Execute(ctx context.Context, query string) (string, error) {
	s := strings.ToLower(query)
	params := strings.Split(s, " ")

	if len(params) < 2 {
		return "", ErrInvalidQuery
	}
	switch params[0] {
	case "get":
		switch params[1] {
		case "vertex":
			if len(params) < 3 {
				return "", ErrGetVertex
			}

			key := params[2]
			v, err := e.Client.GetVertex(ctx, key)
			if err != nil {
				return "", err
			}
			if jsonString, err := json.Marshal(v.Value); err != nil {
				return "", err
			} else {
				return string(jsonString), nil
			}

		case "edge":
			if len(params) < 4 {
				return "", ErrGetEdge
			}

			tail := params[2]
			head := params[3]
			if e, err := e.Client.GetEdge(ctx, tail, head); err != nil {
				return "", err
			} else {
				return fmt.Sprintf("%f", e), nil
			}

		case "graph":
			if len(params) < 7 {
				return "", ErrGetGraph
			}

			seed := params[3]
			step, err := strconv.Atoi(params[4])
			if err != nil {
				return "", ErrGetGraph
			}
			k, err := strconv.Atoi(params[5])
			if err != nil {
				return "", ErrGetGraph
			}

			tfidf, err := strconv.ParseBool(params[6])
			if err != nil {
				return "", ErrGetGraph
			}

			g, err := e.Client.Illuminate(ctx, seed, step, k, tfidf)
			if err != nil {
				return "", err
			}

			switch params[2] {
			case "neighbor":
				if jsonString, err := json.MarshalIndent(g, "", "\t"); err != nil {
					return "", err
				} else {
					return string(jsonString), nil
				}
			case "spt_cost":
				g := g.ShortestPathTree(seed, func(x float32) float32 { return x })
				if jsonString, err := json.MarshalIndent(g, "", "\t"); err != nil {
					return "", err
				} else {
					return string(jsonString), nil
				}
			case "spt_relevance":
				g := g.ShortestPathTree(seed, func(x float32) float32 { return 1 / x })
				if jsonString, err := json.MarshalIndent(g, "", "\t"); err != nil {
					return "", err
				} else {
					return string(jsonString), nil
				}
			case "msp_cost":
				g := g.MinimumSpanningTree(seed, false)
				if jsonString, err := json.MarshalIndent(g, "", "\t"); err != nil {
					return "", err
				} else {
					return string(jsonString), nil
				}
			case "msp_relevance":
				g := g.MinimumSpanningTree(seed, true)
				if jsonString, err := json.MarshalIndent(g, "", "\t"); err != nil {
					return "", err
				} else {
					return string(jsonString), nil
				}
			}
		default:
			return "", ErrGetGraph
		}

	case "put":
		switch params[1] {
		case "vertex":
			if len(params) < 5 {
				return "", ErrPutVertex
			}

			key := params[2]

			var value interface{}
			if v, err := strconv.ParseFloat(params[3], 32); err == nil {
				value = v
			}
			if v, err := strconv.ParseBool(params[3]); err == nil {
				value = v
			}
			if v, err := strconv.Atoi(params[3]); err == nil {
				value = v
			}
			if value == nil {
				value = params[3]
			}

			ttl, err := strconv.Atoi(params[4])
			if err != nil {
				return "", ErrPutVertex
			}
			if err := e.Client.PutVertex(ctx, key, value, time.Duration(ttl)*time.Second); err != nil {
				return "", err
			}
		default:
			return "", ErrPutVertex
		}

	case "add":
		switch params[1] {
		case "edge":
			if len(params) < 6 {
				return "", ErrAddEdge
			}

			tail := params[2]
			head := params[3]
			value, err := strconv.ParseFloat(params[4], 32)
			if err != nil {
				return "", ErrAddEdge
			}
			ttl, err := strconv.Atoi(params[5])
			if err != nil {
				return "", ErrAddEdge
			}

			if err := e.Client.AddEdge(ctx, tail, head, float32(value), time.Duration(ttl)*time.Second); err != nil {
				return "", err
			}
		default:
			return "", ErrInvalidQuery
		}
	default:
		return "", ErrInvalidQuery

	}
	return "", nil
}
