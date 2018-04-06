package framedredblacktree

import (
	"errors"

	"github.com/cheekybits/genny/generic"
)

type GKey generic.Type
type GValue generic.Type

type FRBTKeyLessThan func(p1, p2 GKey) bool

type FRBTreeGKeyXXGValue struct {
	root        *frbtNodeGKeyXXGValue
	generation  uint64
	diversified bool
	lessthan    FRBTKeyLessThan
	size        int
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

var (
	ErrModDiversifiedFRBTreeGKeyXXGValue = errors.New("Modify a diversified FRBTreeGKeyXXGValue")
)

func (t *FRBTreeGKeyXXGValue) Diversify() *FRBTreeGKeyXXGValue {
	t.diversified = true
	return &FRBTreeGKeyXXGValue{
		root:        t.root,
		generation:  t.generation + 1,
		diversified: false,
		lessthan:    t.lessthan}
}

func (t *FRBTreeGKeyXXGValue) IsModifyAllowed() bool {
	return !t.diversified
}

func (t *FRBTreeGKeyXXGValue) Insert(key GKey, value GValue) error {
	if !t.IsModifyAllowed() {
		return ErrModDiversifiedFRBTreeGKeyXXGValue
	}

	_, hierarchy := t.narrowto(key)

	t.size++

	if hierarchy.Len() == 0 {
		//Inserting first node
		t.root = t.makeNode(BLACK, key, value)
	} else {
		inserting := t.makeNode(RED, key, value)
		parent := t.guaranteeWriteAccess(hierarchy.Pop())
		if t.lessthan(key, parent.key) {
			parent.left = inserting
		} else {
			parent.right = inserting
		}
		t.insertFixAscend(inserting, hierarchy)

	}

	return nil
}

func (t *FRBTreeGKeyXXGValue) insertFixAscend(current *frbtNodeGKeyXXGValue, hierarchy *stackGKeyXXGValue) {

}

func (t *FRBTreeGKeyXXGValue) makeNode(color int, key GKey, value GValue) *frbtNodeGKeyXXGValue {
	return &frbtNodeGKeyXXGValue{color: color, shift: t.generation, key: &key, value: &value}
}

func (t *FRBTreeGKeyXXGValue) dupNode(src *frbtNodeGKeyXXGValue) *frbtNodeGKeyXXGValue {
	return &frbtNodeGKeyXXGValue{
		color: src.color,
		shift: t.generation,
		key:   src.key,
		value: src.value,
		left:  src.left,
		right: src.right,
	}
}

func (t *FRBTreeGKeyXXGValue) isShifted(src *frbtNodeGKeyXXGValue) bool {
	return src.shift != t.generation
}

func (t *FRBTreeGKeyXXGValue) guaranteeWriteAccess(src *frbtNodeGKeyXXGValue) *frbtNodeGKeyXXGValue {
	if t.isShifted(src) {
		return t.dupNode(src)
	}
	return src
}

/*narrow down to the nearest node
  return true if exact match is found,
              with an stack topped with result and parents hierarchy
         false if no exact match is found,
              with an stack topped with would be parents hierarchy
*/
func (t *FRBTreeGKeyXXGValue) narrowto(key GKey) (bool, *stackGKeyXXGValue) {
	hierarchystack := newstackGKeyXXGValue()
	current := t.root
	for {
		if current == nil {
			return false, hierarchystack
		}
		hierarchystack.Push(current)
		if t.lessthan(key, current.key) {
			current = current.left
		} else {
			if current.key == key {
				hierarchystack.Push(current)
				return true, hierarchystack
			}
			current = current.right
		}
	}

}

type (
	stackGKeyXXGValue struct {
		top    *stacknodefrbtNodeGKeyXXGValue
		length int
	}
	stacknodefrbtNodeGKeyXXGValue struct {
		value *frbtNodeGKeyXXGValue
		prev  *stacknodefrbtNodeGKeyXXGValue
	}
)

// Create a new stack
func newstackGKeyXXGValue() *stackGKeyXXGValue {
	return &stackGKeyXXGValue{nil, 0}
}

// Return the number of items in the stack
func (this *stackGKeyXXGValue) Len() int {
	return this.length
}

// View the top item on the stack
func (this *stackGKeyXXGValue) Peek() *frbtNodeGKeyXXGValue {
	if this.length == 0 {
		return nil
	}
	return this.top.value
}

// Pop the top item of the stack and return it
func (this *stackGKeyXXGValue) Pop() *frbtNodeGKeyXXGValue {
	if this.length == 0 {
		return nil
	}

	n := this.top
	this.top = n.prev
	this.length--
	return n.value
}

// Push a value onto the top of the stack
func (this *stackGKeyXXGValue) Push(value *frbtNodeGKeyXXGValue) {
	n := &stacknodefrbtNodeGKeyXXGValue{value, this.top}
	this.top = n
	this.length++
}
