package avl

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func TestDeleteLastItem(t *testing.T) {
	item := IntItem(1)
	n0, _ := (*node)(nil).Insert(item)
	n1, _ := n0.Delete(item)
	if n1 != nil {
		t.Fatalf("unexpected non-nil root node after deletion")
	}
}

func TestDeleteNonExistingItem(t *testing.T) {
	n0, _ := (*node)(nil).Insert(IntItem(0))
	n1, _ := n0.Delete(IntItem(1))
	if n1 != n0 {
		t.Fatalf("unexpected root node after deletion")
	}
}

func TestInsertDuplicate(t *testing.T) {
	for _, test := range []struct {
		name    string
		init    []int
		insert  []int
		inOrder []int
	}{
		{
			init:    []int{1, 2, 3},
			insert:  []int{1, 2, 3},
			inOrder: []int{1, 2, 3},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := buildTree(t, test.init, nil)
			for _, x := range test.insert {
				var existing Item
				root, existing = root.Insert(IntItem(x))
				if existing == nil {
					t.Errorf("inserting %d: no duplicate", x)
				}
			}
			assertInOrder(t, root, test.inOrder)
		})
	}
}

func BenchmarkInsert(b *testing.B) {
	for _, test := range []struct {
		name   string
		init   []int
		rand   int
		insert []int
		delete []int
	}{
		{
			name:   "no rebalance",
			init:   []int{1, 2, 3, 5, 6, 7},
			insert: []int{8},
		},
		{
			name:   "rebalance",
			init:   []int{1, 2, 3, 5, 6, 7, 8},
			insert: []int{9},
		},
		{
			name:   "big",
			rand:   1 << 20,
			insert: []int{42},
		},
	} {
		b.Run(test.name, func(b *testing.B) {
			root := buildTree(b, test.init, nil)

			// Ignore values that we are going to insert.
			ignore := make(map[int]bool, len(test.insert))
			for _, n := range test.insert {
				ignore[n] = true
			}
			// Generate test.rand number of random values.
			for i := 0; i < test.rand; i++ {
				for {
					x := rand.Intn(math.MaxInt32)
					if ignore[x] {
						continue
					}
					var existing Item
					root, existing = root.Insert(IntItem(x))
					if existing == nil {
						break
					}
				}
			}

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				temp := root
				for _, x := range test.insert {
					temp, _ = temp.Insert(IntItem(x))
				}
				for _, x := range test.delete {
					temp, _ = temp.Delete(IntItem(x))
				}
			}
		})
	}
}

func TestPredcessorSuccessor(t *testing.T) {
	for _, test := range []struct {
		name       string
		insert     []int
		lookup     int
		empty      bool
		successor  Item
		predcessor Item
	}{
		{
			//   2
			//  / \
			// 1   3
			insert:     []int{1, 2, 3},
			lookup:     2,
			predcessor: IntItem(1),
			successor:  IntItem(3),
		},
		{
			//   2
			//  / \
			// 1   3
			insert:     []int{1, 2, 3},
			lookup:     1,
			predcessor: nil,
			successor:  IntItem(2),
		},
		{
			//   2
			//  / \
			// 1   3
			insert:     []int{1, 2, 3},
			lookup:     3,
			predcessor: IntItem(2),
			successor:  nil,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := buildTree(t, test.insert, nil)
			p := root.Predcessor(IntItem(test.lookup))
			s := root.Successor(IntItem(test.lookup))
			if act, exp := p, test.predcessor; act != exp {
				t.Fatalf("unexpected predcessor: %s; want %s", act, exp)
			}
			if act, exp := s, test.successor; act != exp {
				t.Fatalf("unexpected successor: %s; want %s", act, exp)
			}
		})
	}
}

