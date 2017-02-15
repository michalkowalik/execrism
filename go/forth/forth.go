// Package forth provides minimal Forth evaluator
package forth

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const testVersion = 1

// before we start with Forth, we need some kind of stack implemented:
// element: we need 2 different stacks: one for values, second for words.

type stack struct {
	item []interface{}
}

// get new stack:
func newStack() *stack {
	return &stack{make([]interface{}, 0)}
}

// push new item to stack
func (s *stack) push(i interface{}) {
	s.item = append(s.item, i)
}

// pop an item from stack. return error if trying to pop from empty stack
func (s *stack) pop() (interface{}, error) {
	l := len(s.item)

	if l == 0 {
		return nil, errors.New("Empty Stack")
	}

	res := s.item[l-1]
	s.item = s.item[:l-1]
	return res, nil
}

// get back the stack items as integers:
func (s *stack) getInts() []int {
	ints := make([]int, len(s.item))

	for i := range s.item {
		ints[i] = s.item[i].(int)
	}

	return ints
}

var forthWords = []string{"+", "-", "*", "/", "DUP", "DROP", "SWAP", "OVER", ":"}

// Forth is the main evaluator function
func Forth(val []string) ([]int, error) {
	items := make([]string, 0)
	for _, st := range val {
		fmt.Printf("Statement: %s | ", st)
		items = append(items, itemize(st)...)

	}
	fmt.Printf("Items: %v \n", items)
	return parse(items)
}

// parse does the heavylifting of the statement evaluation:
func parse(items []string) ([]int, error) {
	wordStack := newStack()
	valueStack := newStack()

	for _, item := range items {

		i, err := strconv.Atoi(item)

		// parsed without problems -> integer value
		if err == nil {
			valueStack.push(i)
		} else {
			if isWord(item) {
				wordStack.push(item)
			} else {
				return make([]int, 0), errors.New("Wrong syntax?")
			}
		}
	}

	// show stacks:
	fmt.Printf("wordStack: %v \n", wordStack)
	fmt.Printf("valueStack: %v \n", valueStack)

	return valueStack.getInts(), nil
}

// divide forth statement to item:
// use regular expression to clean multiple white characters
func itemize(st string) []string {
	re1 := regexp.MustCompile("[\t\n\v\f\r\x00\x13]")
	re2 := regexp.MustCompile(" {2,}")
	st = re1.ReplaceAllString(st, " ")
	st = re2.ReplaceAllString(st, " ")
	return strings.Split(st, " ")
}

// isWord check if an item is a defined Forth word
// it feels weird I need to go through complete slice by hand..
func isWord(word string) bool {
	for _, i := range forthWords {
		if i == word {
			return true
		}
	}
	return false
}
