/*
Package avl implements immutable AVL (Adelson-Velsky and Landis) tree.

Immutability means that on any tree modifying operation (insert, update or
delete) it does clone affected path up to the tree root. Since AVL is balanced
binary tree, immutability property leads to O(log n) additional allocations.

Immutability lets applications to update the tree without holding a lock for
the whole time span of operation. That is, it's possible to read the current
tree "state", then update it locally and change the tree "state" in atomic way
later. This technique lets other goroutines to read the tree contents without
blocking while update operation is in progress:

	var (
		mu   sync.RWMutex
		tree avl.Tree
	)
	writer := func() {
		for {
			// Read the tree state while holding read-lock, which means that
			// read-only goroutines are not blocked.
			mu.RLock()
			t := tree
			mu.RUnlock()

			// Modify the tree and update the t variable holding immutable
			// state.
			t, _ = t.Insert(x)
			t, _ = t.Delete(y)

			// Update the tree state while holding write-lock, which means that
			// read-only goroutines have to wait until we finish.
			mu.Lock()
			tree = t
			mu.Unlock()
		}
	}
	reader := func() {
		for {
			// Read the tree state while holding read-lock.
			mu.RLock()
			{
				// Make any read-only calls on tree such as Search() or
				// Successor() etc.
			}
			mu.RUnlock()
		}
	}

	go reader()
	go reader()
	go writer()

Note that usually there is a need to use second mutex to serialize tree updates
across multiple writer goroutines.
*/
package avl
