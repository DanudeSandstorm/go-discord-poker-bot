package util

// Combinations generates all possible combinations of size k from the input slice
func Combinations[T any](input []T, k int) <-chan []T {
	ch := make(chan []T)
	go func() {
		defer close(ch)
		if k > len(input) {
			return
		}

		indices := make([]int, k)
		for i := range indices {
			indices[i] = i
		}

		combination := make([]T, k)
		for i, idx := range indices {
			combination[i] = input[idx]
		}
		ch <- combination

		for {
			// Find the rightmost index that can be incremented
			i := k - 1
			for i >= 0 && indices[i] == len(input)-k+i {
				i--
			}
			if i < 0 {
				break
			}

			// Increment the found index
			indices[i]++
			for j := i + 1; j < k; j++ {
				indices[j] = indices[j-1] + 1
			}

			// Generate next combination
			combination := make([]T, k)
			for i, idx := range indices {
				combination[i] = input[idx]
			}
			ch <- combination
		}
	}()
	return ch
}
