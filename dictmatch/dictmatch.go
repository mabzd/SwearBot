package dictmatch

import (
	"fmt"
)

// Errors
const (
	Success = 0
	WordOverlappedByWildcardErr = 1
	WordExistErr = 2
	WildcardOverlappedByWordErr = 3
	WildcardRootExistErr = 4
)

// Node types
const (
	emptyNode = 0
	endNode = 1
	wildcardNode = 2
)

type Dict struct {
	tree *node
}

type DictErr struct {
	Desc string
	ErrType int
}

type node struct {
	runeMap map[rune]*node
	nodeType int
}

func NewDict() *Dict {
	return &Dict {
		tree: newNode(),
	}
}

func (dict *Dict) AddEntry(word string) *DictErr {
	errType := addRune(dict.tree, []rune(word))
	if errType != Success {
		return &DictErr {
			Desc: fmt.Sprintf("Error when adding '%s': %s", word, newDictErrDesc(errType)),
			ErrType: errType,
		}
	}
	return nil
}

func (dict *Dict) IsMatch(word string) bool {
	return matchRune(dict.tree, []rune(word))
}

func addRune(current *node, runes []rune) int {
	if current.nodeType == wildcardNode {
		return WordOverlappedByWildcardErr
	}
	if len(runes) == 0 {
		if current.nodeType == endNode {
			return WordExistErr
		}
		current.nodeType = endNode
		return Success
	}
	currentRune := runes[0]
	if currentRune == '*' {
		if current.runeMap != nil {
			return WildcardOverlappedByWordErr
		}
		if current.nodeType == endNode {
			return WildcardRootExistErr
		}
		current.nodeType = wildcardNode;
		return Success
	}
	if current.runeMap == nil {
		current.runeMap = make(map[rune]*node)
	}
	next := current.runeMap[currentRune]
	if next == nil {
		next = newNode()
		current.runeMap[currentRune] = next
	}
	return addRune(next, runes[1:])
}

func matchRune(current *node, runes []rune) bool {
	if current.nodeType == wildcardNode {
		return true
	}
	if len(runes) == 0 {
		return current.nodeType == endNode
	}
	currentRune := runes[0]
	next := current.runeMap[currentRune]
	if next == nil {
		return false
	}
	return matchRune(next, runes[1:])
}

func newNode() *node {
	return &node {
		runeMap: nil,
		nodeType: emptyNode,
	}
}

func newDictErrDesc(errType int) string {
	switch (errType) {
	case WordOverlappedByWildcardErr:
		return "Word is overlapped by existing wildcard entry."
	case WordExistErr:
		return "This word already exist."
	case WildcardOverlappedByWordErr:
		return "This wildcard entry is overlapped by existing word."
	case WildcardRootExistErr:
		return "This wildcard entry's root already exist."
	default:
		panic(fmt.Sprintf("Unknown error type: %d", errType))
	}
}