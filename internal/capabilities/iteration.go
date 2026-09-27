package capabilities

import "fmt"

// NextIteration returns the next canonical iteration token after current from the
// declared token sequence. It fails closed on an unknown or final token; the
// budget never wraps or resets (§10.6).
func NextIteration(tokens []string, current string) (string, error) {
	if len(tokens) == 0 {
		return "", fmt.Errorf("iteration: empty token sequence")
	}
	for i, token := range tokens {
		if token != current {
			continue
		}
		if i == len(tokens)-1 {
			return "", fmt.Errorf("iteration: budget exhausted at %s", current)
		}
		return tokens[i+1], nil
	}
	return "", fmt.Errorf("iteration: unknown token %q", current)
}
