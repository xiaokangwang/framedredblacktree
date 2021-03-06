package framedredblacktree

import (
	"errors"
	"runtime"

	"github.com/cheekybits/genny/generic"
)

//type GKey generic.Type
type GKey uint
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
	at         *frbtNodeGKeyXXGValue
	hierarchy  *stackGKeyXXGValue
	lastremove int
}

const (
	RED   = 0
	BLACK = 1
	left  = 2
	right = 3
	root  = 4
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

	anchor := t.narrowto(key, true)
	t.guaranteeAncestorsWriteAccess(anchor)
	t.size++

	if anchor.hierarchy.Len() == 0 {
		//Inserting first node
		t.root = t.makeNode(BLACK, key, value)
	} else {
		inserting := t.makeNode(RED, key, value)
		var parent *frbtNodeGKeyXXGValue

		parent = t.guaranteeWriteAccess(anchor.hierarchy.Pop(), nil)

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
	anchor := t.narrowto(key, false)
	if anchor.at == nil {
		return ErrKeyNotFoundGKeyXXGValue
	}
	t.guaranteeAncestorsWriteAccess(anchor)

	if anchor.at.left == nil {
		replacedOrigColor := effectiveColor(anchor.at)
		replacing := anchor.at.right
		t.replacetreeelement(&anchor, replacing)
		if anchor.lastremove == root {
			if t.root != nil {
				t.root.color = BLACK
			}
			return nil
		}
		if replacedOrigColor == BLACK {
			//hf := anchor.hierarchy.Fork()
			//anchorP := frbtAnchorGKeyXXGValue{at: hf.Pop(), hierarchy: hf}
			t.deleteFixAscendD(anchor, anchor)
		}

	} else if anchor.at.right == nil {
		replacedOrigColor := effectiveColor(anchor.at)
		replacing := anchor.at.left
		t.replacetreeelement(&anchor, replacing)
		if replacedOrigColor == BLACK {
			t.deleteFixAscendD(anchor, anchor)
		}
	} else {
		hifork := anchor.hierarchy.Fork()
		hifork.Push(anchor.at)
		anchorfork := &frbtAnchorGKeyXXGValue{at: anchor.at.right, hierarchy: hifork}
		min := treemin(anchorfork)
		replacedOrigColor := effectiveColor(min.at)
		replacing := min.at.right

		//hf := anchor.hierarchy.Fork()
		//anchorP := frbtAnchorGKeyXXGValue{at: hf.Pop(), hierarchy: hf}

		if min.hierarchy.Peek() == anchor.at && min.at.right == nil {
			//anchorP = *min
		}
		t.guaranteeAncestorsWriteAccess(*min)
		anchor.at.key = min.at.key
		anchor.at.value = min.at.value

		t.replacetreeelement(min, replacing)
		anchor.lastremove = min.lastremove
		if replacedOrigColor == BLACK {
			t.deleteFixAscendD(*min, *min)
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
	var round = 1000
	for anchor.hierarchy.Len() != 0 && anchor.hierarchy.Peek().color == RED {
		round--
		if round == 0 {
			runtime.Breakpoint()
		}

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
				uncle = t.guaranteeWriteAccess(grandparent.right, nil)
				grandparent.right = uncle
				if uncle.color == RED {
					reduncle(uncle)
					continue
				}
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

		} else if parent == grandparent.right {
			if grandparent.left != nil {
				uncle = t.guaranteeWriteAccess(grandparent.left, nil)
				grandparent.left = uncle
				if uncle.color == RED {
					reduncle(uncle)
					continue
				}
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

		} else {
			runtime.Breakpoint()
		}
	}

	t.root.color = BLACK
}

func (t *FRBTreeGKeyXXGValue) deleteFixAscendD(anchor frbtAnchorGKeyXXGValue, replacingParent frbtAnchorGKeyXXGValue) {
	//Arg 1 replacing, arg2 replaced parent
	//t.guaranteeAncestorsWriteAccess(anchor)
	var round = 1000
	for anchor.at != t.root && effectiveColor(anchor.at) == BLACK {
		round--
		if round == 0 {
			runtime.Breakpoint()
		}
		if anchor.at != nil {
			replacingParent.at = anchor.hierarchy.Peek()
		}
		if anchor.lastremove == left {
			if effectiveColor(replacingParent.at.left) == RED {
				replacingParent.at.left.color = BLACK
				return
			}
			sibling := replacingParent.at.right
			if sibling == nil {
				//Sorry, I don't know how to fix it gracefully
				return
			}
			sibling = t.guaranteeWriteAccess(sibling, nil)
			replacingParent.at.right = sibling
			if sibling.color == RED {
				sibling.color = BLACK
				replacingParent.at.color = RED
				//Need Debug
				//replacingParent.hierarchy.Push(replacingParent.at)
				replacingParent.hierarchy.Pop()
				replacingParent = t.leftRotateM(replacingParent)
				replacingParent.hierarchy.Push(replacingParent.at)
				sibling = replacingParent.at.right
			}
			if sibling == nil {
				//Sorry, I don't know how to fix it gracefully
				return
			}

			if effectiveColor(sibling.left) == BLACK && effectiveColor(sibling.right) == BLACK {
				sibling.color = RED
				replacingParent.at = replacingParent.hierarchy.Pop()
				anchor = replacingParent
				if replacingParent.at.color == RED {
					replacingParent.at.color = BLACK
					break
				}
				break
			} else {
				if effectiveColor(sibling.right) == BLACK {
					if sibling.left != nil {
						sibling.left.color = BLACK
					}
					sibling.color = RED
					replacingParent.hierarchy.Push(replacingParent.at)
					replacingParent.at = sibling
					t.rightRotateM(replacingParent)
					sibling = replacingParent.hierarchy.Pop()
					replacingParent.at = replacingParent.hierarchy.Pop()
				}
				sibling.color = replacingParent.at.color
				replacingParent.at.color = BLACK
				if sibling.right != nil {
					sibling.right.color = BLACK
				}
				replacingParent.hierarchy.Pop()
				t.leftRotateM(replacingParent)
				break
			}
		} else if anchor.lastremove == right {
			if effectiveColor(replacingParent.at.right) == RED {
				replacingParent.at.right.color = BLACK
				return
			}
			sibling := replacingParent.at.left
			if sibling == nil {
				//Sorry, I don't know how to fix it gracefully
				return
			}
			sibling = t.guaranteeWriteAccess(sibling, nil)
			replacingParent.at.left = sibling
			if sibling.color == RED {
				sibling.color = BLACK
				replacingParent.at.color = RED
				//Need Debug
				//replacingParent.hierarchy.Push(replacingParent.at)
				replacingParent.hierarchy.Pop()
				replacingParent = t.rightRotateM(replacingParent)
				replacingParent.hierarchy.Push(replacingParent.at)
				//replacingParent.at = replacingParent.hierarchy.Pop()
				sibling = replacingParent.at.left
			}
			if sibling == nil {
				//Sorry, I don't know how to fix it gracefully
				return
			}
			if effectiveColor(sibling.left) == BLACK && effectiveColor(sibling.right) == BLACK {
				sibling.color = RED
				replacingParent.at = replacingParent.hierarchy.Pop()
				anchor = replacingParent
				if replacingParent.at.color == RED {
					replacingParent.at.color = BLACK
					break
				}
				break
			} else {
				if effectiveColor(sibling.left) == BLACK {
					if sibling.right != nil {
						sibling.right.color = BLACK
					}
					sibling.color = RED
					replacingParent.hierarchy.Push(replacingParent.at)
					replacingParent.at = sibling
					t.leftRotateM(replacingParent)
					sibling = replacingParent.hierarchy.Pop()
					replacingParent.at = replacingParent.hierarchy.Pop()
				}
				sibling.color = replacingParent.at.color
				replacingParent.at.color = BLACK
				if sibling.left != nil {
					sibling.left.color = BLACK
				}
				replacingParent.hierarchy.Pop()
				t.rightRotateM(replacingParent)
				break
			}
		} else {
			runtime.Breakpoint()
		}
	}

}

/*
func (t *FRBTreeGKeyXXGValue) Verify() {
	return
	if t.root != nil && t.root.color == RED {
		runtime.Breakpoint()
	}
	t.verify(t.root, 1)
}

func (t *FRBTreeGKeyXXGValue) verify(v *frbtNodeGKeyXXGValue, blackabove uint) uint {
	if v == nil {
		return blackabove + 1
	}
	if v.color == RED {
		if !(effectiveColor(v.left) != RED && effectiveColor(v.right) != RED) {
			runtime.Breakpoint()
		}
	}
	if v.color == BLACK {
		blackabove++
	}
	lb := t.verify(v.left, blackabove)
	rb := t.verify(v.right, blackabove)
	if lb != rb {
		//runtime.Breakpoint()
	}
	return lb
}
*/
func (t *FRBTreeGKeyXXGValue) Walk() {
	t.walk(t.root)
}

func (t *FRBTreeGKeyXXGValue) walk(v *frbtNodeGKeyXXGValue) {
	if v == nil {
		return
	}
	t.walk(v.left)
	t.walk(v.right)
	print(*v.key)
	return
}

func (t *FRBTreeGKeyXXGValue) leftRotateM(anchor frbtAnchorGKeyXXGValue) frbtAnchorGKeyXXGValue {
	t.guaranteeAncestorsWriteAccess(anchor)
	pGrave := t.guaranteeWriteAccess(anchor.at, nil)
	qGrave := t.guaranteeWriteAccess(anchor.at.right, nil)

	b := anchor.at.right.left

	pGrave.right = b
	qGrave.left = pGrave
	if anchor.hierarchy.Len() == 0 {
		t.root = qGrave
	} else {
		if anchor.hierarchy.Peek().left == anchor.at {
			anchor.hierarchy.Peek().left = qGrave
		} else {
			anchor.hierarchy.Peek().right = qGrave
		}
	}

	anchor.hierarchy.Push(qGrave)

	return anchor

}

func (t *FRBTreeGKeyXXGValue) rightRotateM(anchor frbtAnchorGKeyXXGValue) frbtAnchorGKeyXXGValue {
	t.guaranteeAncestorsWriteAccess(anchor)
	qGrave := t.guaranteeWriteAccess(anchor.at, nil)
	pGrave := t.guaranteeWriteAccess(anchor.at.left, nil)

	b := anchor.at.left.right

	pGrave.right = qGrave
	qGrave.left = b

	if anchor.hierarchy.Len() == 0 {
		t.root = pGrave
	} else {
		if anchor.hierarchy.Peek().left == anchor.at {
			anchor.hierarchy.Peek().left = pGrave
		} else {
			anchor.hierarchy.Peek().right = pGrave
		}
	}

	anchor.hierarchy.Push(pGrave)
	anchor.at = qGrave

	return anchor
}

func (t *FRBTreeGKeyXXGValue) replacetreeelement(u *frbtAnchorGKeyXXGValue, v *frbtNodeGKeyXXGValue) {
	t.guaranteeAncestorsWriteAccess(*u)
	if u.hierarchy.Len() == 0 {
		t.root = v
		u.lastremove = root
	} else if u.hierarchy.Peek().left == u.at {
		up := u.hierarchy.Pop()
		up = t.guaranteeWriteAccess(up, nil)
		up.left = v
		u.hierarchy.Push(up)
		u.lastremove = left
	} else {
		up := u.hierarchy.Pop()
		up = t.guaranteeWriteAccess(up, nil)
		up.right = v
		u.hierarchy.Push(up)
		u.lastremove = right
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
	currentnew := a.at
checkfor:
	for {
		if a.hierarchy.Len() == 0 {
			t.root = updatestack.Peek()
			break checkfor
		}
		old := a.hierarchy.Pop()
		if !t.isShifted(old) {
			a.hierarchy.Push(old)
			break checkfor
		}
		new := t.dupNode(old)
		if old.left == current {
			new.left = currentnew
		} else if old.right == current {
			new.right = currentnew
		} else {
			if current != nil {
				runtime.Breakpoint()
			}
		}
		current = old
		currentnew = new
		updatestack.Push(new)
	}

	for updatestack.Len() != 0 {
		a.hierarchy.Push(updatestack.Pop())
	}

}

func (t *FRBTreeGKeyXXGValue) guaranteeWriteAccess(src *frbtNodeGKeyXXGValue, parent *frbtNodeGKeyXXGValue) *frbtNodeGKeyXXGValue {
	if t.isShifted(src) {
		new := t.dupNode(src)
		if parent != nil {
			if src == parent.left {
				parent.left = new
			} else if src == parent.right {
				parent.right = new
			} else {
				runtime.Breakpoint()
			}
		}
		return new
	}
	return src
}

func (t *FRBTreeGKeyXXGValue) Get(key GKey) (GValue, error) {
	result := t.narrowto(key, false)
	if result.at == nil {
		return nil, ErrKeyNotFoundGKeyXXGValue
	}
	return result.at.value, nil
}

/*narrow down to the nearest node
  return true if exact match is found,
              with an stack topped with result and parents hierarchy
         false if no exact match is found,
              with an stack topped with would be parents hierarchy
*/
func (t *FRBTreeGKeyXXGValue) narrowto(key GKey, wishnew bool) frbtAnchorGKeyXXGValue {
	hierarchystack := newstackGKeyXXGValue()
	current := t.root
	var counter = 30
	for {
		counter--
		if counter == 0 {
			runtime.Breakpoint()
		}
		if current == nil {
			return frbtAnchorGKeyXXGValue{at: nil, hierarchy: hierarchystack}
		}
		hierarchystack.Push(current)
		if t.lessthan(key, *current.key) {
			current = current.left
		} else {
			if !wishnew && *current.key == key {
				hierarchystack.Pop()
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
