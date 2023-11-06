package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/lysenkopavlo/telegram-bot/internal/helpers/e"
	"github.com/lysenkopavlo/telegram-bot/internal/storage"
)

const (
	defaultPermission = 0774 // everyone can read and write
)

var ErrNoSavedPages = errors.New("no saved pages")

type Storage struct {
	basePath string
}

// New returns a New Storage
func New(basePath string) Storage {
	return Storage{
		basePath: basePath,
	}
}

// Save saves a given page to new directory
func (s Storage) Save(page *storage.Page) error {

	// choose where to save page
	fPath := filepath.Join(s.basePath, page.UserName)

	// make there a directory
	if err := os.MkdirAll(fPath, defaultPermission); err != nil {
		return e.WrapError("Error while making a directory: %w", err)
	}

	// form a name to file
	fName, err := fileName(page)
	if err != nil {
		return e.WrapError("Error while generating a hash: %w", err)
	}

	fPath = filepath.Join(fPath, fName)

	// create a file to serialize given page
	file, err := os.Create(fPath)
	if err != nil {
		return e.WrapError("Error while creating a file: %w", err)
	}

	// record a page into a file
	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return e.WrapError("Error while encoding a page into file: %w", err)
	}
	return nil
}

// PickRandom returns random article from storage
func (s Storage) PickRandom(userName string) (*storage.Page, error) {
	// where we find page

	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, e.WrapError("Error while reading directory: %w", err)
	}

	if len(files) == 0 {
		return nil, e.WrapError("There is no files: %w", ErrNoSavedPages)
	}

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	n := rng.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(page *storage.Page) error {
	fName, err := fileName(page)
	if err != nil {
		return e.WrapError("Error while generating a hash: %w", err)
	}
	path := filepath.Join(s.basePath, page.UserName, fName)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("Can't remove file at: %s", path)
		return e.WrapError(msg, err)
	}
	return nil
}

func (s Storage) IsExists(page *storage.Page) (bool, error) {
	fName, err := fileName(page)
	if err != nil {
		return false, e.WrapError("Error while generating a hash: %w", err)
	}
	path := filepath.Join(s.basePath, page.UserName, fName)

	switch _, err := os.Stat(path); {
	case err != nil:
		return false, e.WrapError("Error while checking file statistics: %w", err)

	case errors.Is(err, os.ErrNotExist):
		return false, nil
	}

	return true, nil

}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return nil, e.WrapError("Error while opening file %w", err)
	}
	defer file.Close()

	var p *storage.Page

	if err := gob.NewDecoder(file).Decode(&p); err != nil {
		return nil, e.WrapError("Error while decoding file %w", err)
	}
	return p, nil

}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
