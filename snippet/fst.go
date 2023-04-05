package snippet

// FST is a trie data structure that can be used to store a set of strings.
/*
FST（有限状态转换器）是一种用于高效存储和查询字符串数据的数据结构。它可以通过将所有字符串构建成一个有限状态自动机（DFA）来实现。

在FST中，每个节点都代表一个字符串前缀或后缀，并且每个节点都有一个出边列表，其中每个出边都标识了一个字符以及连接到的下一个节点。通过遍历这个有向图，可以在FST中查找任何给定的字符串。

这里我们定义了一个Node结构体，表示FST中的节点。每个节点都有一个isWord字段，用于表示该节点是否是一个字符串的结尾。每个节点还有一个children字段，用于存储所有从该节点出发的边和连接的下一个节点。

我们还定义了一个FST结构体，表示整个有限状态转换器。NewFST函数用于创建一个新的FST实例，并初始化根节点。Insert函数用于将一个字符串插入到FST中，它遍历字符串中的每个字符，并在FST中创建新的节点以表示该字符。最后一个节点被标记为字符串的结尾。Search函数用于在FST中搜索一个字符串，它遍历字符串中的每个字符，并在FST中查找与该字符对应的节点。如果找到了该节点，则继续遍历下一个字符，否则返回false。如果在字符串的末尾找到了一个节点，则返回true。
*/
type Node struct {
	isWord   bool
	children map[rune]*Node
}

type FST struct {
	root *Node
}

func NewFST() *FST {
	return &FST{root: &Node{children: make(map[rune]*Node)}}
}

func (f *FST) Insert(word string) {
	node := f.root
	for _, char := range word {
		if _, ok := node.children[char]; !ok {
			node.children[char] = &Node{children: make(map[rune]*Node)}
		}
		node = node.children[char]
	}
	node.isWord = true
}

func (f *FST) Search(word string) bool {
	node := f.root
	for _, char := range word {
		if _, ok := node.children[char]; !ok {
			return false
		}
		node = node.children[char]
	}
	return node.isWord
}
