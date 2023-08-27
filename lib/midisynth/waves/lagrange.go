package waves

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Lagrange struct {
	points []Point[float64]
}

var _ Wave = (*Lagrange)(nil)

func ParseLagrangeFile(path string) (*Lagrange, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseLagrangeInput(f)
}

func ParseLagrangeInput(in io.Reader) (*Lagrange, error) {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	points := []Point[int]{}
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		x, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, err
		}
		y, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, err
		}
		points = append(points, Point[int]{x, y})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return NewLagrange(points)
}

func NewLagrange(points []Point[int]) (*Lagrange, error) {
	if len(points) < 2 {
		return nil, fmt.Errorf("Not enough points for Lagrange: %v", points)
	}

	x0 := points[0].x
	xlast := points[len(points)-1].x
	xlen := float64(xlast - x0)

	fpoints := make([]Point[float64], len(points))

	for i, p := range points {
		xf := float64(p.x-x0) / xlen
		fpoints[i] = Point[float64]{xf, float64(p.y)}
	}

	lg := &Lagrange{
		points: fpoints,
	}
	lg.normalizePoints()
	return lg, nil
}

func (l *Lagrange) normalizePoints() {
	dx := 1.0 / 2200
	maxy := 0.0
	for x := 0.0; x <= 1.0; x += dx {
		y := l.value(x)
		if y > maxy {
			maxy = y
		} else if -y > maxy {
			maxy = -y
		}
	}
	if maxy == 0.0 {
		return
	}
	for i, p := range l.points {
		l.points[i].y = p.y / maxy
	}
}

func (l *Lagrange) Value(tm float64, ctx *NoteCtx) float64 {
	t := (tm - ctx.Period*float64(int(tm/ctx.Period))) * ctx.Freq
	return l.value(t)
}

type Point[T any] struct {
	x, y T
}

func (l *Lagrange) ljx(j int, x float64) float64 {
	res := 1.0
	for m, p := range l.points {
		if m == j {
			continue
		}
		res *= (x - p.x) / (l.points[j].x - p.x)
	}
	return res
}

func (l *Lagrange) value(x float64) float64 {
	res := 0.0
	for j, p := range l.points {
		res += p.y * l.ljx(j, x)
	}
	return res
}
