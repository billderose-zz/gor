package matrix

import (
	"fmt"
	"math"
)

const epsilon = 1e-12

type row []float64

type Matrix struct {
	Nrow, Ncol int
	Elem       []row
}

type Pivot struct {
	Row, Col int
}

// New returns a 0-matrix with the size initialized
func New(nrow, ncol int) *Matrix {
	m := Matrix{Nrow: nrow, Ncol: ncol, Elem: make([]row, nrow)}
	for i := 0; i < nrow; i++ {
		m.Elem[i] = make(row, ncol)
	}
	return &m
}

// NewPivot returns reference to a Pivot on row, col
func NewPivot(row, col int) *Pivot {
	return &Pivot{Row: row, Col: col}
}

// Rank returns the rank of the matrix m
func (m *Matrix) Rank() int {
	return int(math.Min(float64(m.Nrow), float64(m.Ncol)))
}

// NewFromVector returns a pointer to a matrix populated with entries
// of the provided vector such that m[i][j] = Elem[i+j]
func NewFromVector(Elem []float64, nrow int) (m *Matrix) {
	if nrow < 0 {
		return nil
	} else if nrow == 0 {
		return New(0, 0)
	}
	ncol := len(Elem) / nrow
	m = New(nrow, ncol)
	for i := 0; i < nrow; i++ {
		for j := 0; j < ncol; j++ {
			m.Elem[i][j] = Elem[i*ncol+j]
		}
	}
	return m
}

// SwapRow swaps row i with row j
func (m *Matrix) SwapRow(i, j int) *Matrix {
	m.Elem[i], m.Elem[j] = m.Elem[j], m.Elem[i]
	return m
}

// SwapCol swaps colulmn i with column j
func (m *Matrix) SwapCol(i, j int) *Matrix {
	return m.Transpose().SwapRow(i, j).Transpose()
}

// Row returns the ith row from m
func (m *Matrix) Row(i int) *Matrix {
	return NewFromVector(m.Elem[i], m.Nrow)
}

// Col returns the ith column from m
func (m *Matrix) Col(i int) *Matrix {
	col := make([]float64, 0, m.Ncol)
	for _, row := range m.Elem {
		col = append(col, row[i])
	}
	return NewFromVector(col, m.Nrow)
}

// Transpose returns the transpose of the matrix it is called on
func (m *Matrix) Transpose() *Matrix {
	r := make([]row, len(m.Elem[0])) // first row of transpose, first column of m
	for x := range r {
		r[x] = make(row, len(m.Elem)) // add columns that are the length m's rows
	}

	for y, s := range m.Elem {
		for x, e := range s {
			r[x][y] = e
		}
	}
	return &Matrix{Nrow: m.Ncol, Ncol: m.Nrow, Elem: r}
}

// FindPivot finds the next element to pivot on in m
func (m *Matrix) FindPivot() Pivot {
	rank := m.Rank()
	for j := 0; j < rank; j++ {
		for i, row := range m.Elem {
			if math.Abs(row[j]) > epsilon {
				if math.Abs(row[j]-1) > epsilon {
					pivots := 0
					for k := j - 1; k >= 0; k-- { // scan backwards for pivot in row
						if math.Abs(row[k]-1) < epsilon {
							pivots++
						}
					}
					if pivots == 1 { // already a pivoted on this row
						continue
					}
					return Pivot{Row: i, Col: j}
				} else {
					if i+1 == m.Nrow && j+1 == m.Ncol { // corner element; done pivoting
						continue
					}
					zeroed := true
					for k := i + 1; k < m.Nrow; k++ {
						if math.Abs(m.Elem[k][j]) > epsilon {
							zeroed = false
							break
						}
					}
					if !zeroed || i+1 == m.Nrow { // rest of column not zeroed or is last row
						return Pivot{Row: i, Col: j}
					}
				}
			}
		}
	}
	return Pivot{-1, -1}
}

// Pivot pivots on the Element at index i, j
func (m *Matrix) Pivot(p Pivot) *Matrix {
	pivot := m.Elem[p.Row][p.Col]
	m.ScaleRow(p.Row, 1/pivot) // make pivot entry 1
	for i, _ := range m.Elem { // for each row
		if p.Row != i { // skip pivot row
			if scaleFactor := m.Elem[i][p.Col]; scaleFactor != 0 {
				m.ScaleRow(p.Row, -scaleFactor)
				m.AddRow(p.Row, i)
				m.ScaleRow(p.Row, -1/scaleFactor)
			}
		}
	}
	return m
}

// ScaleRow scales row i by factor
func (m *Matrix) ScaleRow(i int, factor float64) {
	for n := range m.Elem[i] {
		m.Elem[i][n] = m.Elem[i][n] * factor
	}
}

//AddRow adds row i to row j, stores result in row j
func (m *Matrix) AddRow(i, j int) {
	for n := range m.Elem[i] { // for each Element in row
		m.Elem[j][n] += m.Elem[i][n]
	}
}

// RowReduce uses elementary row operations to put m into reduced row form
func (m *Matrix) RowReduce() {
	rank := m.Rank()
	for i := 0; i < rank; i++ {
		if p := m.FindPivot(); p.Col != -1 {
			m.Pivot(p)
		} else {
			break
		}
	}
	row := 0
	for i := range m.Elem { // clean up (you really don't want this for large matrices!)
		for j := range m.Elem[i] {
			if m.Elem[i][j] == 1 {
				m.SwapRow(row, i)
				row++
			}
			if math.Abs(m.Elem[i][j]) < epsilon {
				m.Elem[i][j] = 0.0
			}
		}
	}
}

// String returns the string representation of a Matrix
func (m Matrix) String() (s string) {
	s += fmt.Sprintf("%d-by-%d matrix\n", m.Nrow, m.Ncol)
	for i, row := range m.Elem {
		if i == 0 {
			s += fmt.Sprint(row)
		} else {
			s += "\n" + fmt.Sprint(row)
		}
	}
	return s
}
