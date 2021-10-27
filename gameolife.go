package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"time"

	"golang.org/x/term"
)

//////////////////////////////////////////////////////////////////////////////
//																			//
//	RULES																	//
//  cell with fewer than 2 live neighbors DIES								//
//  cell with more than 3 live neighbors DIES								//
//  any dead cell with exactly 3 live neighbors comes back to LIFE			//
//  only cells with 2 or 3 live neighbors SURVIVES							//
//																			//
//////////////////////////////////////////////////////////////////////////////

//gets size of terminal window in rows/columns
func terminalSize() [2]int {
	width, height, err := term.GetSize(0)
	if err != nil {
		return [2]int{0, 0}
	} else {
		return [2]int{width, height}
	}
}

//struct for generation of cells, a grid of true/false cells for alive/dead
type Generation struct {
	grid   [][]bool
	width  int
	height int
}

//makes a grid and returns a pointer to it
func makeGrid(width, height int) *Generation {
	grid := make([][]bool, height)
	for i := 0; i < height; i++ {
		grid[i] = make([]bool, width)
	}
	return &Generation{grid: grid, width: width, height: height}
}

//initial grid values should only be generated in a small range in the center of the grid, this
//takes terminal size and returns dimentions to define range of initial values in grid
func getInitialRange(termSize [2]int) [2]int {
	width := math.Round(float64(termSize[0]) * 0.25)
	height := math.Round(float64(termSize[1]) * 0.35)
	x := int(width)
	y := int(height)
	return [2]int{x, y}
}

//assign values to initial grid
func assignValues(initGrid *Generation) {

	//define x and y as range of values to randomly update for initial grid
	x := getInitialRange(terminalSize())[0]
	y := getInitialRange(terminalSize())[1]

	for i := y; i < (initGrid.height - y); i++ {
		for j := x; j < (initGrid.width - x); j++ {

			max := 9
			min := 1

			randomNum := rand.Intn(max-min) + min

			if randomNum > 5 {
				initGrid.grid[i][j] = true
			}
		}
	}
}

//some cells in the grid may have neighbors that don't exist, like any neighbor above
//a cell on the top most row, this method checks if given indeces are eithor negative
//or greator than the width/height of the grid and ensures non existant indeces are
//returned false: to be used in below neighborStatus() method
func (g *Generation) isAlive(x, y int) bool {

	if x < 0 || y < 0 {
		return false
	} else if x > (g.width-1) || y > (g.height-1) {
		return false
	} else {
		return g.grid[y][x]
	}

}

//check how many neighbors are alive for given cell, used above isAlive() method
//to ensure non existant neighbors return false
func (g *Generation) neighborStatus(x, y int) bool {

	alive := 0

	//every possible neighbor is now checked and if isAlive() returns true, the alive counter is incremented
	if g.isAlive(x+1, y) {
		alive += 1
	}
	if g.isAlive(x+1, y-1) {
		alive += 1
	}
	if g.isAlive(x+1, y+1) {
		alive += 1
	}
	if g.isAlive(x, y+1) {
		alive += 1
	}
	if g.isAlive(x, y-1) {
		alive += 1
	}
	if g.isAlive(x-1, y) {
		alive += 1
	}
	if g.isAlive(x-1, y+1) {
		alive += 1
	}
	if g.isAlive(x-1, y-1) {
		alive += 1
	}

	//given the number of alive neighbors, return future value of cell
	return alive == 3 || (alive == 2 && g.isAlive(x, y))
}

//struct to store 2 generation grids, one present and one future so chanes can be
//applied without mutating present grid
type Life struct {
	a, b *Generation
	w, h int
}

//creating the initial instance of Generation and a blank future instance to be updated
func Creation(w, h int) *Life {
	return &Life{
		a: makeGrid(w, h), b: makeGrid(w, h),
		w: w, h: h,
	}

}

//function to be called for every instance after initial, reassigning 'a' to the value 'b' points to
//and creating a new blank 'b'
func NewCreation(w int, h int, g *Generation) *Life {
	return &Life{
		a: &*g, b: makeGrid(w, h),
		w: w, h: h,
	}
}

//apply rules and update the future grid, simply get the bool value from the check for each set of
//indeces and assign that value to the same indeces on the new grid
func applyRules(oldGrid, newGrid *Generation) {
	for i := 0; i < oldGrid.height; i++ {
		for j := 0; j < oldGrid.width; j++ {
			newGrid.grid[i][j] = oldGrid.neighborStatus(j, i)
		}
	}
}

