package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/mbordner/aoc2025/common"
	"github.com/mbordner/aoc2025/common/array/bytes"
	"github.com/mbordner/aoc2025/common/bits"
	"github.com/mbordner/aoc2025/common/files"
)

func main() {
	start := time.Now()

	cc, presents, checks := calculateCollisions(parseData("../test1.txt"))
	fmt.Println(len(presents), "unique presents calculated")

	count := 0

	for _, check := range checks[2:] {
		counter, _ := bits.NewCounterPack(check.presents) // initialize the counts needed for the required presents in the space (packed bit counters that can count up to 8 values)
		w := check.w - 1
		h := check.h - 1
		if fits(cc, cc.getShapeIds(), NewGrid(w, h), w, counter, 0) {
			count++
		}
	}

	fmt.Printf("for %d tree spaces, we could fit all the presents under %d of them.\n", len(checks), count)
	elapsed := time.Since(start)
	fmt.Printf("took %s\n", elapsed)
}

// fits recursive fits function that implements a DFS search through states, it prunes out branches for shapes we do not need to fit
func fits(cc *CalculatedCollisions, shapeIds []ShapeID, g Grid, width int, counter bits.CounterPack, pos int) bool {
	//state := fmt.Sprintf("%s,%d,%d", g, counter, pos)
	fit := false

posLoop:
	for ; pos < len(g); pos++ {
		checkGrid := g.GetCheckGrid(width, pos) // get the current check grid view (3x3 grid of left, and above shapes)
		for _, shapeId := range shapeIds {
			presentId, _ := shapeId.Extract()                              // figure out which present we are checking
			neededStill, _ := counter.GetCount(presentId)                  // check how many of this present we still need to fit
			if neededStill > 0 && !cc.CheckCollision(shapeId, checkGrid) { // this shape fits here
				newCounter, _ := counter.DecrementCount(presentId) // decrement the present's count which this shape orientation belongs
				// note: the above call only changes counter if the present's id/counter index is > 0, i.e. neededStill was > 0
				// and if after this, all the packed counters are zero, the overall 64bit int will be 0
				if uint64(newCounter) == uint64(0) { // if we fit all the required shapes, we're done
					fit = true
					break posLoop
				}
				newGrid := g.Clone()   // clone a new grid for the DFS
				newGrid[pos] = shapeId // set the fit shape in the new grid
				if fits(cc, shapeIds, newGrid, width, newCounter, pos+1) {
					fit = true // we can break out early since we know there is a configuration that fits
					break posLoop
				}
			}
		}
	}

	return fit
}

type Grid []ShapeID

func NewGrid(w, h int) Grid {
	return make([]ShapeID, w*h)
}

func (g Grid) String() string {
	vals := make([]string, len(g))
	for i, shape := range g {
		vals[i] = shape.String()
	}
	return strings.Join(vals, "")
}

func (g Grid) Print(width int) {
	height := len(g) / width

}

// GetCheckGrid returns a 3x3 grid of shape ids where p is represented by position 2,2; and positions 0,0 up to 2,2 set to the shapes set on g
func (g Grid) GetCheckGrid(w, pos int) [][]ShapeID {
	cg := make([][]ShapeID, 3)
	for r := 0; r < 3; r++ {
		cg[r] = make([]ShapeID, 3)
	}
	row, col := g.PosToRowCol(w, pos)
	for y, r := 2, row; r > row-3; y, r = y-1, r-1 {
		for x, c := 2, col; c > col-3; x, c = x-1, c-1 {
			if r >= 0 && c >= 0 {
				o := g.RowColToPos(w, r, c)
				if g[o].Empty() == false {
					cg[y][x] = g[o]
				}
			}
		}
	}
	return cg
}

func (g Grid) PosToRowCol(w, p int) (int, int) {
	return p / w, p % w
}

func (g Grid) RowColToPos(w, r, c int) int {
	return w*r + c
}

func (g Grid) Clone() Grid {
	o := make(Grid, len(g))
	copy(o, g)
	return o
}

func (cc *CalculatedCollisions) CheckCollision(shapeId ShapeID, checkGrid [][]ShapeID) bool {
	for y := range checkGrid {
		for x := range checkGrid[y] {
			if cc.collisions[shapeId][y][x][checkGrid[y][x]] {
				return true
			}
		}
	}
	return false
}

// ---- pre-calculated collision detection helper classes and functions:

type ShapeID byte

var EmptyShapeID ShapeID = ShapeID(0) // on the grid this special shape id means no shape is placed

func (si ShapeID) Pack(presentID int, shapeIndex int) ShapeID {
	val := byte(16) // initialize with 5th bit set so that when the ids are (0,0) it will not be considered empty

	p := byte(presentID) << 5
	val |= p
	s := byte(shapeIndex)
	val |= s
	return ShapeID(val)
}

func (si ShapeID) Empty() bool {
	return byte(si) == byte(0)
}

func (si ShapeID) Extract() (int, int) {
	var presentID, shapeIndex int

	presentID = int((byte(si) & byte(0xE0)) >> 5)
	shapeIndex |= int(byte(si) & byte(0x0F))

	return presentID, shapeIndex
}

func (si ShapeID) String() string {
	return fmt.Sprintf("%02x", byte(si))
}

type CollisionGridMap [][]map[ShapeID]bool // [row offset][col offset][other shape id][collides bool]

func NewCollisionGridMap() CollisionGridMap {
	cgm := make([][]map[ShapeID]bool, 3)
	for r := range cgm {
		cgm[r] = make([]map[ShapeID]bool, 3)
		for c := range cgm[r] {
			cgm[r][c] = make(map[ShapeID]bool)
		}
	}
	return cgm
}

