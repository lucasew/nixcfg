package apply

import (
	"os"
)

type FileSystem interface {
	Lstat(name string) (os.FileInfo, error)
	Readlink(name string) (string, error)
	ReadFile(name string) ([]byte, error)
	MkdirAll(path string, perm os.FileMode) error
	Remove(name string) error
	Rename(oldpath, newpath string) error
	Symlink(oldname, newname string) error
	WriteFile(name string, data []byte, perm os.FileMode) error
	Stat(name string) (os.FileInfo, error)
}

type OsFileSystem struct{}

func (OsFileSystem) Lstat(name string) (os.FileInfo, error) {
	return os.Lstat(name)
}

func (OsFileSystem) Readlink(name string) (string, error) {
	return os.Readlink(name)
}

func (OsFileSystem) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (OsFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (OsFileSystem) Remove(name string) error {
	return os.Remove(name)
}

func (OsFileSystem) Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

func (OsFileSystem) Symlink(oldname, newname string) error {
	return os.Symlink(oldname, newname)
}

func (OsFileSystem) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (OsFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}
