package source

type Position struct {
	Line   int
	Column int
}

type Span struct {
	File  *File
	Start Position
	End   Position
}
