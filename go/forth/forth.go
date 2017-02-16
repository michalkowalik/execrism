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
		items = append(items, itemize(st)...)

	}
	fmt.Printf("Items: %v \n", items)
	return parse(items)
}

// parse does the heavylifting of the statement evaluation:
func parse(items []string) ([]int, error) {
	// wordStack := newStack()
	valueStack := newStack()

	for _, item := range items {
		i, err := strconv.Atoi(item)

		// parsed without problems -> integer value
		if err == nil {
			valueStack.push(i)
		} else {
			if isWord(item) {
				if err := eval(item, valueStack); err != nil {
					return nil, errors.New("error while evaluating expression")
				}
			} else {
				return nil, errors.New("Wrong syntax?")
			}
		}
	}

	// show stacks:
	// fmt.Printf("wordStack: %v \n", wordStack)
	fmt.Printf("valueStack: %v \n", valueStack)

	return valueStack.getInts(), nil
}

// evaluate word expression:
func eval(word string, s *stack) error {
	switch strings.ToUpper(word) {
	case "+":
		return execute(s, func(a, b int) (int, error) { return a + b, nil })
	case "-":
		return execute(s, func(a, b int) (int, error) { return a - b, nil })
	case "*":
		return execute(s, func(a, b int) (int, error) { return a * b, nil })
	case "/":
		return execute(s, func(a, b int) (int, error) {
			if b == 0 {
				return 0, errors.New("Dividing by 0")
			}
			return a / b, nil
		})
	case "DUP":
		return dup(s)
	case "DROP":
		return drop(s)
	case "SWAP":
		return swap(s)
	case "OVER":
		return over(s)
	}
	return nil
}

// dup duplicates the top element of the stack:
func dup(s *stack) error {
	op, err := s.pop()
	if err != nil {
		return errors.New("empty stack")
	}
	s.push(op)
	s.push(op)
	return nil
}

// drop removes top element from the stack
func drop(s *stack) error {
	_, err := s.pop()
	if err != nil {
		return errors.New("dropping from empty stack")
	}
	return nil
}

// over copies second to last element of the stack on top of it
func over(s *stack) error {
	op1, err := s.pop()
	if err != nil {
		return errors.New("empty stack")
	}
	op2, err := s.pop()
	if err != nil {
		return errors.New("empty stack")
	}

	s.push(op2)
	s.push(op1)
	s.push(op2)
	return nil
}

// swap swaps top 2 elements of the stack
func swap(s *stack) error {
	op1, err := s.pop()
	if err != nil {
		return errors.New("empty stack")
	}
	op2, err := s.pop()
	if err != nil {
		return errors.New("empty stack")
	}

	s.push(op1)
	s.push(op2)
	return nil
}

// execute provides the arithmetic part of the interpeter
func execute(s *stack, op func(int, int) (int, error)) error {
	op1, err := s.pop()
	if err != nil {
		return errors.New("empty stack")
	}
	op2, err := s.pop()
	if err != nil {
		return errors.New("empty stack")
	}
	res, err := op(op2.(int), op1.(int))
	if err != nil {
		return err
	}
	s.push(res)
	return nil
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
		if i == strings.ToUpper(word) {
			return true
		}
	}
	return false
}
