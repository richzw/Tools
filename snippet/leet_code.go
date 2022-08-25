package snippet

import "fmt"

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func isBST(root *TreeNode) bool {
	var pre *TreeNode
	var LDR func(node *TreeNode) bool
	LDR = func(node *TreeNode) bool {
		if node == nil {
			return true
		}
		leftResult := LDR(node.Left)
		if pre != nil && node.Val <= pre.Val {
			return false
		}
		pre = node
		rightResult := LDR(node.Right)
		return leftResult && rightResult
	}
	return LDR(root)
}

func PreOrder(root *TreeNode) {
	if root == nil {
		//fmt.Println("nil")
		return
	}
	PreOrder(root.Left)
	fmt.Println(root.Val)
	PreOrder(root.Right)
}
func addOneRow(root *TreeNode, val, depth int) *TreeNode {
	if root == nil {
		return nil
	}
	if depth == 1 {
		return &TreeNode{val, root, nil}
	}
	if depth == 2 {
		root.Left = &TreeNode{val, root.Left, nil}
		root.Right = &TreeNode{val, nil, root.Right}
	} else {
		root.Left = addOneRow(root.Left, val, depth-1)
		root.Right = addOneRow(root.Right, val, depth-1)
	}
	return root
}

func addOneRowV2(root *TreeNode, val, depth int) *TreeNode {
	if depth == 1 {
		return &TreeNode{val, root, nil}
	}
	nodes := []*TreeNode{root}
	for i := 1; i < depth-1; i++ {
		tmp := nodes
		nodes = nil
		for _, node := range tmp {
			if node.Left != nil {
				nodes = append(nodes, node.Left)
			}
			if node.Right != nil {
				nodes = append(nodes, node.Right)
			}
		}
	}
	for _, node := range nodes {
		node.Left = &TreeNode{val, node.Left, nil}
		node.Right = &TreeNode{val, nil, node.Right}
	}
	return root
}

func CreateTree() *TreeNode {
	root := &TreeNode{Val: 5}
	node4 := &TreeNode{Val: 4}
	node6 := &TreeNode{Val: 6}

	root.Left = node4
	root.Right = node6

	node1 := &TreeNode{Val: 1}
	node2 := &TreeNode{Val: 2}
	node3 := &TreeNode{Val: 3}
	node2.Left = node1
	node4.Left = node2
	node4.Right = node3

	node8 := &TreeNode{Val: 8}
	node6.Right = node8

	node9 := &TreeNode{Val: 9}
	node10 := &TreeNode{Val: 10}
	node8.Left = node9
	node8.Right = node10

	node13 := &TreeNode{Val: 13}
	node10.Right = node13

	return root
}

func PreOrderInteractive(root *TreeNode) {
	if root == nil {
		return
	}
	stack := make([]*TreeNode, 0)
	stack = append(stack, root)
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		fmt.Printf("%d, ", node.Val)

		if node.Right != nil {
			stack = append(stack, node.Right)
		}
		if node.Left != nil {
			stack = append(stack, node.Left)
		}
	}
}

func inOrderInteractive(root *TreeNode) {
	if root == nil {
		return
	}
	stack := make([]*TreeNode, 0)
	node := root

	for node != nil || len(stack) > 0 {
		if node != nil {
			stack = append(stack, node)
			node = node.Left
			continue
		}

		node = stack[len(stack)-1]
		fmt.Printf("%d, ", node.Val)
		stack = stack[:len(stack)-1]
		node = node.Right
	}
}

func postOrderInteractive(root *TreeNode) {
	if root == nil {
		return
	}
	stack := make([]*TreeNode, 0)
	node := root
	var prev *TreeNode

	for node != nil || len(stack) > 0 {
		if node != nil {
			stack = append(stack, node)
			node = node.Left
			continue
		}

		node = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if node.Right == nil || node.Right == prev {
			fmt.Printf("%d, ", node.Val)
			prev = node
			node = nil
		} else {
			stack = append(stack, node)
			node = node.Right
		}
	}
}

func postOrderInteractive2(root *TreeNode) {
	if root == nil {
		return
	}
	inp := make([]*TreeNode, 0)
	out := make([]int, 0)
	inp = append(inp, root)

	for len(inp) > 0 {
		node := inp[len(inp)-1]
		inp = inp[:len(inp)-1]

		out = append(out, node.Val)

		if node.Left != nil {
			inp = append(inp, node.Left)
		}
		if node.Right != nil {
			inp = append(inp, node.Right)
		}
	}

	for l := len(out) - 1; l >= 0; l-- {
		fmt.Printf("%d, ", out[l])
	}
}

