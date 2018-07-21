package local

import (
	"bufio"
	"github.com/lifei6671/go-storage"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type LocalStorage struct {
}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{}
}

func (m *LocalStorage) Stat(path string, context storage.Context) (storage.FileInfo, error) {
	f, err := os.Stat(path)

	if err != nil {
		return nil, err
	}
	info := NewLocalFileInfo(f)
	info.path = path
	return info, nil

}
func (m *LocalStorage) Write(path string, contents string, context storage.Context) (int64, error) {

	perm := os.FileMode(0755)

	if p, ok := context.Get("perm"); ok {
		if pm, ok := p.(os.FileMode); ok {
			perm = pm
		}
	}

	if err := m.CreateDir(filepath.Dir(path), nil); err != nil {
		return 0, err
	}

	b := []byte(contents)
	err := ioutil.WriteFile(path, b, perm)

	return int64(len(b)), err
}

func (m *LocalStorage) WriteStream(path string, reader io.Reader, size int64, context storage.Context) error {

	perm := os.FileMode(0666)

	if p, ok := context.Get("perm"); ok {
		if pm, ok := p.(os.FileMode); ok {
			perm = pm
		}
	}
	if err := m.CreateDir(filepath.Dir(path), nil); err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)

	if err != nil {
		return err
	}
	read := bufio.NewReader(reader)

	_, err = read.WriteTo(f)

	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

func (m *LocalStorage) WriteBytes(path string, contents []byte, context storage.Context) (int64, error) {
	perm := os.FileMode(0755)

	if p, ok := context.Get("perm"); ok {
		if pm, ok := p.(os.FileMode); ok {
			perm = pm
		}
	}
	if err := m.CreateDir(filepath.Dir(path), nil); err != nil {
		return 0, err
	}

	err := ioutil.WriteFile(path, contents, perm)

	return int64(len(contents)), err
}

func (m *LocalStorage) Append(path string, contents string, context storage.Context) (int64, error) {
	perm := os.FileMode(0666)

	if p, ok := context.Get("perm"); ok {
		if pm, ok := p.(os.FileMode); ok {
			perm = pm
		}
	}
	if err := m.CreateDir(filepath.Dir(path), nil); err != nil {
		return 0, err
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR|os.O_CREATE, perm)

	if err != nil {
		return 0, err
	}

	l, err := f.WriteString(contents)

	if err1 := f.Close(); err == nil {
		err = err1
	}
	return int64(l), err
}

func (m *LocalStorage) AppendStream(path string, reader io.Reader, size int64, context storage.Context) error {
	perm := os.FileMode(0666)

	if p, ok := context.Get("perm"); ok {
		if pm, ok := p.(os.FileMode); ok {
			perm = pm
		}
	}
	if err := m.CreateDir(filepath.Dir(path), nil); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR|os.O_CREATE, perm)

	if err != nil {
		return err
	}

	read := bufio.NewReader(reader)

	_, err = read.WriteTo(f)
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

func (m *LocalStorage) AppendBytes(path string, contents []byte, context storage.Context) (int64, error) {
	perm := os.FileMode(0666)

	if p, ok := context.Get("perm"); ok {
		if pm, ok := p.(os.FileMode); ok {
			perm = pm
		}
	}
	if err := m.CreateDir(filepath.Dir(path), nil); err != nil {
		return 0, err
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR|os.O_CREATE, perm)

	if err != nil {
		return 0, err
	}

	l, err := f.Write(contents)
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return int64(l), err
}
func (m *LocalStorage) ReName(path string, newPath string, context storage.Context) error {
	return os.Rename(path, newPath)
}
func (m *LocalStorage) Copy(src string, dst string, context storage.Context) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}

	if err := m.CreateDir(filepath.Dir(dst), nil); err != nil {
		return err
	}

	desFile, err := os.Create(dst)

	if err != nil {
		return err
	}

	_, err = io.Copy(desFile, srcFile)

	if err1 := desFile.Close(); err == nil {
		err = err1
	}
	if err1 := srcFile.Close(); err == nil {
		err = err1
	}
	return err
}
func (m *LocalStorage) Delete(dst string, context storage.Context) error {
	return os.Remove(dst)
}
func (m *LocalStorage) DeleteDir(dst string, context storage.Context) error {
	return os.RemoveAll(dst)
}
func (m *LocalStorage) CreateDir(dst string, context storage.Context) error {
	if _, err := os.Stat(dst); err == nil {
		return nil
	}
	perm := os.FileMode(0755)

	if p, ok := context.Get("perm"); ok {
		if pm, ok := p.(os.FileMode); ok {
			perm = pm
		}
	}

	err := os.MkdirAll(dst, perm)

	return err
}
func (m *LocalStorage) ListDir(dst string, ctx storage.Context) ([]storage.FileInfo, int64, error) {
	list, err := ioutil.ReadDir(dst)

	if err != nil {
		return nil, 0, err
	}
	num := len(list)

	dirList := make([]storage.FileInfo, num)

	for i, d := range list {
		info := NewLocalFileInfo(d)
		info.path = filepath.Join(dst, d.Name())
		dirList[i] = info
	}

	return dirList, int64(num), nil
}
