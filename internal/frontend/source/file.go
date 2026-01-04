package source

type File struct {
	Path     string
	Text     []byte
	NumBytes int
}

func NewFile(path string, text []byte) *File {
	return &File{
		Path:     path,
		Text:     text,
		NumBytes: len(text),
	}
}
