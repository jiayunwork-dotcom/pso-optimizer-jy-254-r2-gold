// Package pso implements a Particle Swarm Optimizer for continuous
// minimization problems. It supports multiple objectives (tracked as a
// Pareto front) and inequality constraints (handled via penalty).
package pso

import (
	"math/rand"
)

// Vector is a point in the search space.
type Vector []float64

// Objective maps a position to a scalar to be minimized.
type Objective func(Vector) float64

// Constraint returns a violation amount (>= 0); 0 means satisfied.
type Constraint func(Vector) float64

// Bounds constrain one dimension.
type Bounds struct{ Min, Max float64 }

// Problem describes the optimization task.
type Problem struct {
	Dims       int
	Bounds     []Bounds
	Objectives []Objective
	Constraints []Constraint
	Weights    []float64 // for scalar fitness; defaults to equal if nil
}

// Particle is one swarm member.
type Particle struct {
	Position Vector
	Velocity Vector
	Best     Vector
	BestFit  float64
	Objs     Vector // objective values at Best
}

// Result holds the best particle and the discovered Pareto front (objective
// vectors that are not dominated by any other found solution).
type Result struct {
	Best  Particle
	Front []Vector
}

// Optimizer holds PSO hyperparameters.
type Optimizer struct {
	SwarmSize int
	MaxIter   int
	W         float64 // inertia
	C1        float64 // cognitive
	C2        float64 // social
	Penalty   float64
	Seed      int64
}

// New returns an Optimizer with sensible defaults.
func New() *Optimizer {
	return &Optimizer{SwarmSize: 30, MaxIter: 100, W: 0.72, C1: 1.5, C2: 1.5, Penalty: 1e6, Seed: 1}
}

func (p *Problem) weights() []float64 {
	if p.Weights != nil {
		return p.Weights
	}
	w := make([]float64, len(p.Objectives))
	for i := range w {
		w[i] = 1
	}
	return w
}

// fitness returns the scalar fitness (lower is better) and the raw objective
// values for a position.
func (o *Optimizer) fitness(prob *Problem, pos Vector) (float64, Vector) {
	w := prob.weights()
	objs := make(Vector, len(prob.Objectives))
	sum := 0.0
	for i, obj := range prob.Objectives {
		v := obj(pos)
		objs[i] = v
		sum += w[i] * v
	}
	penalty := 0.0
	for _, c := range prob.Constraints {
		if v := c(pos); v > 0 {
			penalty += v
		}
	}
	return sum + o.Penalty*penalty, objs
}

// Run executes the optimization and returns the best particle and Pareto front.
func (o *Optimizer) Run(prob *Problem) Result {
	rng := rand.New(rand.NewSource(o.Seed))
	particles := make([]Particle, o.SwarmSize)
	rangeOf := func(i int) float64 {
		b := prob.Bounds[i]
		return b.Max - b.Min
	}
	for i := range particles {
		pos := make(Vector, prob.Dims)
		vel := make(Vector, prob.Dims)
		for d := 0; d < prob.Dims; d++ {
			b := prob.Bounds[d]
			pos[d] = b.Min + rng.Float64()*rangeOf(d)
			vel[d] = (rng.Float64()*2 - 1) * rangeOf(d) * 0.1
		}
		fit, objs := o.fitness(prob, pos)
		particles[i] = Particle{
			Position: pos, Velocity: vel,
			Best: append(Vector(nil), pos...), BestFit: fit, Objs: objs,
		}
	}

	gbest := particles[0]
	for iter := 0; iter < o.MaxIter; iter++ {
		for i := range particles {
			p := &particles[i]
			for d := 0; d < prob.Dims; d++ {
				r1, r2 := rng.Float64(), rng.Float64()
				b := prob.Bounds[d]
				cog := o.C1 * r1 * (p.Best[d] - p.Position[d])
				soc := o.C2 * r2 * (gbest.Best[d] - p.Position[d])
				p.Velocity[d] = o.W*p.Velocity[d] + cog + soc
				// clamp velocity
				vmax := rangeOf(d) * 0.2
				if p.Velocity[d] > vmax {
					p.Velocity[d] = vmax
				}
				if p.Velocity[d] < -vmax {
					p.Velocity[d] = -vmax
				}
				p.Position[d] += p.Velocity[d]
				if p.Position[d] < b.Min {
					p.Position[d] = b.Min
				}
				if p.Position[d] > b.Max {
					p.Position[d] = b.Max
				}
			}
			fit, objs := o.fitness(prob, p.Position)
			if fit < p.BestFit {
				p.BestFit = fit
				p.Best = append(Vector(nil), p.Position...)
				p.Objs = objs
			}
			if p.BestFit < gbest.BestFit {
				gbest = *p
			}
		}
	}

	// Pareto front over all personal bests.
	objsList := make([]Vector, 0, len(particles))
	for _, p := range particles {
		objsList = append(objsList, p.Objs)
	}
	return Result{Best: gbest, Front: paretoFront(objsList)}
}

// dominates reports whether a is strictly better than or equal to b on all
// objectives and strictly better on at least one (minimization).
func dominates(a, b Vector) bool {
	better := false
	for i := range a {
		if a[i] > b[i] {
			return false
		}
		if a[i] < b[i] {
			better = true
		}
	}
	return better
}

// paretoFront returns the non-dominated objective vectors.
func paretoFront(sols []Vector) []Vector {
	var front []Vector
	for i, s := range sols {
		dominated := false
		for j, o := range sols {
			if i == j {
				continue
			}
			if dominates(o, s) {
				dominated = true
				break
			}
		}
		if !dominated {
			front = append(front, append(Vector(nil), s...))
		}
	}
	return front
}
