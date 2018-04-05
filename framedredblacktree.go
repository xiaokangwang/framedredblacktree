package framedredblacktree

import "github.com/cheekybits/genny/generic"

type GKey generic.Type
type GValue generic.Type

type FRBTKeyLessThan func(p1, p2 *GKey) bool

type FRBTreeGKeyXXGValue struct {
	root        *frbtNodeGKeyXXGValue
	generation  uint64
	diversified bool
	compare     FRBTKeyLessThan
}

type frbtNodeGKeyXXGValue struct {
	left  *frbtNodeGKeyXXGValue
	right *frbtNodeGKeyXXGValue
	shift uint64
	color int
	key   *GKey
	value *GValue
}

const (
	RED   = 0
	BLACK = 1
)

func (t *FRBTreeGKeyXXGValue) Diversify() *FRBTreeGKeyXXGValue {
	t.diversified = true
	return &FRBTreeGKeyXXGValue{
		root:        t.root,
		generation:  t.generation + 1,
		diversified: false,
		compare:     t.compare}
}
