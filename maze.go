package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	Blank  int	= iota;        // not initialised.
	Border;
	Closed;
	Open;
	Visited;
);

type wall struct {
	status			int;
};	

type cell struct {
	status			int;
	RowNum			int;
	ColNum			int;
	EastWall		*wall;
	NorthWall		*wall;
	WestWall		*wall;
	SouthWall		*wall;
};


type maze struct {
	rows 			int;
	cols			int;
	cellWidth		int;
	cellHeight		int;
	matrix			[][]cell;
}


func main () {
	fmt.Println ("starting");
	var _err error;
	var _mazeSize int = 4;
	var _cellWidth int = 4;
	var _cellHeight int = 1;
	if (len (os.Args) > 1) {
		_mazeSize, _err = strconv.Atoi(os.Args[1]);
		if (_err != nil) {
			usage ();
		}
		if (len(os.Args) > 2) {
			_cellWidth, _err = strconv.Atoi(os.Args[2]);
			if (_err != nil) {
				usage ();
			}
		}
		if (len(os.Args) > 3) {
			_cellHeight, _err = strconv.Atoi(os.Args[3]);
			if (_err != nil) {
				usage ();
			}
		}
	}
	_maze := createMaze (_mazeSize);
	_maze.cellWidth = _cellWidth;
	_maze.cellHeight = _cellHeight;

	_maze.recursiveBackTracking ();
	test (_maze);
	_maze.printMaze ();

	fmt.Println ("exiting");
}

func test (_maze maze) {
	// _maze.matrix[1][1].SouthWall.status = Open;
	// _maze.matrix[1][1].NorthWall.status = Open;
	// for i:=0; i<_maze.cols-1; i++ {
	// 	_maze.matrix[0][i].EastWall.status = Open;
	// }
}


func createMaze (size int) maze {

	_matrix := make ([][]cell, size);
	for _row:=0; _row<size; _row++ {
		_matrix[_row] = make ([]cell, size);
		for _col:=0; _col<size; _col++ {
			_matrix[_row][_col].initCell (_row, _col);
		}
	}

	_maze := maze {size, size, 4, 1, _matrix};
	_maze.createMazeBoundary ();
	_maze.setWalls ();
	_maze.setEntryAndExit ();
	return _maze;
}

func (maze maze) createMazeBoundary () {
	var _lastRow int      = maze.rows-1;
	var _numberOfCols int = maze.cols;
	for _col:=0; _col<_numberOfCols; _col++ {
		// set the north boundary
		var  _nWall wall = wall {Border};
		maze.matrix[0][_col].NorthWall = &_nWall;
		// set the south boundary
		var  _sWall wall = wall {Border};
		maze.matrix[_lastRow][_col].SouthWall = &_sWall;
	}

	var _lastCol int  = maze.cols-1;
	var _numberOfRows = maze.rows;
	for i:=0; i<_numberOfRows; i++ {
		// set the west boundary
		var  _wWall wall = wall {Border};
		maze.matrix [i][0].WestWall = &_wWall;
		var  _eWall wall = wall {Border};
		maze.matrix [i][_lastCol].EastWall = &_eWall;
	}
}

func (maze maze) setWalls () {
	var _numberOfCols int = maze.cols;      // we are going to miss the last column
	var _numberOfRows = maze.rows;
	for _row:=0; _row<_numberOfRows; _row++ {
		for _col:=0; _col<_numberOfCols; _col++ {
			if (_col != _numberOfCols-1) {
				// set the east wall. this wall is shared between this cell's east and next cell's west.
				// we will ignore the last column for west-east
				var  _ewWall wall = wall {Closed}; 
				maze.matrix[_row][_col].EastWall = &_ewWall;
				maze.matrix[_row][_col+1].WestWall = &_ewWall
			}
			// we don't need to set south wall for the last row. it should be the boundary
			if (_row != _numberOfRows - 1) {
				var  _snWall = wall {Closed}; 
				maze.matrix[_row][_col].SouthWall = &_snWall;
				maze.matrix[_row+1][_col].NorthWall = &_snWall
			}
		}
	}
}

func (_cell *cell) initCell (row, col int) {
	_cell.status = Blank;
	_cell.EastWall = nil;
	_cell.NorthWall = nil;
	_cell.WestWall = nil;
	_cell.SouthWall = nil;
	_cell.RowNum = row;
	_cell.ColNum = col;
}



func (maze maze) setEntryAndExit  () {
	maze.matrix[0][0].WestWall.status = Open;
	maze.matrix[maze.rows-1][maze.cols-1].EastWall.status = Open;
}