func PostOrder(root *TreeNode) {
	if root == nil {
		//fmt.Println("nil")
		return
	}
	PostOrder(root.Left)
	PostOrder(root.Right)
	fmt.Printf("%d, ", root.Val)
}

type IterNode struct {
	line int
	node *TreeNode
}

func postOrderInteractive3(root *TreeNode) {
	if root == nil {
		return
	}
	stack := make([]*IterNode, 0)
	stack = append(stack, &IterNode{node: root, line: 0})

	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		if cur.line == 0 {
			if cur.node == nil {
				stack = stack[:len(stack)-1]
				continue
			}
		} else if cur.line == 1 {
			stack = append(stack, &IterNode{node: cur.node.Left, line: 0})
		} else if cur.line == 2 {
			stack = append(stack, &IterNode{node: cur.node.Right, line: 0})
		} else if cur.line == 3 {
			fmt.Printf("%d, ", cur.node.Val)
			stack = stack[:len(stack)-1]
		}

		cur.line++
	}
}

func inorderSuccessor(root *TreeNode, p *TreeNode) *TreeNode {
	return findSucc(root, p.Val)
}
func findSucc(root *TreeNode, val int) *TreeNode {
	var ans *TreeNode
	for root != nil {
		if val == root.Val {
			if root.Right != nil {
				p := root.Right
				for p.Left != nil {
					p = p.Left
				}
				return p
			}

			break
		}
		if root.Val > val && (ans == nil || ans.Val > root.Val) {
			ans = root
		}
		if val < root.Val {
			root = root.Left
		} else {
			root = root.Right
		}
	}
	return ans
}

func inorderTrave(root, succ *TreeNode, target int) {
	if root == nil {
		return
	}

	inorderTrave(root.Left, succ, target)
	if root.Val > target && succ.Left != nil {
		succ.Left = root
		return
	}
	inorderTrave(root.Right, succ, target)
}
func inorderSuccessorV2(root *TreeNode, target int) *TreeNode {
	succ := &TreeNode{Val: 0}
	inorderTrave(root, succ, target)
	return succ.Left
}

func getPar(root *TreeNode, pars map[int]*TreeNode) {
	if root == nil {
		return
	}

	if root.Left != nil {
		getPar(root.Left, pars)
		pars[root.Left.Val] = root
	}

	if root.Right != nil {
		getPar(root.Right, pars)
		pars[root.Right.Val] = root
	}
}

func testLeastAncestor(root, a, b *TreeNode) int {
	parMap := make(map[int]*TreeNode)
	getPar(root, parMap)

	path := make(map[int]bool)
	val := a.Val
	for {
		if v, ok := parMap[val]; ok {
			path[v.Val] = true
			val = v.Val
		} else {
			break
		}
	}

	val = b.Val
	for {
		if v, ok := parMap[val]; ok {
			if _, ok = path[v.Val]; ok {
				return v.Val
			} else {
				val = v.Val
			}
		} else {
			break
		}
	}

	return 0
}

func pathSu(root *TreeNode, target int) bool {
	if root == nil {
		return target == 0
	}
	if root.Left == nil && root.Right == nil {
		return root.Val == target
	}

	return pathSu(root.Left, target-root.Val) || pathSu(root.Right, target-root.Val)
}

func canFinish(n int, pre [][]int) bool {
	in := make([]int, n)
	frees := make([][]int, n)
	next := make([]int, 0, n)
	for _, v := range pre {
		in[v[0]]++
		frees[v[1]] = append(frees[v[1]], v[0])
	}
	for i := 0; i < n; i++ {
		if in[i] == 0 {
			next = append(next, i)
		}
	}
	for i := 0; i != len(next); i++ {
		c := next[i]
		v := frees[c]
		for _, vv := range v {
			in[vv]--
			if in[vv] == 0 {
				next = append(next, vv)
			}
		}
	}
	return len(next) == n
}

