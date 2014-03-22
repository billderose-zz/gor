package matrix

import (
	//"math/rand"
	//"fmt"
	"math"
	"testing"
)

func TestNew(t *testing.T) {
	nrow := 5
	ncol := 6
	m := New(nrow, ncol)
	if m.Nrow != nrow {
		t.Errorf("New matrix should have %d", nrow, "rows. Has %d", m.Nrow)
	}
	if m.Ncol != ncol {
		t.Errorf("New matrix should have %d", ncol, "rows. Has %d", m.Ncol)
	}
	for i := 0; i < nrow; i++ {
		for j := 0; j < ncol; j++ {
			if m.Elem[i][j] != 0 {
				t.Errorf("New matrix not iniitalized to all 0")
			}
		}
	}
}

func TestNewFromVector(t *testing.T) {
	vec := []float64{1, 2, 3, 4, 5, 6}
	nrow := 3
	ncol := 2
	m := NewFromVector(vec, nrow)
	for i := 0; i < nrow; i++ {
		for j := 0; j < ncol; j++ {
			if m.Elem[i][j] != vec[i*ncol+j] {
				t.Errorf("New matrix not iniitalized with correct entries."+
					"Entry %f,%f should be %f, but is %f.", i, j, vec[i*ncol+j], m.Elem[i][j])
			}
		}
	}

	m = NewFromVector(vec, -1)
	if m != nil {
		t.Errorf("NewFromVector expected to return nil matrix when provided with negative nrow argument")
	}

	m = NewFromVector(vec, 0)
	if m.Nrow != 0 {
		t.Errorf("NewFromVector expected to return 0-by-0 matrix when provided with nrow = 0. Instead, "+
			"produced matrix %d-by-%d", m.Nrow, m.Ncol)
	}
	if m.Ncol != 0 {
		t.Errorf("NewFromVector expected to return 0-by-0 matrix when provided with nrow = 0. Instead, "+
			"produced matrix %d-by-%d", m.Nrow, m.Ncol)
	}
}

func TestSwapRow(t *testing.T) {
	vec := []float64{1, 2, 3, 4, 5, 6}
	nrow := 3
	//ncol := 2
	m := NewFromVector(vec, nrow)
	m.SwapRow(1, 2)
	if m.Elem[1][0] != 5 || m.Elem[1][1] != 6 {
		t.Errorf("SwapRow incorrectly exchanging rows")
	}
	if m.Elem[2][0] != 3 || m.Elem[2][1] != 4 {
		t.Errorf("SwapRow incorrectly exchanging rows")
	}
}

func TestSwapCol(t *testing.T) {
	vec := []float64{1, 2, 3, 4, 5, 6}
	nrow := 3
	//ncol := 2
	m := NewFromVector(vec, nrow)
	m = m.SwapCol(0, 1) // this is an API inconsistency. Library should handle all memory
	if m.Elem[0][0] != 2 || m.Elem[1][0] != 4 || m.Elem[2][0] != 6 {
		t.Errorf("SwapCol incorrectly exchanging columns")
	}
	if m.Elem[0][1] != 1 || m.Elem[1][1] != 3 || m.Elem[2][1] != 5 {
		t.Errorf("SwapCol incorrectly exchanging columns")
	}
}

var pivotTests = []struct {
	vec  []float64
	nrow int
	pRow int
	pCol int
}{
	{[]float64{0, 0,
		1, -1,
		5, 6,
		7, 9}, 4, 1, 0},
	{[]float64{4.4, 5,
		6, 7.5}, 2, 0, 0},
	{[]float64{1, 9.0 / 10.0, 5.0 / 6.0,
		0, 5, 3}, 2, 1, 1},
	{[]float64{-1, 0, 4, 5, 6, 7}, 3, 0, 0},
	{[]float64{1, 2, 3, 4, 5, 6, 7, 8}, 8, 0, 0},
	{[]float64{-2, 0, 4, 5, 9.2, 0.3}, 2, 0, 0},
	{[]float64{0, 0, 0, 1, 9.2, 0}, 2, 1, 0},
	{[]float64{0, 0, 0, 0, 9.2, 0}, 2, 1, 1},
	{[]float64{1, 0, 0, 0, 1, 0, 0, 0, 1}, 3, -1, -1},
	{[]float64{3, 4, 5, 7, 9, 1, 0, 8, 5, 4, 3, -7}, 4, 0, 0},
	{[]float64{5.5, 6.7, 9, 0.1, 3, 4, 5.6, 8}, 2, 0, 0},
	{[]float64{0, 6.7, 9, 0.1, 3, 4, 5.6, 8, .3, .5, .6, .11, 3.0 / 4.0, 5.0 / 7.0, 8.0 / 9.0, 0}, 4, 1, 0},
}

var reducedTests = []struct {
	vec []float64
}{
	{[]float64{1, 0, 0, 1, 0, 0, 0, 0}},
	{[]float64{1, 0, 0, 1}},
	{[]float64{1, 0, 0.2933333, 0, 1, 0.6}},
	{[]float64{1, 0, 0, 1, 0, 0}},
	{[]float64{1, 0, 0, 0, 0, 0, 0, 0}},
	{[]float64{1, 0, -2, 0, 1, 1.11957}},
	{[]float64{1, 9.2, 0, 0, 0, 0}},
	{[]float64{0, 1, 0, 0, 0, 0}},
	{[]float64{1, 0, 0, 0, 1, 0, 0, 0, 1}},
	{[]float64{1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0}},
	{[]float64{1, 0, -0.8, -28, 0, 1, 2, 23}},
	{[]float64{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}},
}

func TestFindPivot(t *testing.T) {
	for i, test := range pivotTests {
		m := NewFromVector(test.vec, test.nrow)
		p := m.FindPivot()

		if p.Col != test.pCol {
			t.Errorf("Incorrectly finding  pivot column on test case %d. "+
				"Found %d, should be %d\n", i, p.Col, test.pCol)
		}
		if p.Row != test.pRow {
			t.Errorf("Incorrectly finding  pivot row on test case %d. "+
				"Found %d, should be %d\n", i, p.Row, test.pRow)
		}
	}
}

func TestRowReduce(t *testing.T) {
	for i, test := range pivotTests {
		m := NewFromVector(test.vec, test.nrow)
		m.RowReduce()
		check := NewFromVector(reducedTests[i].vec, test.nrow)
		if check.Ncol != m.Ncol || check.Nrow != m.Nrow {
			t.Error("RowReduce altering the number of rows/columns")
		}
		for j := range m.Elem {
			for k := range m.Elem[j] {
				if math.Abs(m.Elem[j][k]-check.Elem[j][k]) > 1e-5 {
					t.Error("Incorrectly row-reducing matrix,"+
						"should be:\n", check, "\nIs:\n", m)
				}
			}
		}

	}
}
