// Package forth provides minimal Forth evaluator
package forth

/*
Warning: this is a very first version of the solution.
The code will contain tons of completely unnecessary comments
*/

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
// Note: as it seems, a stack for ints should be enough.
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
	// dictionary of user defined words:
	// userWords := make(map[string]string)
	valueStack := newStack()

	// val will contain one or more Forth statements.
	// Each of them needs to be parsed and evaluated separately
	// The only thing that joins them is the common stack and user defined words dictionary
	for _, st := range val {
		items := itemize(st)
		fmt.Printf("Items: %v \n", items)
		if err := parse(items, valueStack); err != nil {
			return nil, err
		}
	}
	return valueStack.getInts(), nil
}

// parse does the heavylifting of the statement evaluation.
// Evaluation result stays in the modified valueStack.
// return value is used to check whether error occured
func parse(items []string, valueStack *stack) error {

	// potentially, the range solution will need to be replaced
	// with a plain loop with indices.
	// the ":" word starts a definition of the new word: everything up to the ";"
	// goes as a def to the map. Then update the counter on the current position in the
	// slice.
	// Also -> check for words in the forthWords and in the map.
	// should it be a user def. word: -> create a new items string slice with the user word
	// replaced by it's definition and recursively call the parse function.
	for _, item := range items {
		i, err := strconv.Atoi(item)

		// parsed without problems -> integer value
		if err == nil {
			valueStack.push(i)
		} else {
			if isWord(item) {
				if err := eval(item, valueStack); err != nil {
					return err
				}
			} else {
				return errors.New("Wrong syntax?")
			}
		}
	}

	// show stacks:
	fmt.Printf("valueStack: %v \n", valueStack)

	return nil
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
		return err
	}
	return nil
}

// over copies second to last element of the stack on top of it
func over(s *stack) error {
	op1, err := s.pop()
	if err != nil {
		return err
	}
	op2, err := s.pop()
	if err != nil {
		return err
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