func tSort(g map[int][]int) []int {
	var linearOrder []int
	inDegree := map[int]int{}

	for n := range g {
		inDegree[n] = 0
	}

	for _, adjacent := range g {
		for _, v := range adjacent {
			inDegree[v]++
		}
	}

	next := []int{}
	for u, v := range inDegree {
		if v != 0 {
			continue
		}

		next = append(next, u)
	}

	for len(next) > 0 {
		u := next[0]
		next = next[1:]

		linearOrder = append(linearOrder, u)

		for _, v := range g[u] {
			inDegree[v]--

			if inDegree[v] == 0 {
				next = append(next, v)
			}
		}
	}
	return linearOrder
}

func canDone(n int, arr [][]int) bool {
	ind := make([]int, n)
	que := make([]int, 0)
	gra := make(map[int][]int)

	for _, v := range arr {
		ind[v[1]]++
		gra[v[0]] = append(gra[v[0]], v[1])
	}

	for i, v := range ind {
		if v == 0 {
			que = append(que, i)
		}
	}

	ret := make([]int, 0)
	for len(que) > 0 {
		elem := que[0]
		que = que[1:]

		ret = append(ret, elem)
		for _, neig := range gra[elem] {
			ind[neig]--
			if ind[neig] == 0 {
				que = append(que, neig)
			}
		}
	}

	return len(ret) == n
}

type Node struct {
	Val  int
	Next *Node
}

func midList(head *Node) *Node {
	if head == nil || head.Next == nil {
		return head
	}
	fast := head
	slow := head

	for fast.Next != nil && fast.Next.Next != nil {
		fast = fast.Next.Next
		slow = slow.Next
	}

	tmp := slow.Next
	slow.Next = nil
	slow = tmp
	return slow
}

func mergeTwoLists(l1 *Node, l2 *Node) *Node {
	if l1 == nil {
		return l2
	}
	if l2 == nil {
		return l1
	}
	if l1.Val < l2.Val {
		l1.Next = mergeTwoLists(l1.Next, l2)
		return l1
	}
	l2.Next = mergeTwoLists(l1, l2.Next)
	return l2
}

func ListM(left, right *Node) *Node {
	dummy := &Node{}
	cur := dummy
	for left != nil && right != nil {
		if left.Val < right.Val {
			cur.Next = left
			left = left.Next
		} else {
			cur.Next = right
			right = right.Next
		}
		cur = cur.Next
	}

	if left != nil {
		cur.Next = left
	}
	if right != nil {
		cur.Next = right
	}

	return dummy.Next
}

func ListMergeS(head *Node) *Node {
	if head == nil || head.Next == nil {
		return head
	}

	mid := midList(head)
	left := ListMergeS(head)
	right := ListMergeS(mid)

	return ListM(left, right)
}

func MinTotalV2(trian [][]int) int {
	dy := len(trian)
	dx := len(trian[dy-1])
	dp := make([][]int, dy+1)
	for i := range dp {
		dp[i] = make([]int, dx+1)
	}

	for i := dy - 1; i >= 0; i-- {
		for j := 0; j <= i; j++ {
			dp[i][j] = MinVal(dp[i+1][j], dp[i+1][j+1]) + trian[i][j]
		}
	}

	return dp[0][0]
}

func parti(arr []int, l, r int) int {
	pivot := l
	left := l - 1
	right := r + 1

	for {
		for {
			left++
			if arr[left] >= arr[pivot] {
				break
			}
		}

		for {
			right--
			if arr[right] <= arr[pivot] {
				break
			}
		}

		if left >= right {
			return right
		}

		arr[left], arr[right] = arr[right], arr[left]
	}
}

func quicS(arr []int, left, right int) {
	if left >= right {
		return
	}

	pivot := parti(arr, left, right)
	quicS(arr, left, pivot)
	quicS(arr, pivot+1, right)
	return
}

func bsLower(arr []int, target int) int {
	left := -1
	right := len(arr) - 1
	for left < right {
		mid := left + (right-left+1)>>1
		if arr[mid] <= target {
			left = mid
		} else {
			right = mid - 1
		}
	}

	//if arr[left] != target {
	//	return -1
	//} else {
	//	return left
	//}

	return left
}

func bsUpper(arr []int, target int) int {
	left := 0
	right := len(arr)
	for left < right {
		mid := left + (right-left)>>2
		if arr[mid] >= target {
			right = mid
		} else {
			left = mid + 1
		}
	}

	return right
}
