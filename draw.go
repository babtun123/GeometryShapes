// Samuel Shodiya

package main

import (
	//   "math"
	"errors"
	"fmt"
	"os"
)

type Color int

var (
	red    Color = 0
	green  Color = 1
	blue   Color = 2
	yellow Color = 3
	orange Color = 4
	purple Color = 5
	brown  Color = 6
	black  Color = 7
	white  Color = 8
)

// Initialize cmaps
var cmap = map[Color][3]int{
	red:    {255, 0, 0},
	green:  {0, 255, 0},
	blue:   {0, 0, 255},
	yellow: {255, 255, 0},
	orange: {255, 165, 0},
	purple: {128, 0, 128},
	brown:  {165, 42, 42},
	black:  {0, 0, 0},
	white:  {255, 255, 255},
}

// Initialize point struct with its x and y axis
type Point struct {
	x, y int
}

// geometry interface as given
type geometry interface {
	draw(scn screen) (err error)
	shape() (s string)
}

// Declare all struct for all shapes
type Triangle struct {
	pt1, pt2, pt3 Point
	c             Color
}
type Rectangle struct {
	ll, ur Point
	c      Color
}
type Circle struct {
	cp Point
	r  int
	c  Color
}

func (t Triangle) shape() (s string) {
	s = "Triangle"
	return
}

// https://gabrielgambetta.com/computer-graphics-from-scratch/07-filled-triangles.html
func interpolate(l0, d0, l1, d1 int) (values []int) {
	a := float64(d1-d0) / float64(l1-l0)
	d := float64(d0)

	count := l1 - l0 + 1
	for ; count > 0; count-- {
		values = append(values, int(d))
		d = d + a
	}
	return
}

// https://gabrielgambetta.com/computer-graphics-from-scratch/07-filled-triangles.html
func (tri Triangle) draw(scn screen) (err error) {
	if outOfBounds(tri.pt1, scn) || outOfBounds(tri.pt2, scn) || outOfBounds(tri.pt3, scn) {
		return errOutOfBounds
	}
	if colorUnknown(tri.c) {
		return errColorUnknown
	}

	y0 := tri.pt1.y
	y1 := tri.pt2.y
	y2 := tri.pt3.y

	// Sort the points so that y0 <= y1 <= y2
	if y1 < y0 {
		tri.pt2, tri.pt1 = tri.pt1, tri.pt2
	}
	if y2 < y0 {
		tri.pt3, tri.pt1 = tri.pt1, tri.pt3
	}
	if y2 < y1 {
		tri.pt3, tri.pt2 = tri.pt2, tri.pt3
	}

	x0, y0, x1, y1, x2, y2 := tri.pt1.x, tri.pt1.y, tri.pt2.x, tri.pt2.y, tri.pt3.x, tri.pt3.y

	x01 := interpolate(y0, x0, y1, x1)
	x12 := interpolate(y1, x1, y2, x2)
	x02 := interpolate(y0, x0, y2, x2)

	// Concatenate the short sides

	x012 := append(x01[:len(x01)-1], x12...)

	// Determine which is left and which is right
	var x_left, x_right []int
	m := len(x012) / 2
	if x02[m] < x012[m] {
		x_left = x02
		x_right = x012
	} else {
		x_left = x012
		x_right = x02
	}

	// Draw the horizontal segments
	for y := y0; y <= y2; y++ {
		for x := x_left[y-y0]; x <= x_right[y-y0]; x++ {
			scn.drawPixel(x, y, tri.c)
		}
	}
	return
}

func (c Circle) shape() (s string) {
	s = "Circle"
	return
}

// Circle draw function
// http://fredericgoset.ovh/mathematiques/courbes/en/filled_circle.html
/*
For My circle function, I got it from the link above.
The code is originally written in c/c++
The algorithm used is: The Bresenham's line algorithm
*/
func (c Circle) draw(scn screen) (err error) {
	if outOfBounds(c.cp, scn) {
		return errOutOfBounds
	}
	if colorUnknown(c.c) {
		return errColorUnknown
	}

	x := 0
	y := c.r
	m := 5 - 4*c.r

	for x <= y {
		fillSymmetricPoints(scn, c.cp.x, c.cp.y, x, y, c.c)
		if m > 0 {
			y--
			m -= 8 * y
		}
		x++
		m += 8*x + 4
	}
	return nil
}