func TestBalance(t *testing.T) {
	for _, test := range []struct {
		name      string
		insert    []int
		delete    []int
		min       int
		max       int
		inOrder   []int
		preOrder  []int
		postOrder []int
	}{
		{
			name: "left-left",
			//    2
			//   / \
			//  1   4
			//     / \
			//    3   5
			insert:    []int{1, 2, 3, 4, 5},
			min:       1,
			max:       5,
			inOrder:   []int{1, 2, 3, 4, 5},
			preOrder:  []int{2, 1, 4, 3, 5},
			postOrder: []int{1, 3, 5, 4, 2},
		},
		{
			name: "left-left deletion",
			//      4
			//     / \
			//    2   5
			//     \
			//      3
			insert:    []int{1, 2, 3, 4, 5},
			delete:    []int{1},
			min:       2,
			max:       5,
			inOrder:   []int{2, 3, 4, 5},
			preOrder:  []int{4, 2, 3, 5},
			postOrder: []int{3, 2, 5, 4},
		},
		{
			name: "right-right",
			//      4
			//     / \
			//    2   5
			//   / \
			//  1   3
			insert:    []int{5, 4, 3, 2, 1},
			min:       1,
			max:       5,
			inOrder:   []int{1, 2, 3, 4, 5},
			preOrder:  []int{4, 2, 1, 3, 5},
			postOrder: []int{1, 3, 2, 5, 4},
		},
		{
			name: "right-right deletion",
			//    2
			//   / \
			//  1   4
			//     /
			//    3
			insert:    []int{5, 4, 3, 2, 1},
			delete:    []int{5},
			min:       1,
			max:       4,
			inOrder:   []int{1, 2, 3, 4},
			preOrder:  []int{2, 1, 4, 3},
			postOrder: []int{1, 3, 4, 2},
		},
		{
			name: "left-right",
			//   2
			//  / \
			// 1   3
			insert:    []int{3, 1, 2},
			min:       1,
			max:       3,
			inOrder:   []int{1, 2, 3},
			preOrder:  []int{2, 1, 3},
			postOrder: []int{1, 3, 2},
		},
		{
			name: "left-right deletion",
			//   3            2
			//  / \          / \
			// 1  [5]  x=>  1   3
			//  \
			//   2
			insert:    []int{3, 1, 5, 2},
			delete:    []int{5},
			min:       1,
			max:       3,
			inOrder:   []int{1, 2, 3},
			preOrder:  []int{2, 1, 3},
			postOrder: []int{1, 3, 2},
		},
		{
			name: "right-left",
			//   2
			//  / \
			// 1   3
			insert:    []int{1, 3, 2},
			min:       1,
			max:       3,
			inOrder:   []int{1, 2, 3},
			preOrder:  []int{2, 1, 3},
			postOrder: []int{1, 3, 2},
		},
		{
			name: "right-left deletion",
			//    1           2
			//   / \         / \
			// [0]  3  x=>  1   3
			//     /
			//    2
			insert:    []int{1, 0, 3, 2},
			delete:    []int{0},
			min:       1,
			max:       3,
			inOrder:   []int{1, 2, 3},
			preOrder:  []int{2, 1, 3},
			postOrder: []int{1, 3, 2},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := buildTree(t, test.insert, test.delete)
			assertMin(t, root, test.min)
			assertMax(t, root, test.max)
			assertInOrder(t, root, test.inOrder)
			assertPreOrder(t, root, test.preOrder)
			assertPostOrder(t, root, test.postOrder)
		})
	}
}

func buildTree(t testing.TB, insert, delete []int) *node {
	var root *node
	for _, n := range insert {
		var existing Item
		root, existing = root.Insert(IntItem(n))
		if existing != nil {
			t.Fatalf("malformed input: %d inserted already", n)
		}
	}
	for _, n := range delete {
		var prev Item
		root, prev = root.Delete(IntItem(n))
		if prev == nil {
			t.Fatalf("malformed input: %d wasn't inserted", n)
		}
	}
	return root
}

func assertItem(t *testing.T, name string, exp int, getter func() Item) {
	act := int(getter().(IntItem))
	if act != exp {
		t.Errorf(
			"unexpected %s value: %d; want %d",
			name, act, exp,
		)
	}
}

func assertMin(t *testing.T, root *node, exp int) {
	assertItem(t, "min", exp, root.Min)
}
func assertMax(t *testing.T, root *node, exp int) {
	assertItem(t, "max", exp, root.Max)
}

func assertOrder(t *testing.T, name string, exp []int, iterator func(func(Item) bool)) {
	var i int
	iterator(func(x Item) bool {
		act := int(x.(IntItem))
		if exp := exp[i]; act != exp {
			t.Errorf(
				"%s[%d]=%d; want %d",
				name, i, act, exp,
			)
		}
		i++
		return true
	})
	if n := len(exp); i != n {
		t.Errorf("unexpected traversed items count: %d; want %d", i, n)
	}
}

func assertInOrder(t *testing.T, root *node, exp []int) {
	assertOrder(t, "inOrder", exp, root.InOrder)
}
func assertPreOrder(t *testing.T, root *node, exp []int) {
	assertOrder(t, "preOrder", exp, root.PreOrder)
}
func assertPostOrder(t *testing.T, root *node, exp []int) {
	assertOrder(t, "postOrder", exp, root.PostOrder)
}

type IntItem int

func (a IntItem) Compare(b Item) int {
	return int(a) - int(b.(IntItem))
}

func (a IntItem) String() string {
	return fmt.Sprintf("%d", int(a))
}
