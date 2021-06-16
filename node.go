package avl

// Item holds a piece of information needed to be stored (or searched by) in a
// tree.
//
// It's common to use different Item types for store and lookup while all of
// the types are consistent in comparisons:
//
//	type User struct {
//		ID int
//	}
//
//	func (u User) Compare(x avl.Item) int {
//		return u.ID - x.(User).ID
//	}
//
//	type ID int
//
//	func (id ID) Compare(x avl.Item) int {
//		return int(id) - x.(User).ID
//	}
//
//	tree, _ = tree.Insert(User{ID: 42})
//	user := tree.Search(ID(42))
//
// That is, Item can represent both the key for searching and value for storing
// (or searching).
type Item interface {
	// Compare compares item itself with another item usually stored in a tree.
	// It reports whether the receiver is less, greater or equal to the given
	// Item by returning values less than, greater than or equal to zero
	// respectively.
	Compare(Item) int
}

// node is a node of a tree.
type node struct {
	value Item
	left  *node
	right *node
	h     int // Subtree height.
}

// Size returns the size of a subtree rooted by n.
// Note that this method runs in O(N) to not bring additional O(N) space
// penalty to store the size field at each node.
func (n *node) Size() int {
	if n == nil {
		return 0
	}
	return 1 + n.left.Size() + n.right.Size()
}

// Insert inserts a new node with value x in the tree.
// It returns new tree root node if insertion happened or already existing node
// having the same value, meaning x was not inserted.
func (n *node) Insert(x Item) (root *node, existing Item) {
	if n == nil {
		return &node{
			value: x,
			h:     1,
		}, nil
	}
	cmp := x.Compare(n.value)
	switch {
	case cmp < 0:
		var m *node
		m, existing = n.left.Insert(x)
		if existing == nil {
			root = n.clone()
			root.left = m
		}
	case cmp > 0:
		var m *node
		m, existing = n.right.Insert(x)
		if existing == nil {
			root = n.clone()
			root.right = m
		}
	default:
		existing = n.value
	}
	if root == nil {
		// x is not inserted.
		return n, existing
	}

	root.adjustHeight()

	return root.rebalance(), nil
}

// Update updates a node having value x in the tree.
// It replaces the value of a node if it already exists in the tree or inserts
// new one with value x. It returns new tree root and an old value if it
// was present in the tree and replaced by x.
func (n *node) Update(x Item) (root *node, prev Item) {
	if n == nil {
		return &node{
			value: x,
			h:     1,
		}, nil
	}
	root = n.clone()
	cmp := x.Compare(n.value)
	switch {
	case cmp < 0:
		root.left, prev = n.left.Update(x)
	case cmp > 0:
		root.right, prev = n.right.Insert(x)
	default:
		root.value, prev = x, root.value
	}

	root.adjustHeight()

	return root.rebalance(), nil
}

// Delete deletes a node having value x from the tree.
// It returns new tree root node and a value of deleted node if such node was
// present in the tree. Otherwise it returns n and nil.
func (n *node) Delete(x Item) (root *node, existed Item) {
	if n == nil {
		return nil, nil
	}
	cmp := x.Compare(n.value)
	switch {
	case cmp < 0:
		var m *node
		m, existed = n.left.Delete(x)
		if existed != nil {
			root = n.clone()
			root.left = m
		}
	case cmp > 0:
		var m *node
		m, existed = n.right.Delete(x)
		if existed != nil {
			root = n.clone()
			root.right = m
		}
	default:
		root = n.destroy()
		existed = n.value
	}
	if existed == nil {
		// x is not present in n.
		return n, nil
	}
	if root == nil {
		// x was the last element of n.
		return nil, existed
	}

	root.adjustHeight()

	return root.rebalance(), existed
}

// Max returns max value of the tree.
func (n *node) Max() Item {
	if n == nil {
		return nil
	}
	if n.right != nil {
		return n.right.Max()
	}
	return n.value
}

// Min returns min value of the tree.
func (n *node) Min() Item {
	if n == nil {
		return nil
	}
	if n.left != nil {
		return n.left.Min()
	}
	return n.value
}

// Search searches for a node having value x and return its value.
// Note that x and node's value essentially can be a different types sharing
// comparison logic.
func (n *node) Search(x Item) Item {
	if n == nil {
		return nil
	}
	cmp := x.Compare(n.value)
	switch {
	case cmp < 0:
		return n.left.Search(x)
	case cmp > 0:
		return n.right.Search(x)
	default:
		return n.value
	}
}

