//-----------------------------------------------------------------------------
/*

Connect Four

*/
//-----------------------------------------------------------------------------

package cf

import (
	"errors"
	"strings"
)

//-----------------------------------------------------------------------------

const nRows = 6
const nCols = 7
const nPlayers = 2
const nConnect = 4

//-----------------------------------------------------------------------------

// posnIndex returns the index for a single board position.
func posnIndex(row, col int) int {
	return (row * nCols) + col
}

// posnBit returns the bitmask for a single board position.
func posnBit(row, col int) uint64 {
	return 1 << posnIndex(row, col)
}

// posnBitCheck returns the bitmask for a single board position, with row and column checks.
func posnBitCheck(row, col int) (uint64, error) {
	if row < 0 || row >= nRows {
		return 0, errors.New("invalid row")
	}
	if col < 0 || col >= nCols {
		return 0, errors.New("invalid col")
	}
	return posnBit(row, col), nil
}

//-----------------------------------------------------------------------------

// bmHorz returns a horizontal bitmap.
func bmHorz(row, col int) (uint64, error) {
	// left to right
	var mask uint64
	for i := 0; i < nConnect; i++ {
		bit, err := posnBitCheck(row, col+i)
		if err != nil {
			return 0, err
		}
		mask |= bit
	}
	return mask, nil
}

// bmVert returns a vertical bitmap.
func bmVert(row, col int) (uint64, error) {
	// top to bottom
	var mask uint64
	for i := 0; i < nConnect; i++ {
		bit, err := posnBitCheck(row-i, col)
		if err != nil {
			return 0, err
		}
		mask |= bit
	}
	return mask, nil
}

// bmDiag0 returns a diagonal bitmap.
func bmDiag0(row, col int) (uint64, error) {
	// bottom left to top right
	var mask uint64
	for i := 0; i < nConnect; i++ {
		bit, err := posnBitCheck(row+i, col+i)
		if err != nil {
			return 0, err
		}
		mask |= bit
	}
	return mask, nil
}

// bmDiag1 returns a diagonal bitmap.
func bmDiag1(row, col int) (uint64, error) {
	// top left to bottom right
	var mask uint64
	for i := 0; i < nConnect; i++ {
		bit, err := posnBitCheck(row-i, col+i)
		if err != nil {
			return 0, err
		}
		mask |= bit
	}
	return mask, nil
}

// genWins generates the win bitmasks for each position.
func genWins() [][]uint64 {
	wins := make([][]uint64, nRows*nCols)
	for row := 0; row < nRows; row++ {
		for col := 0; col < nCols; col++ {
			idx := posnIndex(row, col)
			// horizontal
			for i := 0; i < nConnect; i++ {
				mask, err := bmHorz(row, col-i)
				if err == nil {
					wins[idx] = append(wins[idx], mask)
				}
			}
			// vertical
			mask, err := bmVert(row, col)
			if err == nil {
				wins[idx] = append(wins[idx], mask)
			}
			// diagonal (bottom left to top right)
			for i := 0; i < nConnect; i++ {
				mask, err := bmDiag0(row-i, col-i)
				if err == nil {
					wins[idx] = append(wins[idx], mask)
				}
			}
			// diagonal (top left to bottom right)
			for i := 0; i < nConnect; i++ {
				mask, err := bmDiag1(row+i, col-i)
				if err == nil {
					wins[idx] = append(wins[idx], mask)
				}
			}
		}
	}
	return wins
}

var winSet [][]uint64

func init() {
	winSet = genWins()
	/*
	   for i, v := range winSet {
	     for _, k := range v {
	       g := NewGame()
	       g.player[0] = k
	       fmt.Printf("%d\n%s\n", i, g)
	     }
	   }
	*/
}

//-----------------------------------------------------------------------------

// Game represents the connect4 game state.
type Game struct {
	player    [nPlayers]uint64 // player bit fields
	colHeight [nCols]int       // height of pieces in each column
	last      int              // index of last move
}

func (g *Game) String() string {
	rows := make([]string, nRows)
	for i := range rows {
		cols := make([]rune, nCols)
		for j := range cols {
			posn := posnBit(nRows-1-i, j)
			cols[j] = '.'
			if g.player[0]&posn != 0 {
				cols[j] = 'O'
			}
			if g.player[1]&posn != 0 {
				cols[j] = 'X'
			}
		}
		rows[i] = string(cols)
	}
	return strings.Join(rows, "\n")
}

// NewGame returns a new connect4 game state.
func NewGame() *Game {
	return &Game{}
}

// Add a player piece to the game.
func (g *Game) Add(player, col int) error {
	if col < 0 || col >= nCols {
		return errors.New("invalid col")
	}
	if g.colHeight[col] >= nRows {
		return errors.New("no col space")
	}
	g.last = posnIndex(g.colHeight[col], col)
	g.player[player] |= 1 << g.last
	g.colHeight[col] += 1
	return nil
}

// Win returns if the player won with the last move,
func (g *Game) Win(player int) bool {
	for _, mask := range winSet[g.last] {
		if g.player[player]&mask == mask {
			return true
		}
	}
	return false
}

//-----------------------------------------------------------------------------
