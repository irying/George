package main

type TrieNode struct {
	children map[interface{}] *TrieNode
	isEnd bool
	num int // 记录多少单词途径该节点
	count int // 记录单词的出现次数
}

type Trie struct {
	root *TrieNode
}

func newTrieNode() *TrieNode  {
	return &TrieNode{children: make(map[interface{}]*TrieNode),isEnd:false,num:0,count:0}
}

func NewTrie() *Trie {
	return &Trie{root:newTrieNode()}
}

func (trie *Trie) Insert(word interface{})  {
	node := trie.root
	length := len(word)
	for i := 0; i < length; i++  {
		_, ok := node.children[word[i]]
		if !ok {
			node.children[word[i]] = newTrieNode()
		}
		node = node.children[word[i]]
		node.num++
	}
	node.isEnd = true
	node.count++
}

func (trie *Trie) Search(word []interface{}) bool  {
	node := trie.root
	length := len(word)
	for i := 0; i < length; i++ {
		_, ok := node.children[word[i]]
		if !ok {
			return false
		}
		node = node.children[word[i]]
	}

	return node.isEnd
}

func (trie *Trie) StartsWith(prefix []interface{}) bool  {
	node := trie.root
	length := len(prefix)
	for i := 0; i < length; i++  {
		_, ok := node.children[prefix[i]]
		if !ok {
			return false
		}
		node = node.children[prefix[i]]
	}

	return  true
}

func main() {
}
