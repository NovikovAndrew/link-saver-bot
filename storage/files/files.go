package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"bot-saver/package/error"
	"bot-saver/storage"
)

const defaultPerm = 0774

var ErrNoSavedPage = errors.New("no saved pages")

type Storage struct {
	basePath string
}

func New(basePath string) storage.Storage {
	return &Storage{
		basePath: basePath,
	}
}

func (s Storage) Save(p *storage.Page) (err error) {
	defer func() { err = e.WrapIfErr("can't to save page", err) }()

	filePath := filepath.Join(s.basePath, p.UserName)

	if err = os.MkdirAll(filePath, defaultPerm); err != nil {
		return err
	}

	fName, err := fileName(p)

	if err != nil {
		return err
	}

	file, err := os.Create(fName)

	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(p); err != nil {
		return err
	}

	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("can't to PickRandom", err) }()
	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)

	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, ErrNoSavedPage
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files) - 1)
	file := files[n]

	return s.decode(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)

	if err != nil {
		return e.Wrap("can't to get fileName hash", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(path); err != nil {
		return e.Wrap(fmt.Sprintf("can't to delete file, %s", path), err)
	}

	return err
}

func (s Storage) IsExist(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("can't to get fileName hash", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err := os.Stat(path); {
	case errors.Is(err, ErrNoSavedPage):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't to find filepath, filepath %s", path)
		return false, e.Wrap(msg, err)
	}

	return true, nil
}

func (s Storage) decode(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("can't to open file", err)
	}

	defer func() { _ = f.Close() }()

	var page storage.Page

	if err := gob.NewDecoder(f).Decode(&page); err != nil {
		return nil, e.Wrap("can't to decode file", err)
	}

	return &page, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