// Predcessor finds a node which is in-order predcessor of a node having value
// x. It returns value of found node or nil.
func (n *node) Predcessor(x Item) Item {
	if n == nil {
		return nil
	}
	cmp := x.Compare(n.value)
	switch {
	case cmp < 0:
		return n.left.Predcessor(x)
	case cmp > 0:
		p := n.right.Predcessor(x)
		if p == nil {
			p = n.value
		}
		return p
	default:
		return n.left.Max()
	}
}

// Successor finds a node which is in-order successor of a node having value x.
// It returns value of found node or nil.
func (n *node) Successor(x Item) Item {
	if n == nil {
		return nil
	}
	cmp := x.Compare(n.value)
	switch {
	case cmp < 0:
		s := n.left.Successor(x)
		if s == nil {
			s = n.value
		}
		return s
	case cmp > 0:
		return n.right.Successor(x)
	default:
		return n.right.Min()
	}
}

// InOrder prepares in-order traversal of the tree and calls fn with value of
// each visited node.
func (n *node) InOrder(fn func(Item) bool) {
	if n == nil {
		return
	}
	n.left.InOrder(fn)
	fn(n.value)
	n.right.InOrder(fn)
}

// PreOrder prepares pre-order traversal of the tree and calls fn with value of
// each visited node.
func (n *node) PreOrder(fn func(Item) bool) {
	if n == nil {
		return
	}
	fn(n.value)
	n.left.PreOrder(fn)
	n.right.PreOrder(fn)
}

// PostOrder prepares post-order traversal of the tree and calls fn with value
// of each visited node.
func (n *node) PostOrder(fn func(Item) bool) {
	if n == nil {
		return
	}
	n.left.PostOrder(fn)
	n.right.PostOrder(fn)
	fn(n.value)
}

func (n *node) destroy() *node {
	switch {
	case n.left != nil && n.right != nil:
		//    (a)           e
		//    / \          / \
		//   b   c  =>    b   c
		//  / \          /
		// d  [e]       d
		m := n.left.Max()

		root := new(node)
		root.value = m
		root.left, _ = n.left.Delete(m)
		root.right = n.right

		return root

	case n.left != nil:
		return n.left

	case n.right != nil:
		return n.right

	default:
		return nil
	}
}

func (n *node) adjustHeight() {
	n.h = max(n.left.height(), n.right.height()) + 1
}

func (n *node) height() int {
	if n == nil {
		return 0
	}
	return n.h
}

func (n *node) balance() int {
	if n == nil {
		return 0
	}
	return n.right.height() - n.left.height()
}

func (n *node) rebalance() (root *node) {
	// b is greater than 1 when tree is right-heavy.
	// b is less than -1 when tree is left-heavy.
	// note that balance is simply right.height() - left.height().
	b := n.balance()
	switch {
	case b < -1 && n.left.balance() <= 0:
		//     (a)      b
		//     /       / \
		//    b   =>  c   a
		//   /
		//  c
		return n.rotateRight()

	case b > 1 && n.right.balance() >= 0:
		//  (a)           b
		//    \          / \
		//     b    =>  a   c
		//      \
		//       c
		return n.rotateLeft()

	case b < -1 && n.left.balance() > 0:
		//     a        (a)        b
		//    /         /         / \
		//  (c)   =>   b     =>  c   a
		//    \       /
		//     b     c
		n = n.clone()
		n.left = n.left.rotateLeft()
		return n.rotateRight()

	case b > 1 && n.right.balance() < 0:
		//  a       (a)           b
		//   \        \          / \
		//   (c) =>    b    =>  a   c
		//   /          \
		//  b            c
		n = n.clone()
		n.right = n.right.rotateRight()
		return n.rotateLeft()

	case b > 1 || b < -1:
		panic("avl: internal error: balancing error")
	}
	return n
}

func (n *node) rotateRight() *node {
	//     (a)        b
	//     / \       / \
	//    b   c =>  d   a
	//   / \           / \
	//  d   e         e   c
	root := n.left.clone()
	node := n.clone()
	node.left = root.right
	root.right = node

	node.adjustHeight()
	root.adjustHeight()

	return root
}

func (n *node) rotateLeft() *node {
	//     c         (a)
	//    / \        / \
	//   a   e  <=  b   c
	//  / \            / \
	// b   d          d   e
	root := n.right.clone()
	node := n.clone()
	node.right = root.left
	root.left = node

	node.adjustHeight()
	root.adjustHeight()

	return root
}

func (n *node) clone() *node {
	if n == nil {
		return nil
	}
	cp := *n
	return &cp
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
