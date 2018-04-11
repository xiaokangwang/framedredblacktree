package framedredblacktree

import "testing"

import "os"

import "math/rand"

//import "github.com/davecgh/go-spew/spew"

var (
	testput, _ = os.Create("/home/shelikhoo/proj/src/github.com/xiaokangwang/framedredblacktree/test.txt")
)

func TestInit(t *testing.T) {
	tree := NewFRBTreeGKeyXXGValue(func(p1, p2 GKey) bool {
		return (p1) < (p2)
	})
	tree.Insert(GKey(1), 1)
	tree.Insert(GKey(3), 3)
	tree.Insert(GKey(4), 4)
	tree.Insert(GKey(5), 5)
	tree.Insert(GKey(2), 2)
	tree.Insert(GKey(7), 7)
	tree.Insert(GKey(8), 8)
	tree.Insert(GKey(6), 6)
	tree.Insert(GKey(12), 12)
	//spew.Fdump(testput, tree)
	tree.Insert(GKey(1), 13)
	tree.Insert(GKey(7), 14)
	tree.Insert(GKey(9), 15)
}

func TestNoDataLoss(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	tree := NewFRBTreeGKeyXXGValue(func(p1, p2 GKey) bool {
		return (p1) < (p2)
	})
	var current uint
	current = 1000
	for current < 2000 {
		tree.Insert(GKey(current), current)
		_, err := tree.Get(GKey(current))
		if err != nil {
			t.Fatalf("err is %s %v", err, current)
		}
		current++
	}
	current = 1000
	for current > 0 {
		tree.Insert(GKey(current), current)
		_, err := tree.Get(GKey(current))
		if err != nil {
			t.Fatalf("err is %s %v", err, current)
		}
		current--
	}

	current = 240000

	for current > 0 {
		thistry := uint(rand.Int())
		tree.Insert(GKey(thistry), thistry)
		_, err := tree.Get(GKey(thistry))
		if err != nil {
			t.Fatalf("err is %s %v", err, current)
		}
		current--
	}

}

func TestNoDataLossEx(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	tree := NewFRBTreeGKeyXXGValue(func(p1, p2 GKey) bool {
		return (p1) < (p2)
	})
	var current uint

	current = 240000
	reference := make([]uint, 0, 240000)
	for current > 0 {
		thistry := uint(rand.Int())
		tree.Insert(GKey(thistry), thistry)
		reference = append(reference, thistry)
		current--
	}
	for _, testvar := range reference {
		_, err := tree.Get(GKey(testvar))
		if err != nil {
			t.Fatalf("err is %s %v", err, testvar)
		}
	}

}

func TestDelInit(t *testing.T) {

	tree := NewFRBTreeGKeyXXGValue(func(p1, p2 GKey) bool {
		return (p1) < (p2)
	})
	tree.Insert(GKey(1), 1)
	tree.Insert(GKey(2), 2)
	tree.Insert(GKey(3), 3)
	tree.Insert(GKey(4), 4)
	tree.Insert(GKey(5), 5)
	tree.Drop(GKey(1))
	tree.Drop(GKey(5))
	//spew.Fdump(testput, tree)
}

func TestNoDataLossD(t *testing.T) {
	t.SkipNow()
	if testing.Short() {
		t.SkipNow()
	}
	tree := NewFRBTreeGKeyXXGValue(func(p1, p2 GKey) bool {
		return (p1) < (p2)
	})
	var current uint
	current = 1000
	for current < 2000 {
		tree.Insert(GKey(current), current)
		_, err := tree.Get(GKey(current))
		if err != nil {
			t.Fatalf("err is %s %v", err, current)
		}
		if current%3 == 0 || current%7 == 0 {
			tree.Drop(GKey(current))
			_, err = tree.Get(GKey(current))
			if err == nil {
				t.Fatalf("err is %s %v", err, current)
			}
		}
		current++
	}
	current = 1000
	for current > 0 {
		tree.Insert(GKey(current), current)
		_, err := tree.Get(GKey(current))
		if err != nil {
			t.Fatalf("err is %s %v", err, current)
		}
		if current%3 == 0 || current%7 == 0 {
			tree.Drop(GKey(current))
			_, err = tree.Get(GKey(current))
			if err == nil {
				t.Fatalf("err is %s %v", err, current)
			}
		}
		current--
	}

}