func printBuffer (buffer [][]rune) {
	var _bufferHeight int = len(buffer);

	for _line:=0; _line<_bufferHeight; _line++ {
		fmt.Printf ("%s\n", string(buffer[_line]));
	}
}

func allocatePrintBuffer (maze maze) [][]rune {
	var _lineNum int = (maze.rows*(maze.cellHeight+1))+1;
	var _drawBuffer [][]rune = make ([][]rune, _lineNum);
	for _row:=0; _row<_lineNum; _row++ {
		_drawBuffer [_row] = make ([]rune, (maze.cols*maze.cellWidth)+1);
	}
	return _drawBuffer;
}


// ==============================================================
// ==============================================================
// =============== Printing =====================================
// ==============================================================
// ==============================================================
func (_maze maze) printMaze () [][]rune {
//	var _lineNum int = (maze.rows*maze.cellHeight)+1;
	var _drawBuffer [][]rune = allocatePrintBuffer (_maze);

	drawTopBorder (_maze, _drawBuffer[0]);

	var _row int;
	for _row=0; _row < _maze.rows; _row++ {
		drawACellRow (_maze, _row, _drawBuffer);
	}
	var _bufferHeight int = len(_drawBuffer);
	drawBottomBorder (_maze, _drawBuffer [_bufferHeight-1]); 
	printBuffer (_drawBuffer);
	return _drawBuffer;
}

// ==============================================================
// const topLeftCorner     = 0x250c;
// const topRightCorner    = 0x2510;
// const bottomLeftCorner  = 0x2514;
// const bottomRightCorner = 0x2518;

const topLeftCornerBold     = 0x250f;
const topRightCornerBold    = 0x2513;
const bottomLeftCornerBold  = 0x2517;
const bottomRightCornerBold = 0x251b;

const horizontalLine     = 0x2500;
const horizontalLineBold = 0x2501;
const verticalLine       = 0x2502;
const verticalLineBold   = 0x2503;

// const horizontallTopT   = 0x252c;
// const verticalRightT    = 0x2524;
// const verticalLeftT     = 0x251c;
// const horizontalBottomT = 0x2534;

const horizontallTopTBold   = 0x252f;
const verticalRightTBold    = 0x2528;
const verticalLeftTBold     = 0x2520;
const horizontalBottomTBold = 0x2537;

const Cross			    = 0x253c;
// ==============================================================
func drawTopBorder (maze maze, line []rune) {
	var _lineLen int = len (line);

	// the top 2 corners
	line [0]          = topLeftCornerBold;
	line [_lineLen-1] = topRightCornerBold;

	for i:=1; i<maze.cols; i++ {
		line [i*maze.cellWidth] = horizontallTopTBold;
	}

	for i:=1; i<_lineLen-1; i++ {
		if (line[i] == 0) {
			line [i] = horizontalLineBold;
		}
	}
}


func drawBottomBorder (maze maze, line []rune) {
	var _lineLen int = len (line);

	// the top 2 corners
	line [0]          = bottomLeftCornerBold;
	line [_lineLen-1] = bottomRightCornerBold;

	for i:=1; i<maze.cols; i++ {
		line [i*maze.cellWidth] = horizontalBottomTBold;
	}

	for i:=1; i<_lineLen-1; i++ {
		if (line[i] == 32) {
			line [i] = horizontalLineBold;
		}
	}
}

// this func draws the middle/vertical and the bottom border of a row of cells
// the top border is drawn by the previous row of cells above.
func drawACellRow (maze maze, _row int, drawBuffer [][]rune) {
	var _lineNum int = (_row*(maze.cellHeight+1))+1;
	var _line []rune = drawBuffer [_lineNum];
	var _lineLen int = len (_line);
	// draw the vertical walls
	for h:=0; h<maze.cellHeight; h++ {
		// set the whole line to spaces 
		for i:=0; i<_lineLen; i++ {
			_line [i] = ' ';
		}
		// row 0 has the west wall as the maze entry point
		if (_row > 0) {
			_line [0] = verticalLineBold;
		}
		if (_row == maze.rows-1) {
			_line[_lineLen-1] = ' ';
		} else {
			_line [_lineLen-1] = verticalLineBold;
		}

		for i:=0; i<maze.cols; i++ {
			if (maze.matrix[_row][i].status != Visited) {
				_line[((i+1)*maze.cellWidth)-2] = 'x';
			}
			if (maze.matrix[_row][i].EastWall.status == Closed) {
					_line[((i+1)*maze.cellWidth)] = verticalLine;
			}
		}
		_lineNum++;
		_line = drawBuffer [_lineNum];
	}
	// draw the south walls
	for i:=0; i<_lineLen; i++{
		_line [i] = ' ';
	}
	for i:=0; i<maze.cols; i++ {
		if (maze.matrix[_row][i].SouthWall.status == Closed) {
			var _index int;
			for w:=0; w < maze.cellWidth; w++ {
					_index = (i*maze.cellWidth)+w+1;
					_line[_index] = horizontalLine;
			}
		}
		_line [(i+1)*maze.cellWidth] = Cross;
	}
	_line [0]=verticalLeftTBold;
	_line [_lineLen-1]=verticalRightTBold;

}


