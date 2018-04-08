package framedredblacktree

import (
	"errors"
	"runtime"

	"github.com/cheekybits/genny/generic"
)

type GKey generic.Type
type GValue generic.Type

type FRBTKeyLessGKeyXXGValue func(p1, p2 GKey) bool

type FRBTreeGKeyXXGValue struct {
	root        *frbtNodeGKeyXXGValue
	generation  uint64
	diversified bool
	lessthan    FRBTKeyLessGKeyXXGValue
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
	ErrKeyNotFoundGKeyXXGValue           = errors.New("Key not found")
)

func NewFRBTreeGKeyXXGValue(lessfunc FRBTKeyLessGKeyXXGValue) *FRBTreeGKeyXXGValue {
	return &FRBTreeGKeyXXGValue{lessthan: lessfunc}
}

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
		if t.lessthan(key, *parent.key) {
			parent.left = inserting
		} else {
			parent.right = inserting
		}
		anchor.hierarchy.Push(parent)
		anchor.at = inserting
		t.insertFixAscendD(anchor)

	}

	return nil
}

func (t *FRBTreeGKeyXXGValue) Drop(key GKey) error {
	if !t.IsModifyAllowed() {
		return ErrModDiversifiedFRBTreeGKeyXXGValue
	}
	anchor := t.narrowto(key)
	if anchor.at == nil {
		return ErrKeyNotFoundGKeyXXGValue
	}
	t.guaranteeAncestorsWriteAccess(anchor)

	if anchor.at.left == nil {
		replacedOrigColor := effectiveColor(anchor.at)
		replacing := anchor.at.right
		t.replacetreeelement(anchor, replacing)
		if replacedOrigColor == BLACK {
			hf := anchor.hierarchy.Fork()
			anchorP := frbtAnchorGKeyXXGValue{at: hf.Pop(), hierarchy: hf}
			t.deleteFixAscendD(anchor, anchorP)
		}

	} else if anchor.at.right == nil {
		replacedOrigColor := effectiveColor(anchor.at)
		replacing := anchor.at.left
		t.replacetreeelement(anchor, replacing)
		if replacedOrigColor == BLACK {
			hf := anchor.hierarchy.Fork()
			anchorP := frbtAnchorGKeyXXGValue{at: hf.Pop(), hierarchy: hf}
			t.deleteFixAscendD(anchor, anchorP)
		}
	} else {
		hifork := anchor.hierarchy.Fork()
		hifork.Push(anchor.at)
		anchorfork := &frbtAnchorGKeyXXGValue{at: anchor.at.right, hierarchy: hifork}
		min := treemin(anchorfork)
		replacedOrigColor := effectiveColor(min.at)
		replacing := min.at.right

		hf := anchor.hierarchy.Fork()
		anchorP := frbtAnchorGKeyXXGValue{at: hf.Pop(), hierarchy: hf}

		if min.hierarchy.Peek() == anchor.at && min.at.right == nil {
			anchorP = *min
		}
		t.guaranteeAncestorsWriteAccess(*min)
		anchor.at.key = min.at.key
		anchor.at.value = min.at.value
		t.replacetreeelement(*min, replacing)
		if replacedOrigColor == BLACK {
			t.deleteFixAscendD(anchor, anchorP)
		}
	}

	t.size--
	return nil
}

func treemin(anchor *frbtAnchorGKeyXXGValue) *frbtAnchorGKeyXXGValue {

	for anchor.at.left != nil {
		anchor.hierarchy.Push(anchor.at)
		anchor.at = anchor.at.left
	}
	return anchor
}

func effectiveColor(v *frbtNodeGKeyXXGValue) int {
	if v == nil {
		return BLACK
	}
	return v.color
}

func (t *FRBTreeGKeyXXGValue) insertFixAscendD(anchor frbtAnchorGKeyXXGValue) {
	//https://www.geeksforgeeks.org/red-black-tree-set-2-insert/ step 3
	t.guaranteeAncestorsWriteAccess(anchor)

	for anchor.hierarchy.Len() != 0 && anchor.hierarchy.Peek().color == RED {
		//small parent or big parent?
		parent := anchor.hierarchy.Pop()
		grandparent := anchor.hierarchy.Pop()
		var uncle *frbtNodeGKeyXXGValue
		reduncle := func(uncle *frbtNodeGKeyXXGValue) {
			parent.color = BLACK
			uncle.color = BLACK
			grandparent.color = RED
			anchor.at = grandparent
		}
		if parent == grandparent.left {
			if grandparent.right != nil {
				uncle = t.guaranteeWriteAccess(grandparent.right)
				grandparent.right = uncle
				if uncle.color == RED {
					reduncle(uncle)
					continue
				}
				if anchor.at == parent.right {
					anchor.hierarchy.Push(grandparent)
					anchor.at = parent
					t.leftRotateM(anchor)
					parent = anchor.hierarchy.Pop()
					grandparent = anchor.hierarchy.Pop()
				}
				parent.color = BLACK
				grandparent.color = RED
				anchor.at = grandparent
				t.rightRotateM(anchor)
				break
			}

		} else if parent == grandparent.right {
			if grandparent.left != nil {
				uncle = t.guaranteeWriteAccess(grandparent.left)
				grandparent.left = uncle
				if uncle.color == RED {
					reduncle(uncle)
					continue
				}
				if anchor.at == parent.left {
					anchor.hierarchy.Push(grandparent)
					anchor.at = parent
					t.rightRotateM(anchor)
					parent = anchor.hierarchy.Pop()
					grandparent = anchor.hierarchy.Pop()
				}
				parent.color = BLACK
				grandparent.color = RED
				anchor.at = grandparent
				t.leftRotateM(anchor)
				break
			}

		} else {
			runtime.Breakpoint()
		}
	}

	t.root.color = BLACK
}

