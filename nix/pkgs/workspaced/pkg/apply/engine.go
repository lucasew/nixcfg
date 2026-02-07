package apply

type Engine struct {
	FS FileSystem
}

func NewEngine(fs FileSystem) *Engine {
	if fs == nil {
		fs = OsFileSystem{}
	}
	return &Engine{FS: fs}
}
