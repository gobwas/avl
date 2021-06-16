package avl

import "fmt"

func ExampleTree() {
	var tree Tree
	tree, _ = tree.Insert(IntItem(1))
	tree, _ = tree.Insert(IntItem(2))
	tree, _ = tree.Insert(IntItem(3))
	tree, _ = tree.Insert(IntItem(4))
	tree, _ = tree.Delete(IntItem(4))
	tree.InOrder(func(x Item) bool {
		fmt.Print(x, " ")
		return true
	})
	// Output:
	// 1 2 3
}
