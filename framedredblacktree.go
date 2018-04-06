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

type frbtAnchorGKeyXXGValue struct {
	at        *frbtNodeGKeyXXGValue
	hierarchy *stackGKeyXXGValue
}

const (
	RED   = 0
	BLACK = 1
)

var (
	ErrModDiversifiedFRBTreeGKeyXXGValue = errors.New("Modify a diversified FRBTreeGKeyXXGValue")
	ErrKeyNotFound                       = errors.New("Key not found")
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

	anchor := t.narrowto(key)

	t.size++

	if anchor.hierarchy.Len() == 0 {
		//Inserting first node
		t.root = t.makeNode(BLACK, key, value)
	} else {
		inserting := t.makeNode(RED, key, value)
		parent := t.guaranteeWriteAccess(anchor.hierarchy.Pop())
		if t.lessthan(key, parent.key) {
			parent.left = inserting
		} else {
			parent.right = inserting
		}
		anchor.hierarchy.Push(parent)
		anchor.at = inserting
		t.insertFixAscend(anchor)

	}

	return nil
}

func (t *FRBTreeGKeyXXGValue) insertFixAscend(anchor frbtAnchorGKeyXXGValue) {
	//https://www.geeksforgeeks.org/red-black-tree-set-2-insert/ step 3
}

func (t *FRBTreeGKeyXXGValue) deleteFixAscend(anchor frbtAnchorGKeyXXGValue) {

}

func (t *FRBTreeGKeyXXGValue) leftRotateM(anchor frbtAnchorGKeyXXGValue) frbtAnchorGKeyXXGValue {
	pGrave := t.guaranteeWriteAccess(anchor.at)
	qGrave := t.guaranteeWriteAccess(anchor.at.right)

	b := anchor.at.right.left

	pGrave.right = b
	qGrave.left = pGrave

	anchor.hierarchy.Push(qGrave)

	return anchor

}

func (t *FRBTreeGKeyXXGValue) rightRotateM(anchor frbtAnchorGKeyXXGValue) frbtAnchorGKeyXXGValue {
	qGrave := t.guaranteeWriteAccess(anchor.at)
	pGrave := t.guaranteeWriteAccess(anchor.at.left)

	b := anchor.at.left.right

	pGrave.right = qGrave
	qGrave.left = b

	anchor.hierarchy.Push(pGrave)
	anchor.at = qGrave

	return anchor
}

func (t *FRBTreeGKeyXXGValue) transplant(u frbtAnchorGKeyXXGValue, v frbtAnchorGKeyXXGValue) {

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
func (t *FRBTreeGKeyXXGValue) narrowto(key GKey) frbtAnchorGKeyXXGValue {
	hierarchystack := newstackGKeyXXGValue()
	current := t.root
	for {
		if current == nil {
			return frbtAnchorGKeyXXGValue{at: nil, hierarchy: hierarchystack}
		}
		hierarchystack.Push(current)
		if t.lessthan(key, current.key) {
			current = current.left
		} else {
			if current.key == key {
				return frbtAnchorGKeyXXGValue{at: current, hierarchy: hierarchystack}
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

func (this *stackGKeyXXGValue) Fork() *stackGKeyXXGValue {
	return &stackGKeyXXGValue{this.top, this.length}
}

func (this *stackGKeyXXGValue) Peek() *frbtNodeGKeyXXGValue {
	if this.length == 0 {
		return nil
	}
	return this.top.value
}

// Return the number of items in the stack
func (this *stackGKeyXXGValue) Len() int {
	return this.length
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
