// Command pso-optimizer runs a particle swarm optimizer on a built-in
// benchmark or constrained problem and prints the best solution found.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"pso-optimizer/internal/pso"
)

func main() {
	problem := flag.String("problem", "sphere", "problem: sphere | zdt1 | constrained")
	dims := flag.Int("dims", 10, "number of dimensions (sphere/zdt1)")
	iters := flag.Int("iters", 200, "max iterations")
	swarm := flag.Int("swarm", 40, "swarm size")
	seed := flag.Int64("seed", 1, "RNG seed")
	flag.Parse()

	opt := pso.New()
	opt.MaxIter = *iters
	opt.SwarmSize = *swarm
	opt.Seed = *seed

	var prob *pso.Problem
	switch strings.ToLower(*problem) {
	case "sphere":
		prob = pso.Sphere(*dims)
	case "zdt1":
		prob = pso.ZDT1(*dims)
	case "constrained":
		prob = pso.ConstrainedQuad()
	default:
		fatal("unknown problem %q (sphere|zdt1|constrained)", *problem)
	}

	res := opt.Run(prob)
	fmt.Printf("problem:   %s (dims=%d)\n", *problem, prob.Dims)
	fmt.Printf("best fit:  %g\n", res.Best.BestFit)
	fmt.Printf("best pos:  %s\n", formatVec(res.Best.Best))
	if len(res.Front) > 0 {
		fmt.Printf("pareto front size: %d\n", len(res.Front))
		fmt.Printf("front (objective vectors):\n")
		for _, f := range res.Front {
			fmt.Printf("  %s\n", formatVec(f))
		}
	}
}

func formatVec(v pso.Vector) string {
	parts := make([]string, len(v))
	for i, x := range v {
		parts[i] = fmt.Sprintf("% .4f", x)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "pso-optimizer: "+format+"\n", args...)
	os.Exit(1)
}
