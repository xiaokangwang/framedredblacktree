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

type (
	stackGKey struct {
		top    *nodeGKey
		length int
	}
	nodeGKey struct {
		value GKey
		prev  *nodeGKey
	}
)

// Create a new stack
func New() *stackGKey {
	return &stackGKey{nil, 0}
}

// Return the number of items in the stack
func (this *stackGKey) Len() int {
	return this.length
}

// View the top item on the stack
func (this *stackGKey) Peek() GKey {
	if this.length == 0 {
		return nil
	}
	return this.top.value
}

// Pop the top item of the stack and return it
func (this *stackGKey) Pop() GKey {
	if this.length == 0 {
		return nil
	}

	n := this.top
	this.top = n.prev
	this.length--
	return n.value
}

// Push a value onto the top of the stack
func (this *stackGKey) Push(value GKey) {
	n := &nodeGKey{value, this.top}
	this.top = n
	this.length++
}
