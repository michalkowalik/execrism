package hamming

import "errors"

const testVersion = 5

// Distance return Hamming distance between 2 DNA strands
func Distance(a, b string) (int, error) {
	if len(a) != len(b) {
		return -1, errors.New("a and b should be of equal length")
	}

	hammingDistance := 0
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			hammingDistance++
		}
	}

	return hammingDistance, nil
}
