package devtoolkit

// Pair is a generic pair of values
type Pair[F any, S any] struct {
	First  F
	Second S
}

func NewPair[F any, S any](first F, second S) Pair[F, S] {
	return Pair[F, S]{First: first, Second: second}
}

func (p *Pair[F, S]) GetFirst() F {
	return p.First
}

func (p *Pair[F, S]) GetSecond() S {
	return p.Second
}

func (p *Pair[F, S]) GetAll() (F, S) {
	return p.First, p.Second
}

// Triple is a generic triple of values
type Triple[F any, S any, T any] struct {
	First  F
	Second S
	Third  T
}

func NewTriple[F any, S any, T any](first F, second S, third T) Triple[F, S, T] {
	return Triple[F, S, T]{First: first, Second: second, Third: third}
}

func (t Triple[F, S, T]) GetFirst() F {
	return t.First
}

func (t Triple[F, S, T]) GetSecond() S {
	return t.Second
}

func (t Triple[F, S, T]) GetThird() T {
	return t.Third
}

func (t Triple[F, S, T]) GetAll() (F, S, T) {
	return t.First, t.Second, t.Third
}
