package avl

// Tree is an immutable container holding root of an AVL tree.
// Modifying operations (Insert(), Update() and Delete()) are immutable and
// return copy of the tree.
//
// Note that Tree holds pointer to the root of an AVL tree internally, which
// makes Tree so called reference type. That is, there is no cases when you may
// need to pass pointer to instance of the Tree.
type Tree struct {
	root *node
	size int
}

// Size returns the size of a tree.
// The time complexity is O(1).
func (t Tree) Size() int {
	return t.size
}

// Insert inserts a new node with value x in the tree.
// It returns a copy of the tree and already existing item, which non-nil value
// means that x was not inserted.
func (t Tree) Insert(x Item) (_ Tree, existing Item) {
	t.root, existing = t.root.Insert(x)
	if existing == nil {
		t.size++
	}
	return t, existing
}

// Update updates a node having value x in the tree.
// It replaces the value of a node in the tree if it already exists or inserts
// new one with value x. It returns a copy of the tree and an old value if it
// was present and replaced by x.
func (t Tree) Update(x Item) (_ Tree, prev Item) {
	t.root, prev = t.root.Update(x)
	if prev == nil {
		t.size++
	}
	return t, prev
}

// Delete deletes a node having value x from the tree.
// It returns a copy of the tree and a value of deleted node if such node was
// present.
func (t Tree) Delete(x Item) (_ Tree, existed Item) {
	t.root, existed = t.root.Delete(x)
	if existed != nil {
		t.size--
	}
	return t, existed
}

// Max returns max value of the tree.
func (t Tree) Max() Item {
	return t.root.Max()
}

// Min returns min value of the tree.
func (t Tree) Min() Item {
	return t.root.Min()
}

// Search searches for a node having value x and return its value.
// Note that x and node's value essentially can be a different types sharing
// comparison logic.
func (t Tree) Search(x Item) Item {
	return t.root.Search(x)
}

// Predecessor finds a node in the tree which is an in-order predecessor of a
// node having value x. It returns value of found node or nil.
func (t Tree) Predecessor(x Item) Item {
	return t.root.Predecessor(x)
}

// Successor finds a node in the tree which is an in-order successor of a node
// having value x. It returns value of found node or nil.
func (t Tree) Successor(x Item) Item {
	return t.root.Successor(x)
}

// InOrder prepares in-order traversal of the tree and calls fn with value of
// each visited node. If fn returns false it stops traversal.
func (t Tree) InOrder(fn func(Item) bool) {
	t.root.InOrder(fn)
}

// PreOrder prepares pre-order traversal of the tree and calls fn with value of
// each visited node. If fn returns false it stops traversal.
func (t Tree) PreOrder(fn func(Item) bool) {
	t.root.PreOrder(fn)
}

// PostOrder prepares post-order traversal of the tree and calls fn with value
// of each visited node. If fn returns false it stops traversal.
func (t Tree) PostOrder(fn func(Item) bool) {
	t.root.PostOrder(fn)
}
