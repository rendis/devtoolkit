package devtoolkit

import "errors"

type (
	LinkFn[T any]   func(T) error
	SaveStep[T any] func(T) error
)

var (
	ErrNilLinkFn = errors.New("nil link function")
)

// ProcessChain defines an interface for a chain of operations (links) that can be executed
// on data of type T. It allows adding links, setting a save step, executing the chain,
// and retrieving the sequence of added links.
type ProcessChain[T any] interface {
	// AddLink adds a new link, identified by a string key, to the chain of operations.
	// It returns an error if the provided link function is nil.
	AddLink(string, LinkFn[T]) error

	// SetSaveStep sets a save step function that is executed after each link in the chain.
	// This step is used to persist the state of the data after each operation.
	SetSaveStep(SaveStep[T])

	// Execute runs the process chain on data of type T, sequentially executing
	// each link in the order they were added.
	// It returns a slice of string keys representing the successfully executed links and an error if the execution
	// of any link fails.
	Execute(T) ([]string, error)

	// GetChain returns a slice of string keys representing the sequence of links added to the chain.
	GetChain() []string
}

// NewProcessChain creates and returns a new instance of a process chain for data of type T.
func NewProcessChain[T any]() ProcessChain[T] {
	return &processChain[T]{
		links: make(map[string]LinkFn[T]),
	}
}

type processChain[T any] struct {
	links    map[string]LinkFn[T]
	linksSeq []string
	saveStep SaveStep[T]
}

func (p *processChain[T]) AddLink(s string, l LinkFn[T]) error {
	if l == nil {
		return ErrNilLinkFn
	}
	p.links[s] = l
	p.linksSeq = append(p.linksSeq, s)
	return nil
}

func (p *processChain[T]) SetSaveStep(s SaveStep[T]) {
	p.saveStep = s
}

func (p *processChain[T]) Execute(t T) ([]string, error) {
	var successExecutedLinks []string

	for _, link := range p.linksSeq {
		if err := p.links[link](t); err != nil {
			return successExecutedLinks, err
		}

		if p.saveStep != nil {
			if err := p.saveStep(t); err != nil {
				return successExecutedLinks, err
			}
		}

		successExecutedLinks = append(successExecutedLinks, link)
	}

	return successExecutedLinks, nil
}

func (p *processChain[T]) GetChain() []string {
	var cp = make([]string, len(p.linksSeq))
	copy(cp, p.linksSeq)
	return cp
}