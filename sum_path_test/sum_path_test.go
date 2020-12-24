package sum_path_test

import (
	"errors"
	"fmt"
	"testing"
)

type BTree struct {
	left  *BTree
	right *BTree
	value int
}

type Stack []*BTree

func (s *Stack) Push(t *BTree) {
	*s = append(*s, t)
}

func (s *Stack) Pop() (*BTree, error) {
	temp := *s
	if len(temp) == 0 {
		return nil, errors.New("can't pop an empty stack")
	}
	t := temp[len(temp)-1]
	*s = temp[:len(temp)-1]
	return t, nil
}

func (s Stack) Top() (*BTree, error) {
	if len(s) == 0 {
		return nil, errors.New("can't top an empty stack")
	}
	return s[len(s)-1], nil
}

func (s Stack) Len() int {
	return len(s)
}

func (s Stack) IsEmpty() bool {
	return len(s) == 0
}

func PrintArray(array []int, length int) {
	if length <= 0 || len(array) == 0 {
		return
	} else if length > len(array) {
		length = len(array)
	}

	fmt.Printf("%d", array[0])
	for i := 1; i < length; i++ {
		fmt.Printf("->%d", array[i])
	}
	fmt.Println()
}

func PrintSumPath(root *BTree, sum int) {
	if root == nil {
		return
	}

	// preparations
	var lastNode *BTree
	s := Stack{}
	s.Push(root)

	// resume that maximum depth of binary tree is 100
	l := make([]int, 100)
	i := 0

	currentSum := 0
	currentNode := root
	for {
		if s.IsEmpty() {
			break
		}

		if currentNode != nil {
			s.Push(currentNode)

			l[i] = currentNode.value
			i += 1

			currentSum += currentNode.value

			if currentSum == sum {
				PrintArray(l, i)
			}

			currentNode = currentNode.left
		} else {
			topNode, err := s.Top()
			if err != nil {
				return
			}

			if topNode.right != nil && topNode.right != lastNode {
				currentNode = topNode.right
			} else {
				lastNode = topNode
				_, err := s.Pop()
				if err != nil {
					return
				}
				if i > 0 {
					i--
				}

				currentSum -= topNode.value
			}
		}

	}
}

func TestPrintSumPath(t *testing.T) {
	root := &BTree{
		value: 2,
		left: &BTree{
			value: 7,
			left: &BTree{
				value: 2,
			},
			right: &BTree{
				value: 6,
				left: &BTree{
					value: 5,
				},
				right: &BTree{
					value: 11,
				},
			},
		},
		right: &BTree{
			value: 5,
			right: &BTree{
				value: 9,
				left: &BTree{
					value: 4,
				},
			},
		},
	}

	PrintSumPath(root, 15)
	PrintSumPath(root, 20)
}
