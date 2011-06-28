package polecalc

import "math"

type Vector2 struct {
	X, Y float64
}

func (v Vector2) Add(u Vector2) Vector2 {
	return Vector2{v.X + u.X, v.Y + u.Y}
}

func (v Vector2) Sub(u Vector2) Vector2 {
	return Vector2{v.X - u.X, v.Y - u.Y}
}

func (v Vector2) Mult(s float64) Vector2 {
	return Vector2{s * v.X, s * v.Y}
}

func (v Vector2) Dot(u Vector2) float64 {
	return v.X * u.X + v.Y * u.Y
}

func (v Vector2) Norm() float64 {
	return math.Sqrt(v.Dot(v))
}

func (v Vector2) NormSquared() float64 {
	return v.Dot(v)
}
