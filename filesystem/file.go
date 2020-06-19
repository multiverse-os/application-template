package filesystem

type File struct {
	Directory *Directory
	Path      Path
	Data      []byte
}