type CollisionsMap map[ShapeID]CollisionGridMap

type CalculatedCollisions struct {
	collisions CollisionsMap
	shapeIds   []ShapeID
	shapes     map[ShapeID]Shape
}

func (cc *CalculatedCollisions) getShapeIds() []ShapeID {
	sids := make([]ShapeID, len(cc.shapeIds))
	copy(sids, cc.shapeIds)
	return sids
}

func (cc *CalculatedCollisions) addShape(presentID int, shapeIndex int, shape Shape) {
	var shapeId ShapeID
	shapeId = shapeId.Pack(presentID, shapeIndex)
	cc.shapes[shapeId] = shape
	cc.collisions[shapeId] = NewCollisionGridMap()
}

func (cc *CalculatedCollisions) calculateCollisions(sid, oid ShapeID) {
	shapeS := cc.shapes[sid]
	shapeO := cc.shapes[oid]

	grid := common.ConvertGrid([]string{
		`.....`,
		`.....`,
		`.....`,
		`.....`,
		`.....`,
	})

	shapeSCenterX, shapeSCenterY := 3, 3

	// add shape s to empty grid
	for y, j := 2, 0; y < 5; y, j = y+1, j+1 {
		copy(grid[y][2:], shapeS[j])
	}

	for yOffset := 0; yOffset > -3; yOffset-- {
		for xOffset := 0; xOffset > -3; xOffset-- {
			cc.collisions[sid][2+yOffset][2+xOffset][sid] = true           // always would collide with itself
			cc.collisions[sid][2+yOffset][2+xOffset][EmptyShapeID] = false // never collide with empty shape

			shapeOCenterX, shapeOCenterY := shapeSCenterX+xOffset, shapeSCenterY+yOffset
			collision := false
		collisionCheck:
			for h, j := 0, shapeOCenterY-1; j <= shapeOCenterY+1; h, j = h+1, j+1 {
				for g, i := 0, shapeOCenterX-1; i <= shapeOCenterX+1; g, i = g+1, i+1 {
					if grid[j][i] != '.' {
						if shapeO[h][g] == '#' {
							collision = true // s centered at the lower corner of a 3x3 gird
							// would collide with o if o was centered at -yOffset,-xOffset into this 3x3 grid
							break collisionCheck
						}
					}
				}
			}
			cc.collisions[sid][2+yOffset][2+xOffset][oid] = collision
		}
	}
}

// calculateCollisions is where all the collisions are precalculated, every shape against every other shape
func calculateCollisions(presents Presents, checks Checks) (*CalculatedCollisions, Presents, Checks) {
	cc := &CalculatedCollisions{collisions: make(CollisionsMap), shapes: make(map[ShapeID]Shape)}

	for p, present := range presents { // for each of the presents
		for s, shape := range present.shapes { // get all calculated unique orientations of the present's shape
			cc.addShape(p, s, shape) // create a unique identifier from the present id, and shape index, and build up the collision map objects
		}
	}

	cc.shapeIds = make([]ShapeID, 0, len(cc.shapes))
	for sid := range cc.shapes {
		cc.shapeIds = append(cc.shapeIds, sid) // add all unique shape ids that were generated to a list
	}

	shapePairs := common.GetPairSets(cc.shapeIds) // for each pair of shapes, calculate and store if they collide in the maps
	for _, pair := range shapePairs {
		cc.calculateCollisions(pair[0], pair[1])
		cc.calculateCollisions(pair[1], pair[0])
	}

	return cc, presents, checks
}

// ----- data file parsing code and helper classes :

type Shape [][]byte
type Present struct {
	id     int
	shapes []Shape
}

type Presents map[int]*Present

func (s Shape) String() string {
	return fmt.Sprintf("%s\n%s\n%s", string(s[0]), string(s[1]), string(s[2]))
}

type Check struct {
	w        int
	h        int
	presents []int
}

type Checks []Check

func parseData(filename string) (Presents, Checks) {
	presents := make(Presents)
	lines := files.MustGetLines(filename)
	for p, l := 0, 0; p < 6; p, l = p+1, l+5 {
		present := &Present{id: p, shapes: make([]Shape, 0, 8)}

		shapeBytes := make(Shape, 3)
		shapeBytes[0] = []byte(lines[l+1])
		shapeBytes[1] = []byte(lines[l+2])
		shapeBytes[2] = []byte(lines[l+3])

		present.shapes = append(present.shapes, shapeBytes)

		visited := make(common.VisitedState[string, bool])
		visited.Set(shapeBytes.String(), true)

		var next Shape = shapeBytes
		for i := 0; i < 3; i++ {
			next = Shape(bytes.Rotate(next))
			ns := next.String()
			if has := visited.Has(ns); !has {
				present.shapes = append(present.shapes, next)
				visited.Set(ns, true)
			}
		}

		for _, d := range []bytes.Direction{bytes.Horizontal, bytes.Vertical} {
			next = Shape(bytes.Flip(d, shapeBytes))
			for i := 0; i < 4; i++ {
				next = Shape(bytes.Rotate(next))
				ns := next.String()
				if has := visited.Has(ns); !has {
					present.shapes = append(present.shapes, next)
					visited.Set(ns, true)
				}
			}
		}

		presents[present.id] = present
	}

	replacer := strings.NewReplacer("x", ",", ":", "", " ", ",")
	checks := make(Checks, 0, len(lines)-30)
	for _, line := range lines[30:] {
		check := Check{}
		vals := common.IntVals[int](replacer.Replace(line))
		check.w = vals[0]
		check.h = vals[1]
		check.presents = vals[2:]
		checks = append(checks, check)
	}

	return presents, checks
}
