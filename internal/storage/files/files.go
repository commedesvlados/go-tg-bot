package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"github/commedesvlados/go-tg-bot/internal/storage"
	"github/commedesvlados/go-tg-bot/pkg/lib/e"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type Storage struct {
	basePath string
}

const defaultPerm = 0774 // read/ write for all users

func NewStorage(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) error {
	fPath := filepath.Join(s.basePath, page.Username)

	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return e.Wrap("can't create directory", err)
	}

	fName, err := filename(page)
	if err != nil {
		return e.Wrap("can't get filename", err)
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return e.Wrap("can't create file", err)
	}
	defer file.Close() // Unhandled error

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return e.Wrap("can't encode data", err)
	}

	return nil
}

func (s Storage) PickRandom(username string) (*storage.Page, error) {
	path := filepath.Join(s.basePath, username)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		// dir does not exist
		return nil, storage.ErrNoSavedPages
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, e.Wrap("can't pick random page", err)
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	// TODO - rand.New(rand.NewSource(0-9))
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(page *storage.Page) error {
	filename, err := filename(page)
	if err != nil {
		return e.Wrap("cant' remove file", err)
	}

	path := filepath.Join(s.basePath, page.Username, filename)

	if err := os.Remove(path); err != nil {
		errmsg := fmt.Sprintf("cant' remove file %s", path)
		return e.Wrap(errmsg, err)
	}

	return nil
}

func (s Storage) IsExists(page *storage.Page) (bool, error) {
	filename, err := filename(page)
	if err != nil {
		return false, e.Wrap("cant' check if file exists", err)
	}

	path := filepath.Join(s.basePath, page.Username, filename)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		errmsg := fmt.Sprintf("cant' check if file exists: %s", path)
		return false, e.Wrap(errmsg, err)
	}

	return true, nil
}

func (s Storage) decodePage(filepath string) (*storage.Page, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, e.Wrap("can't decode page", err)
	}
	defer f.Close()

	var p storage.Page
	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

func filename(p *storage.Page) (string, error) {
	return p.Hash()
}