func TestNoDataLossExD(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	tree := NewFRBTreeGKeyXXGValue(func(p1, p2 GKey) bool {
		return (p1) < (p2)
	})
	var current uint

	current = 240000
	rand.Seed(671) //671 342
	reference := make([]uint, 0, 100)
	del := make(map[uint]int)
	for current > 0 {
		thistry := uint(rand.Int())
		tree.Insert(GKey(thistry), thistry)
		//tree.Verify()
		reference = append(reference, thistry)
		current--
		//spew.Fdump(testput, tree)
	}
	//spew.Fdump(testput, reference)
	//spew.Fdump(testput, tree)

	for _, testvar := range reference {
		if rand.Float32() > 0.3 {
			//spew.Fdump(testput, GKey(testvar))
			err := tree.Drop(GKey(testvar))
			//spew.Fdump(testput, tree)
			//tree.Verify()
			if err != nil {
				t.Fatalf("err is %s %v %v", err, testvar, del[testvar])
			}
			sum, isdeled := del[testvar]
			if isdeled {
				del[testvar] = sum + 1
			} else {
				del[testvar] = 1
			}
		}

	}

	for _, testvar := range reference {
		_, err := tree.Get(GKey(testvar))
		_, isdeled := del[testvar]
		if isdeled {
			if err == nil {
				//t.Fatalf("err is %s %v", err, testvar)
			}
		} else {
			if err != nil {
				t.Fatalf("err is %s %v", err, testvar)
			}
		}

	}

}

func TestNoDataLossDiversified(t *testing.T) {
	if testing.Short() {
		//t.SkipNow()
	}
	tree := NewFRBTreeGKeyXXGValue(func(p1, p2 GKey) bool {
		return (p1) < (p2)
	})
	tree.Insert(GKey(1), 1)
	tree.Insert(GKey(2), 2)
	tree.Insert(GKey(3), 3)
	tree1 := tree.Diversify()
	tree1.Insert(GKey(4), 4)
	tree1.Insert(GKey(5), 5)
	tree1.Insert(GKey(6), 6)
	tree2 := tree.Diversify()
	tree2.Insert(GKey(7), 7)
	tree2.Insert(GKey(8), 8)
	tree2.Insert(GKey(9), 9)
	tree.Walk()
	println()
	tree1.Walk()
	println()
	tree2.Walk()
	println()
	tree2.Drop(GKey(8))
	tree2.Walk()
}

func subTestNoDataLossExDFork(t *testing.T, tree *FRBTreeGKeyXXGValue, gen uint) {

	var current uint

	current = 100000
	reference := make([]uint, 0, 100)
	del := make(map[uint]int)
	for current > 0 {
		thistry := uint(rand.Int())
		tree.Insert(GKey(thistry), thistry)
		//tree.Verify()
		reference = append(reference, thistry)
		current--
		//spew.Fdump(testput, tree)
	}
	//spew.Fdump(testput, reference)
	//spew.Fdump(testput, tree)
	/*
		for _, testvar := range reference {
			if rand.Float32() > 0.3 {
				//spew.Fdump(testput, GKey(testvar))
				err := tree.Drop(GKey(testvar))
				//spew.Fdump(testput, tree)
				//tree.Verify()
				if err != nil {
					t.Fatalf("err is %s %v %v", err, testvar, del[testvar])
				}
				sum, isdeled := del[testvar]
				if isdeled {
					del[testvar] = sum + 1
				} else {
					del[testvar] = 1
				}
			}

		}
	*/
	if gen < 3 {
		subTestNoDataLossExDFork(t, tree.Diversify(), gen+1)
		subTestNoDataLossExDFork(t, tree.Diversify(), gen+1)
		subTestNoDataLossExDFork(t, tree.Diversify(), gen+1)
	}

	for _, testvar := range reference {
		_, err := tree.Get(GKey(testvar))
		_, isdeled := del[testvar]
		if isdeled {
			if err == nil {
				//t.Fatalf("err is %s %v", err, testvar)
			}
		} else {
			if err != nil {
				t.Fatalf("err is %s %v", err, testvar)
			}
		}

	}

}

func TestDiversifiedMutiDim(t *testing.T) {
	//t.SkipNow()
	rand.Seed(12997)
	if testing.Short() {
		//t.SkipNow()
	}
	tree := NewFRBTreeGKeyXXGValue(func(p1, p2 GKey) bool {
		return (p1) < (p2)
	})
	subTestNoDataLossExDFork(t, tree, 0)
}
