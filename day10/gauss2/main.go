package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/alex-ant/gomath/gaussian-elimination"
	"github.com/alex-ant/gomath/rational"
)

func main() {
	nr := func(i int64) rational.Rational {
		return rational.New(i, 1)
	}

	text := `[0 0 1 0 1 0 1 29]
[0 0 0 0 1 1 0 33]
[0 0 0 1 1 0 0 22]
[1 1 1 0 0 0 0 34]
[1 0 0 0 1 1 0 45]
[1 0 0 0 1 0 1 35]`

	text = `[1 1 1 0 10]
[1 0 1 1 11]
[1 0 1 1 11]
[1 1 0 0 5]
[1 1 1 0 10]
[0 0 1 0 5]`

	replacer := strings.NewReplacer("[", "", "]", "")
	lines := strings.Split(text, "\n")
	equations := make([][]rational.Rational, len(lines))
	for l, line := range lines {
		line = replacer.Replace(line)
		tokens := strings.Fields(line)
		equations[l] = make([]rational.Rational, len(tokens))
		for i, t := range tokens {
			val, _ := strconv.ParseInt(t, 10, 64)
			equations[l][i] = nr(val)
		}
	}

	res, gausErr := gaussian.SolveGaussian(equations, false)
	if gausErr != nil {
		log.Fatal(gausErr)
	}

	for _, v := range res {
		log.Println(v)
	}
	// Output:
	// [{1 1}]
	// [{2 1}]
	// [{3 1}]
	// [{4 1}]
}
