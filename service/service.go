package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/anaregdesign/lantern-cli/parser"
	"github.com/anaregdesign/lantern/client"
)

var (
	ErrInvalidObjective = fmt.Errorf("invalid objective")
	ErrInvalidVerb      = fmt.Errorf("invalid verb")
	ErrNotImplemented   = fmt.Errorf("not implemented")
)

type CLIService struct {
	client *client.Lantern
}

func NewCLIService(client *client.Lantern) *CLIService {
	return &CLIService{
		client: client,
	}
}

func (c *CLIService) Run(ctx context.Context, str string) error {
	s := parser.NewSource(str)
	verb, err := parser.Verb(s)
	if err != nil {
		return ErrInvalidVerb
	}
	switch verb {
	case "get":
		obj, err := parser.Objective(s)
		if err != nil {
			return ErrInvalidObjective
		}
		switch obj {
		case "vertex":
			p, err := parser.GetVertexParam(s)
			if err != nil {
				return err
			}
			v, err := c.client.GetVertex(ctx, p.Key)
			if err != nil {
				return err
			}
			if jsonString, err := json.Marshal(v.Value); err != nil {
				return err
			} else {
				fmt.Println(string(jsonString))
				return nil
			}
		case "edge":
			p, err := parser.GetEdgeParam(s)
			if err != nil {
				return err
			}
			weight, err := c.client.GetEdge(ctx, p.Tail, p.Head)
			if err != nil {
				return err
			}
			fmt.Printf("%f\n", weight)
			return nil

		default:
			return ErrInvalidObjective
		}
	case "add":
		obj, err := parser.Objective(s)
		if err != nil {
			return ErrInvalidObjective
		}
		switch obj {
		case "edge":
			p, err := parser.AddEdgeParam(s)
			if err != nil {
				return err
			}
			if err := c.client.AddEdge(ctx, p.Tail, p.Head, p.Weight, p.TTL); err != nil {
				return err
			}
			return nil
		default:
			return ErrInvalidObjective
		}
	case "put":
		obj, err := parser.Objective(s)
		if err != nil {
			return ErrInvalidObjective
		}
		switch obj {
		case "vertex":
			p, err := parser.PutVertexParam(s)
			if err != nil {
				return err
			}
			if err := c.client.PutVertex(ctx, p.Key, p.Value, p.TTL); err != nil {
				return err
			}
			return nil
		case "edge":
			return ErrNotImplemented
		}

	case "delete":
		obj, err := parser.Objective(s)
		if err != nil {
			return ErrInvalidObjective
		}
		switch obj {
		case "vertex":
			return ErrNotImplemented
		case "edge":
			return ErrNotImplemented
		default:
			return ErrInvalidObjective
		}

	case "illuminate":
		obj, err := parser.IlluminateObjective(s)
		if err != nil {
			return err
		}
		p, err := parser.IlluminateParam(s)
		if err != nil {
			return err
		}
		g, err := c.client.Illuminate(ctx, p.Seed, p.Step, p.K, p.Tfidf)
		if err != nil {
			return err
		}

		switch obj {
		case "neighbor":
			if jsonString, err := json.MarshalIndent(g, "", "\t"); err != nil {
				return err
			} else {
				fmt.Println(string(jsonString))
				return nil
			}
		case "spt_cost":
			g = g.ShortestPathTree(p.Seed, func(x float32) float32 { return x })
			if jsonString, err := json.MarshalIndent(g, "", "\t"); err != nil {
				return err
			} else {
				fmt.Println(string(jsonString))
				return nil
			}

		case "spt_relevance":
			g = g.ShortestPathTree(p.Seed, func(x float32) float32 { return 1 / x })
			if jsonString, err := json.MarshalIndent(g, "", "\t"); err != nil {
				return err
			} else {
				fmt.Println(string(jsonString))
				return nil
			}
		case "mst_cost":
			g = g.MinimumSpanningTree(p.Seed, false)
			if jsonString, err := json.MarshalIndent(g, "", "\t"); err != nil {
				return err
			} else {
				fmt.Println(string(jsonString))
				return nil
			}
		case "mst_relevance":
			g = g.MinimumSpanningTree(p.Seed, true)
			if jsonString, err := json.MarshalIndent(g, "", "\t"); err != nil {
				return err
			} else {
				fmt.Println(string(jsonString))
				return nil
			}
		}
	default:
		return ErrInvalidVerb
	}
	return nil
}
