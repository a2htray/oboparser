package oboparser

func stringIn(token string, tokens []string) bool {
	in := false
	for _, v := range tokens {
		if v == token {
			in = true
		}
	}
	return in
}

