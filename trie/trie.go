package main

import "fmt"

type TrieNode struct {
    children map[rune]*TrieNode
    endOfWord bool
}

func NewTrieNode() *TrieNode {
    return &TrieNode{
        children: make(map[rune]*TrieNode),
        endOfWord: false,
    }
}

func (root *TrieNode) Insert(word string) {
    node := root
    for _, ch := range word {
        if _, ok := node.children[ch]; !ok {
            node.children[ch] = NewTrieNode()
        }
        node = node.children[ch]
    }
    node.endOfWord = true
}

func (root *TrieNode) Search(word string) bool {
    node := root
    for _, ch := range word {
        if _, ok := node.children[ch]; !ok {
            return false
        }
        node = node.children[ch]
    }
    return node != nil && node.endOfWord
}

func (root *TrieNode) getAllWords(node *TrieNode, currentWord string, words *[]string) {
    if node.endOfWord {
        *words = append(*words, currentWord)
    }
    for ch, childNode := range node.children {
        root.getAllWords(childNode, currentWord + string(ch), words)
    }
}

func main() {
    root := NewTrieNode()

    words := []string{"hello", "dog", "hell", "cat", "a", "hel", "help", "helps", "helping"}
    for _, word := range words {
        root.Insert(word)
    }

    fmt.Println(root.Search("help"))
    fmt.Println(root.Search("hel"))
    fmt.Println(root.Search("zebra"))

    var trieWords []string
    root.getAllWords(root, "", &trieWords)
    fmt.Println(trieWords)
}