package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"sync"
)

var wg sync.WaitGroup

type Point struct {
	x float64
	y float64
}
type Triangle struct {
	A Point
	B Point
	C Point
}

func (t *Triangle) print() {
	fmt.Printf("Triangle: A( %f, %f), B(%f, %f), C(%f, %f)", t.A.x, t.A.y, t.B.x, t.B.y, t.C.x, t.C.y)
}

type Stack struct {
	items []Triangle
	sem   chan int
	size  int
}

func newStack() Stack {
	sem := make(chan int, 1)
	items := make([]Triangle, 0)
	return Stack{items, sem, 0}
}

func (stack *Stack) push(t Triangle) {
	stack.sem <- 1
	stack.items = append(stack.items, t)
	stack.size = stack.size + 1
	<-stack.sem

}

func (stack *Stack) peek() *Triangle {
	stack.sem <- 1
	temp := stack.items[stack.size-1]
	<-stack.sem
	return &temp

}

func triangles10000() (result [10000]Triangle) {
	rand.Seed(2120)
	for i := 0; i < 10000; i++ {
		result[i].A = Point{rand.Float64() * 100., rand.Float64() * 100.}
		result[i].B = Point{rand.Float64() * 100., rand.Float64() * 100.}
		result[i].C = Point{rand.Float64() * 100., rand.Float64() * 100.}
	}
	return
}

func (t Triangle) Perimeter() float64 {
	side1 := math.Sqrt(math.Pow(t.A.x-t.B.x, 2) + math.Pow(t.A.y-t.B.y, 2))
	side2 := math.Sqrt(math.Pow(t.B.x-t.C.x, 2) + math.Pow(t.B.y-t.C.y, 2))
	side3 := math.Sqrt(math.Pow(t.C.x-t.A.x, 2) + math.Pow(t.C.y-t.A.y, 2))
	return side1 + side2 + side3
}

func (t Triangle) Area() float64 {
	return 0.5 * ((t.B.x-t.A.x)*(t.C.y-t.A.y) - (t.C.x-t.A.x)*(t.B.y-t.A.y))
}

func classifyTriangles(highRatio *Stack, lowRatio *Stack,
	ratioThreshold float64, triangles []Triangle) {
	for i := range triangles {
		if triangles[i].Perimeter()/triangles[i].Area() > 1.0 {
			highRatio.push(triangles[i])
		} else {
			lowRatio.push(triangles[i])
		}
	}
	wg.Done()

}

func main() {
	triangles := triangles10000()
	high := newStack()
	low := newStack()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go classifyTriangles(&high, &low, 1, triangles[i:i+1000])
	}

	wg.Wait()
	fmt.Printf("\nHigh Ratio Stack:\n   Size: " + strconv.Itoa(high.size) + "\n   Top Element: ")
	high.peek().print()
	fmt.Printf("\nLow Ratio Stack:\n   Size: " + strconv.Itoa(low.size) + "\n   Top Element: ")
	low.peek().print()

}
