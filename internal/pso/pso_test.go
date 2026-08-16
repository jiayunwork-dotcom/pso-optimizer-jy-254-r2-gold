package pso

import (
	"math"
	"testing"
)

func TestOptimizeSphere(t *testing.T) {
	prob := &Problem{
		Dims:   5,
		Bounds: []Bounds{{Min: -5, Max: 5}, {Min: -5, Max: 5}, {Min: -5, Max: 5}, {Min: -5, Max: 5}, {Min: -5, Max: 5}},
		Objectives: []Objective{func(v Vector) float64 {
			s := 0.0
			for _, x := range v {
				s += x * x
			}
			return s
		}},
	}
	opt := New()
	opt.SwarmSize = 40
	opt.MaxIter = 200
	opt.Seed = 42
	res := opt.Run(prob)
	if res.Best.BestFit > 1e-3 {
		t.Errorf("sphere best fit too high: %g", res.Best.BestFit)
	}
}

func TestConstrained(t *testing.T) {
	// minimize x^2 + y^2 subject to x + y >= 1
	prob := &Problem{
		Dims:   2,
		Bounds: []Bounds{{Min: -2, Max: 2}, {Min: -2, Max: 2}},
		Objectives: []Objective{func(v Vector) float64 {
			return v[0]*v[0] + v[1]*v[1]
		}},
		Constraints: []Constraint{func(v Vector) float64 {
			viol := 1 - (v[0] + v[1]) // >0 when x+y < 1
			if viol > 0 {
				return viol
			}
			return 0
		}},
	}
	opt := New()
	opt.SwarmSize = 40
	opt.MaxIter = 200
	opt.Seed = 7
	res := opt.Run(prob)
	x, y := res.Best.Best[0], res.Best.Best[1]
	if x+y < 0.999 {
		t.Errorf("constraint violated: x+y = %g", x+y)
	}
	// optimal is (0.5,0.5) with value 0.5
	if math.Abs(x-0.5) > 0.1 || math.Abs(y-0.5) > 0.1 {
		t.Errorf("expected ~ (0.5,0.5), got (%g,%g)", x, y)
	}
}

func TestParetoFront(t *testing.T) {
	sols := []Vector{
		{1, 2}, // dominated by {1,1}
		{1, 1},
		{2, 0}, // non-dominated vs {1,1}
		{3, 3}, // dominated by both
	}
	front := paretoFront(sols)
	if len(front) != 2 {
		t.Fatalf("want 2 front members, got %d: %v", len(front), front)
	}
}

func TestResultBestIndependentOfPosition(t *testing.T) {
	opt := New()
	opt.SwarmSize = 8
	opt.MaxIter = 20
	opt.Seed = 3
	res := opt.Run(Sphere(2))
	if len(res.Best.Best) == 0 || len(res.Best.Position) == 0 {
		t.Fatal("expected non-empty Best and Position")
	}
	saved := res.Best.Position[0]
	res.Best.Best[0] = saved + 99
	if res.Best.Position[0] != saved {
		t.Errorf("writing Best[0] changed Position[0]: got %g want %g", res.Best.Position[0], saved)
	}
}

func TestFitnessDefaultEqualWeights(t *testing.T) {
	opt := New()
	prob := Sphere(2)
	fit, objs := opt.fitness(prob, Vector{3, 4})
	if len(objs) != 1 || objs[0] != 25 {
		t.Fatalf("raw objective: got %v want [25]", objs)
	}
	if fit != 25 {
		t.Errorf("fitness with default weights: got %g want 25", fit)
	}
}

func TestFitnessAppliesPenalty(t *testing.T) {
	opt := New()
	prob := ConstrainedQuad()
	fit, _ := opt.fitness(prob, Vector{0, 0})
	if fit < opt.Penalty*0.5 {
		t.Errorf("infeasible point fitness too small: %g", fit)
	}
	ok, _ := opt.fitness(prob, Vector{0.5, 0.5})
	if ok > 1 {
		t.Errorf("feasible point fitness too large: %g", ok)
	}
}

func TestDominatesWhenTiedOnOneObjective(t *testing.T) {
	a := Vector{1, 1}
	b := Vector{1, 2}
	if !dominates(a, b) {
		t.Fatal("expected {1,1} to dominate {1,2}")
	}
	if dominates(b, a) {
		t.Fatal("expected {1,2} not to dominate {1,1}")
	}
	if dominates(a, a) {
		t.Fatal("a vector should not dominate itself")
	}
}

func TestRunRespectsUpperBound(t *testing.T) {
	prob := &Problem{
		Dims:   1,
		Bounds: []Bounds{{Min: 0, Max: 1}},
		Objectives: []Objective{func(v Vector) float64 {
			return -v[0]
		}},
	}
	opt := New()
	opt.SwarmSize = 12
	opt.MaxIter = 40
	opt.Seed = 11
	res := opt.Run(prob)
	x := res.Best.Best[0]
	if x > 1+1e-9 {
		t.Errorf("best position %g exceeds upper bound 1", x)
	}
	if x < 0-1e-9 {
		t.Errorf("best position %g below lower bound 0", x)
	}
}
