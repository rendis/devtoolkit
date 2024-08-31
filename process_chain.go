package devtoolkit

import (
	"context"
	"errors"
	"time"
)

type (
	LinkFn[T any]   func(context.Context, T) error
	SaveStep[T any] func(context.Context, T, []string) error
)

var (
	ErrNilLinkFn = errors.New("nil link function")
)

// ProcessChain defines an interface for a chain of operations (links) that can be executed
// on data of type T. It allows adding links, setting a save Step, executing the chain,
// and retrieving the sequence of added links.
type ProcessChain[T any] interface {
	// AddLink adds a new link, identified by a string key, to the chain of operations.
	// It returns an error if the provided link function is nil.
	AddLink(string, LinkFn[T]) error

	// AddLinkWithWait adds a new link, identified by a string key, to the chain of operations.
	// It returns an error if the provided link function is nil.
	AddLinkWithWait(s string, l LinkFn[T], wait time.Duration) error

	// AddLinks adds multiple links to the chain of operations.
	// It returns an error if any of the provided link functions is nil.
	AddLinks(links []LinkInfo[T]) error

	// SetSaveStep sets a save Step function that is executed after each link in the chain.
	// This Step is used to persist the state of the data after each operation.
	SetSaveStep(SaveStep[T])

	// GetChain returns a slice of string keys representing the sequence of links added to the chain.
	GetChain() []string

	// Execute runs the process chain on data of type T, sequentially executing
	// each link in the order they were added.
	// It returns a slice of string keys representing the successfully executed links and an error if the execution
	// of any link fails.
	Execute(context.Context, T) ([]string, error)

	// ExecuteWithIgnorableLinks runs the process chain on data of type T, sequentially executing
	// each link in the order they were added, except for the ignorable links.
	// It returns a slice of string keys representing the successfully executed links and an error if the execution
	// of any link fails.
	ExecuteWithIgnorableLinks(context.Context, T, []string) ([]string, error)
}

type ProcessChainOptions struct {
	AddLinkNameToError bool // default: false
}

func setProcessChainOptionsDefaults(opts *ProcessChainOptions) *ProcessChainOptions {
	if opts == nil {
		opts = &ProcessChainOptions{
			AddLinkNameToError: false,
		}
	}
	return opts
}

// NewProcessChain creates and returns a new instance of a process chain for data of type T.
func NewProcessChain[T any](opts *ProcessChainOptions) ProcessChain[T] {
	opts = setProcessChainOptionsDefaults(opts)
	return &processChain[T]{
		addLinkNameToError: opts.AddLinkNameToError,
	}
}

type LinkInfo[T any] struct {
	Mame string
	Step LinkFn[T]
	Wait time.Duration
}

type processChain[T any] struct {
	links              []*LinkInfo[T]
	saveStep           SaveStep[T]
	addLinkNameToError bool
}

func (p *processChain[T]) AddLink(s string, l LinkFn[T]) error {
	return p.AddLinkWithWait(s, l, 0)
}

func (p *processChain[T]) AddLinkWithWait(s string, l LinkFn[T], wait time.Duration) error {
	if l == nil {
		return ErrNilLinkFn
	}
	li := &LinkInfo[T]{Mame: s, Step: l, Wait: wait}
	p.links = append(p.links, li)
	return nil
}

func (p *processChain[T]) AddLinks(links []LinkInfo[T]) error {
	for _, link := range links {
		if err := p.AddLinkWithWait(link.Mame, link.Step, link.Wait); err != nil {
			return err
		}
	}
	return nil
}

func (p *processChain[T]) SetSaveStep(s SaveStep[T]) {
	p.saveStep = s
}

func (p *processChain[T]) GetChain() []string {
	var chain []string
	for _, link := range p.links {
		chain = append(chain, link.Mame)
	}
	return chain
}

func (p *processChain[T]) Execute(ctx context.Context, t T) ([]string, error) {
	return p.execute(ctx, t, nil)
}

func (p *processChain[T]) ExecuteWithIgnorableLinks(ctx context.Context, t T, ignorableLinks []string) ([]string, error) {
	var ignorableLinksMap = make(map[string]struct{})

	for _, link := range ignorableLinks {
		ignorableLinksMap[link] = struct{}{}
	}

	return p.execute(ctx, t, ignorableLinksMap)
}

func (p *processChain[T]) execute(ctx context.Context, t T, ignorableLinks map[string]struct{}) ([]string, error) {
	var successExecutedLinks []string

	for _, link := range p.links {
		linkName := link.Mame

		if _, ok := ignorableLinks[linkName]; ok {
			successExecutedLinks = append(successExecutedLinks, linkName)
			continue
		}

		if err := link.Step(ctx, t); err != nil {
			if p.addLinkNameToError {
				err = errors.New(linkName + ": " + err.Error())
			}
			return successExecutedLinks, err
		}

		successExecutedLinks = append(successExecutedLinks, linkName)

		if link.Wait > 0 {
			time.Sleep(link.Wait)
		}

		if p.saveStep != nil {
			if err := p.saveStep(ctx, t, successExecutedLinks); err != nil {
				if p.addLinkNameToError {
					err = errors.New("saveStep: " + err.Error())
				}
				return successExecutedLinks[:len(successExecutedLinks)-1], err
			}
		}
	}

	return successExecutedLinks, nil
}
