package tree

// PreOrderRecur 递归前序遍历
func PreOrderRecur(root *TreeNode) (treeNode []*TreeNode) {
	if root == nil {
		return nil
	}

	treeNode = append(treeNode, root)
	PreOrderRecur(root.Left)
	PreOrderRecur(root.Right)
	return
}

// PreOrderUnRecur 非递归前序遍历
func PreOrderUnRecur(root *TreeNode) (treeNode []*TreeNode) {
	if root == nil {
		return
	}
	stack := make([]*TreeNode, 0)
	stack = append(stack, root)
	for len(stack) > 0 {
		curNode := stack[len(stack)-1]
		treeNode = append(treeNode, curNode)

		stack = stack[0:len(stack)-1]
		if curNode.Right != nil {
			stack = append(stack, curNode.Left)
		}
		if curNode.Left != nil {
			stack = append(stack, curNode.Left)
		}
	}
	return treeNode
}

func InOrderRecur(root *TreeNode) (treeNode []*TreeNode) {
	if root == nil {
		return nil
	}

	InOrderRecur(root.Left)
	treeNode = append(treeNode, root)
	InOrderRecur(root.Right)
	return
}

func InOrderUnRecur(root *TreeNode) (treeNode []*TreeNode) {
	if root == nil {
		return nil
	}
	cur := root
	stack := make([]*TreeNode, 0)
	stack = append(stack, cur)

	for len(stack) > 0 {
		if cur.Left != nil {
			stack = append(stack, cur.Left)
			continue
		}
		curNode := stack[len(stack)-1]
		stack = stack[0:len(stack)-1]
		treeNode = append(treeNode, curNode)
		if curNode.Right != nil {
			stack = append(stack, curNode.Right)
		}
	}
	return
}


func PostOrderRecur(root *TreeNode) (treeNode []*TreeNode) {
	if root == nil {
		return nil
	}

	PostOrderRecur(root.Left)
	PostOrderRecur(root.Right)
	treeNode = append(treeNode, root)
	return
}

func PostOrderUnRecur(root *TreeNode) (treeNode []*TreeNode) {
	if root == nil {
		return nil
	}

	helpStack := make([]*TreeNode, 0)
	resultStack := make([]*TreeNode, 0)

	helpStack = append(helpStack, root)
	for len(helpStack) > 0 {
		curNode := helpStack[len(helpStack)-1]
		helpStack = helpStack[0:len(helpStack)-1]
		if curNode.Left != nil {
			helpStack = append(helpStack, curNode.Left)
		}
		if curNode.Right != nil {
			helpStack = append(helpStack, curNode.Right)
		}
		resultStack = append(resultStack, curNode)
	}
	return
}



