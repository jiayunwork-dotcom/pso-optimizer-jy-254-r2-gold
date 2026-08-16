package pso

import "math"

// Sphere is a classic single-objective benchmark: minimize sum(x_i^2).
func Sphere(dims int) *Problem {
	b := make([]Bounds, dims)
	for i := range b {
		b[i] = Bounds{Min: -5, Max: 5}
	}
	return &Problem{
		Dims:   dims,
		Bounds: b,
		Objectives: []Objective{func(v Vector) float64 {
			s := 0.0
			for _, x := range v {
				s += x * x
			}
			return s
		}},
	}
}

// ZDT1 is a two-objective benchmark on the unit hypercube.
func ZDT1(dims int) *Problem {
	b := make([]Bounds, dims)
	for i := range b {
		b[i] = Bounds{Min: 0, Max: 1}
	}
	return &Problem{
		Dims:   dims,
		Bounds: b,
		Objectives: []Objective{
			func(v Vector) float64 { return v[0] },
			func(v Vector) float64 {
				g := 1.0
				if dims > 1 {
					s := 0.0
					for _, x := range v[1:] {
						s += x
					}
					g = 1 + 9*s/float64(dims-1)
				}
				f1 := v[0]
				return g * (1 - math.Sqrt(f1/g))
			},
		},
	}
}

// ConstrainedQuad minimizes x^2 + y^2 subject to x + y >= 1.
func ConstrainedQuad() *Problem {
	return &Problem{
		Dims:   2,
		Bounds: []Bounds{{Min: -2, Max: 2}, {Min: -2, Max: 2}},
		Objectives: []Objective{func(v Vector) float64 {
			return v[0]*v[0] + v[1]*v[1]
		}},
		Constraints: []Constraint{func(v Vector) float64 {
			if viol := 1 - (v[0] + v[1]); viol > 0 {
				return viol
			}
			return 0
		}},
	}
}
