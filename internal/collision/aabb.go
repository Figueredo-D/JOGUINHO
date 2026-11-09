package collision

// Rect representa uma caixa delimitadora 2D (AABB).
type Rect struct {
	X, Y, W, H float64
}

// NewRect cria uma nova caixa delimitadora.
func NewRect(x, y, w, h float64) Rect {
	return Rect{X: x, Y: y, W: w, H: h}
}

// Intersects verifica se duas caixas delimitadoras colidem.
func (r Rect) Intersects(other Rect) bool {
	return r.X < other.X+other.W &&
		r.X+r.W > other.X &&
		r.Y < other.Y+other.H &&
		r.Y+r.H > other.Y
}

// Contains verifica se um ponto (px, py) está dentro da caixa.
func (r Rect) Contains(px, py float64) bool {
	return px >= r.X && px <= r.X+r.W &&
		py >= r.Y && py <= r.Y+r.H
}