func usage () {
	fmt.Println ("usage:");
	fmt.Println ("    maze [matrix size. this will be the size of rows and columns] ");
	fmt.Println ("         [cellS width for drawing purpose]");
	fmt.Println ("         [cellS height for drawing purpose]");
	fmt.Println ("    eg. maze 4 4");
}


// the lower limit is inclusive and upper limit is excluded.
// so the random int is anything from lower to upper-1
func GenerateRandomInt (lower int, upper int) int {
	rand.Seed (time.Now().UnixNano());
	var _rand int = rand.Intn(upper - lower) + lower;
	return _rand;
}


// ==============================================================
// ==============================================================
// =============== recrusive backtracking =======================
// ==============================================================
// ==============================================================
func (_maze maze) unvisitedNeighbours (current *cell) []*cell {
	var _unvisited  []*cell = make ([]*cell, 0);
	var _row		int    = current.RowNum;
	var _col	    int    = current.ColNum;

	// check the western neighbour
	if (_col > 0) {
		if (_maze.matrix [_row][_col-1].status != Visited) {
			_unvisited = append(_unvisited, &_maze.matrix [_row][_col-1]);
		}
	}

	// check the eastern neighbour
	if (_col < _maze.cols-1) {
		if (_maze.matrix[_row][_col+1].status != Visited) {
			_unvisited = append(_unvisited, &_maze.matrix [_row][_col+1]);
		}
	}

	// check the northern neighbour
	if (_row > 0) {
		if (_maze.matrix[_row-1][_col].status != Visited) {
			_unvisited = append(_unvisited, &_maze.matrix [_row-1][_col]);
		}
	}

	// check the southern neighbour
	if (_row < _maze.rows-1) {
		if (_maze.matrix[_row+1][_col].status != Visited) {
			_unvisited = append(_unvisited, &_maze.matrix [_row+1][_col]);
		}
	}

	return (_unvisited);
}

func (_maze maze) recursiveBackTracking () {
	var _row	int = GenerateRandomInt (0, _maze.rows);
	var _col    int = GenerateRandomInt (0, _maze.cols);

	_maze.recursiveBackTrackingProcess (_row,_col);

}


func (_maze maze) openCommonWall (cell1, cell2 *cell) {
	if (cell1.NorthWall == cell2.SouthWall) {
		cell1.NorthWall.status = Open;
	} else if (cell1.SouthWall == cell2.NorthWall) {
		cell1.SouthWall.status = Open;
	} else if (cell1.EastWall == cell2.WestWall) {
		cell1.EastWall.status = Open;
	} else {
		cell1.WestWall.status = Open;
	}
}

// ==============================================================
// process (cell)
// 		start with the cell (x,y]
// 		mark the cell visited
// 		get the unvisited neibours
// 		while > 0 pick one at random
// 			create a opening
//			make it the current cell
// 			process (cell)
//			get the unvisited neibours
// return
// ==============================================================
func (_maze maze) recursiveBackTrackingProcess (row, col int) {
	var _currentCell *cell = &_maze.matrix[row][col];
	_maze.matrix[row][col].status = Visited;
	var _unvisited  []*cell = _maze.unvisitedNeighbours (_currentCell);
	var _unvisitedLen	int = len(_unvisited);

	for ; _unvisitedLen>0; { 
		var _randomIndex int = 0;

		if (_unvisitedLen > 1) {
			_randomIndex = GenerateRandomInt (0, _unvisitedLen);
		}
		var _randomCell *cell = _unvisited [_randomIndex];
		_maze.openCommonWall (_currentCell, _randomCell);
		_maze.recursiveBackTrackingProcess (_randomCell.RowNum, _randomCell.ColNum);
		_unvisited = _maze.unvisitedNeighbours (_currentCell);
		_unvisitedLen = len(_unvisited);
	}
}

// ==============================================================
// ==============================================================
// =============== recrusive backtracking =======================
// ==============================================================
// ==============================================================