// http://fredericgoset.ovh/mathematiques/courbes/en/filled_circle.html
// Helper function to Fill in Color to the circle
func fillSymmetricPoints(scn screen, fromMiddlePointX, fromMiddlePointY, x, y int, c Color) {
	for i := fromMiddlePointX - x; i <= fromMiddlePointX+x; i++ {
		scn.drawPixel(i, fromMiddlePointY-y, c)
		scn.drawPixel(i, fromMiddlePointY+y, c)
	}
	for i := fromMiddlePointX - y; i <= fromMiddlePointX+y; i++ {
		scn.drawPixel(i, fromMiddlePointY-x, c)
		scn.drawPixel(i, fromMiddlePointY+x, c)
	}
}

func (r Rectangle) shape() (s string) {
	s = "Rectangle"
	return
}

func (r Rectangle) draw(scn screen) (err error) {
	if outOfBounds(r.ll, scn) || outOfBounds(r.ur, scn) {
		return errOutOfBounds
	}
	if colorUnknown(r.c) {
		return errColorUnknown
	}
	// Calculate the coordinates of the rectangle corners
	var xMin, yMin, xMax, yMax int
	// compare xMin and xMax
	if r.ll.x < r.ur.x {
		xMin = r.ll.x
		xMax = r.ur.x
	} else {
		xMin = r.ur.x
		xMax = r.ll.x
	}
	// compare yMin and yMax
	if r.ll.y < r.ur.y {
		yMin = r.ll.y
		yMax = r.ur.y
	} else {
		yMin = r.ur.y
		yMax = r.ll.y
	}
	// Fill rectangle with Color
	for y := yMin; y <= yMax; y++ {
		for x := xMin; x <= xMax; x++ {
			scn.drawPixel(x, y, r.c)
		}
	}
	return
}

type screen interface {
	initialize(maxX, maxY int)
	getMaxXY() (maxX, maxY int)
	drawPixel(x, y int, c Color) (err error)
	getPixel(x, y int) (c Color, err error)
	clearScreen()
	screenShot(f string) (err error)
}

// initialize dsiplay struct
type Display struct {
	maxX, maxY int
	matrix     [][]Color
}

// initialize function for Display
func (dis *Display) initialize(maxX, maxY int) {
	// set the display max X and Y to para vals
	dis.maxX = maxX
	dis.maxY = maxY
	dis.matrix = make([][]Color, dis.maxY)
	for i := 0; i < dis.maxY; i++ {
		// initialize new array that will inside the matrix
		dis.matrix[i] = make([]Color, dis.maxX)
	}
	// use clearscreen function to initialize the screen to white
	dis.clearScreen()
}

// get the max X and Y
func (dis *Display) getMaxXY() (maxX, maxY int) {
	return dis.maxX, dis.maxY
}

// drawPixel Function
func (dis *Display) drawPixel(x, y int, c Color) (err error) {
	dis.matrix[x][y] = c
	return nil
}

// getPixel Function
func (dis *Display) getPixel(x, y int) (c Color, err error) {
	return dis.matrix[x][y], nil
}

// clearScreen Function
func (dis *Display) clearScreen() {
	// Set the colors of each pixel to white
	for i := 0; i < dis.maxY; i++ {
		for j := 0; j < dis.maxX; j++ {
			dis.matrix[i][j] = white
		}
	}
}

// screenshot function
func (dis *Display) screenShot(f string) error {
	outputFile, err := os.Create(f + ".ppm")
	if err != nil {
		return err
	}
	defer outputFile.Close()
	fmt.Fprintf(outputFile, "P3\n%d %d\n255\n", dis.maxX, dis.maxY)
	// Loop through the rows and columns for the image pixel matrix and set whatever Color value
	// to the cmap[c] and the array index
	for _, row := range dis.matrix {
		for _, c := range row {
			fmt.Fprintf(outputFile, "%d %d %d ", cmap[c][0], cmap[c][1], cmap[c][2])
		}
		fmt.Fprintln(outputFile)
	}
	return nil
}

var errOutOfBounds = errors.New("geometry out of bounds")
var errColorUnknown = errors.New("color unknown")

func outOfBounds(point Point, s screen) bool {
	maxX, maxY := s.getMaxXY()
	if point.x < 0 || point.x >= maxX || point.y < 0 || point.y >= maxY {
		return true
	} else {
		return false
	}
}

func colorUnknown(c Color) bool {
	if c < 0 || c > 8 {
		return true
	} else {
		return false
	}
}

// display
// TODO: you must implement the struct for this variable, and the interface it implements (screen)
var display Display