//print the grid onto the ternimal
func printGrid(l *Life) {

	output := `

`
	//taking the ansi escape codes and making colord strings based on escape code assigned
	//to the indeces in the hash map
	for i := 0; i < l.h; i++ {
		for j := 0; j < l.w; j++ {

			if j == (l.w - 1) {
				output += " \n"
			} else if l.a.grid[i][j] {
				output += "%"
			} else {
				output += "."
			}

		}
	}

	fmt.Print(output)
}

func funStyles() {

	fmt.Print(`
▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒
▒▒▒▒██████████▒█████████▒████▒▒████▒█████████████▒▒▒
▒▒▒▒█▒▒▒▒▒▒▒▒▒▒█▒▒▒▒▒▒▒█▒█▒▒█▒▒█▒▒█▒█▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒
▒▒▒▒█▒▒▒▒▒▒▒▒▒▒█▒▒▒▒▒▒▒█▒█▒▒█▒▒█▒▒█▒█▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒
▒▒▒▒█▒▒▒▒▒▒▒▒▒▒█▒▒▒▒▒▒▒█▒█▒▒████▒▒█▒█▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒
▒▒▒▒█▒▒███████▒█████████▒█▒▒▒▒▒▒▒▒█▒█████████████▒▒▒
▒▒▒▒█▒▒▒▒▒▒▒▒█▒█▒▒▒▒▒▒▒█▒█▒▒▒▒▒▒▒▒█▒█▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒
▒▒▒▒█▒▒▒▒▒▒▒▒█▒█▒▒▒▒▒▒▒█▒█▒▒▒▒▒▒▒▒█▒█▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒
▒▒▒▒██████████▒█▒▒▒▒▒▒▒█▒█▒▒▒▒▒▒▒▒█▒█████████████▒▒▒
▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒
▒▒▒▒██████████▒██████▒▒▒▒█▒▒███████▒█████▒███████▒▒▒
▒▒▒▒█▒▒▒▒▒▒▒▒█▒█▒▒▒▒▒▒▒▒▒█▒▒▒▒▒█▒▒▒▒█▒▒▒▒▒█▒▒▒▒▒▒▒▒▒
▒▒▒▒█▒▒▒▒▒▒▒▒█▒██████▒▒▒▒█▒▒▒▒▒█▒▒▒▒█████▒███████▒▒▒
▒▒▒▒█▒▒▒▒▒▒▒▒█▒█▒▒▒▒▒▒▒▒▒█▒▒▒▒▒█▒▒▒▒█▒▒▒▒▒█▒▒▒▒▒▒▒▒▒
▒▒▒▒██████████▒█▒▒▒▒▒▒▒▒▒██████████▒█▒▒▒▒▒███████▒▒▒
▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒
▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒


`)
	time.Sleep(1 * time.Second)
	fmt.Println("..... (¯v¯)♥")
	time.Sleep(300 * time.Millisecond)
	fmt.Println(".......•.¸.•´")
	time.Sleep(300 * time.Millisecond)
	fmt.Println("....¸.•´")
	time.Sleep(300 * time.Millisecond)
	fmt.Println("... (")
	time.Sleep(300 * time.Millisecond)
	fmt.Println(" ☻ /")
	time.Sleep(300 * time.Millisecond)
	fmt.Println("/▌♥♥")
	time.Sleep(300 * time.Millisecond)
	fmt.Println("/ \\♥♥")
	time.Sleep(2 * time.Second)
}

func main() {

	//FLAGS
	iterations := flag.Int("i", 500, "number of iterations the game should go through, default is 500")
	speed := flag.Duration("s", 200, "how quickly each iteraton/generation should pass in milliseconds, default is 200")

	flag.Parse()

	iter := *iterations
	spd := *speed

	//setting seed for rng
	rand.Seed(time.Now().UnixNano())

	funStyles()

	//EXECUTION

	//get terminal size
	termSize := terminalSize()
	//create, populate, and print initial grid
	gen := Creation(termSize[0]-1, termSize[1]-1)
	assignValues(gen.a)
	printGrid(gen)

	//develop future grid
	applyRules(gen.a, gen.b)

	newGen := NewCreation(gen.w, gen.h, gen.b)

	for i := 0; i < iter; i++ {

		time.Sleep(spd * time.Millisecond)

		printGrid(newGen)
		applyRules(newGen.a, newGen.b)

		newGen = NewCreation(newGen.w, newGen.h, newGen.b)
	}
}