func (t *FRBTreeGKeyXXGValue) deleteFixAscendD(anchor frbtAnchorGKeyXXGValue, replacingParent frbtAnchorGKeyXXGValue) {
	//Arg 1 replacing, arg2 replaced parent
	t.guaranteeAncestorsWriteAccess(anchor)

	for anchor.at != t.root && effectiveColor(anchor.at) == BLACK {
		if anchor.at != nil {
			replacingParent.at = anchor.hierarchy.Peek()
		}
		if anchor.at == replacingParent.at.left {
			sibling := replacingParent.at.right
			sibling = t.guaranteeWriteAccess(sibling)
			replacingParent.at.right = sibling
			if sibling.color == RED {
				sibling.color = BLACK
				replacingParent.at.color = RED
				t.leftRotateM(replacingParent)
				sibling = replacingParent.at.right
			}

			if effectiveColor(sibling.left) == BLACK && effectiveColor(sibling.right) == BLACK {
				sibling.color = RED
				anchor = replacingParent
			} else {
				if effectiveColor(sibling.right) == BLACK {
					if sibling.left != nil {
						sibling.left.color = BLACK
					}
					sibling.color = RED
					replacingParent.hierarchy.Push(replacingParent.at)
					replacingParent.at = sibling
					t.rightRotateM(replacingParent)
					sibling = replacingParent.at.right
				}
				sibling.color = replacingParent.at.color
				replacingParent.at.color = BLACK
				if sibling.right != nil {
					sibling.right.color = BLACK
				}
				t.leftRotateM(replacingParent)
				break
			}
		} else {
			sibling := replacingParent.at.left
			sibling = t.guaranteeWriteAccess(sibling)
			replacingParent.at.left = sibling
			if sibling.color == RED {
				sibling.color = BLACK
				replacingParent.at.color = RED
				t.rightRotateM(replacingParent)
				sibling = replacingParent.at.left
			}

			if effectiveColor(sibling.left) == BLACK && effectiveColor(sibling.right) == BLACK {
				sibling.color = RED
				anchor = replacingParent
			} else {
				if effectiveColor(sibling.left) == BLACK {
					if sibling.right != nil {
						sibling.right.color = BLACK
					}
					sibling.color = RED
					replacingParent.hierarchy.Push(replacingParent.at)
					replacingParent.at = sibling
					t.leftRotateM(replacingParent)
					sibling = replacingParent.at.left
				}
				sibling.color = replacingParent.at.color
				replacingParent.at.color = BLACK
				if sibling.left != nil {
					sibling.left.color = BLACK
				}
				t.rightRotateM(replacingParent)
				break
			}
		}
	}

}

func (t *FRBTreeGKeyXXGValue) leftRotateM(anchor frbtAnchorGKeyXXGValue) frbtAnchorGKeyXXGValue {
	t.guaranteeAncestorsWriteAccess(anchor)
	pGrave := t.guaranteeWriteAccess(anchor.at)
	qGrave := t.guaranteeWriteAccess(anchor.at.right)

	b := anchor.at.right.left

	pGrave.right = b
	qGrave.left = pGrave

	anchor.hierarchy.Push(qGrave)

	return anchor

}

func (t *FRBTreeGKeyXXGValue) rightRotateM(anchor frbtAnchorGKeyXXGValue) frbtAnchorGKeyXXGValue {
	t.guaranteeAncestorsWriteAccess(anchor)
	qGrave := t.guaranteeWriteAccess(anchor.at)
	pGrave := t.guaranteeWriteAccess(anchor.at.left)

	b := anchor.at.left.right

	pGrave.right = qGrave
	qGrave.left = b

	anchor.hierarchy.Push(pGrave)
	anchor.at = qGrave

	return anchor
}

func (t *FRBTreeGKeyXXGValue) replacetreeelement(u frbtAnchorGKeyXXGValue, v *frbtNodeGKeyXXGValue) {
	t.guaranteeAncestorsWriteAccess(u)
	if u.hierarchy.Len() == 0 {
		t.root = v
	} else if u.hierarchy.Peek().left == u.at {
		up := u.hierarchy.Pop()
		up = t.guaranteeWriteAccess(up)
		up.left = v
		u.hierarchy.Push(up)
	} else {
		up := u.hierarchy.Pop()
		up = t.guaranteeWriteAccess(up)
		up.right = v
		u.hierarchy.Push(up)
	}
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

func (t *FRBTreeGKeyXXGValue) guaranteeAncestorsWriteAccess(a frbtAnchorGKeyXXGValue) {
	updatestack := stackGKeyXXGValue{}
	current := a.at

checkfor:
	for {
		if a.hierarchy.Len() != 0 {
			t.root = current
			break checkfor
		}
		old := a.hierarchy.Pop()
		if !t.isShifted(old) {
			break checkfor
		}
		new := t.dupNode(old)
		if old.left == current {
			new.left = current
		} else if old.right == current {
			new.right = current
		} else {
			runtime.Breakpoint()
		}
		updatestack.Push(new)
	}

	for updatestack.Len() != 0 {
		a.hierarchy.Push(updatestack.Pop())
	}

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
		if t.lessthan(key, *current.key) {
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
