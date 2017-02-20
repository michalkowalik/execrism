// Package binarysearchtree provides a binary search tree implementation
package binarysearchtree

const testVersion = 1

// SearchTreeData implements the binary search tree
type SearchTreeData struct {
	left  *SearchTreeData
	data  int
	right *SearchTreeData
}

// Bst returns a new binary search tree
func Bst(d int) *SearchTreeData {
	return &SearchTreeData{data: d}
}

// Insert inserts new item to the binary tree
func (bt *SearchTreeData) Insert(d int) {
	bt = insertRec(bt, d)
}

func insertRec(node *SearchTreeData, d int) *SearchTreeData {
	if node == nil {
		return Bst(d)
	}
	if d <= node.data {
		node.left = insertRec(node.left, d)
	} else {
		node.right = insertRec(node.right, d)
	}
	return node
}

// MapString returns sorted values of all nodes mapped through "f" function
func (bt *SearchTreeData) MapString(f func(int) string) []string {
	sts := make([]string, 0)
	i := make([]int, 0)
	bt.traverse(&i)
	for _, item := range i {
		sts = append(sts, f(item))
	}
	return sts
}

// MapInt returns sorted values of all nodes mapped through "f" function
func (bt *SearchTreeData) MapInt(f func(int) int) []int {
	ints := make([]int, 0)
	bt.traverse(&ints)
	for i := range ints {
		ints[i] = f(ints[i])
	}
	return ints
}

func (bt *SearchTreeData) traverse(in *[]int) {
	if bt == nil {
		return
	}
	bt.left.traverse(in)
	*in = append(*in, bt.data)
	bt.right.traverse(in)
}
