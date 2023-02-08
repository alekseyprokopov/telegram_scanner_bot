package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"scanner_bot/storage"
	"time"
)

type Storage struct {
	basePath string
}

// права на чтение и запись
const defaultPerm = 0774

var ErrNoSavedPages = errors.New("no saved pages")

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) error {

	filePath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(filePath, defaultPerm); err != nil {
		return fmt.Errorf("can't save create dir: %w", err)

	}

	fName, err := fileName(page)
	if err != nil {
		return fmt.Errorf("can't save create fileName: %w", err)
	}

	filePath = filepath.Join(filePath, fName)

	file, err := os.Create(filePath)

	if err != nil {
		return fmt.Errorf("can't create file: %w", err)
	}
	defer file.Close()

	// write page in file

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return fmt.Errorf("can't encode file: %w", err)
	}
	return nil

}

func (s Storage) Pick(userName string) (page *storage.Page, err error) {
	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)

	if err != nil {
		return nil, fmt.Errorf("can't read dir: %w", err)
	}

	if len(files) == 0 {
		return nil, ErrNoSavedPages
	}

	//генерация случайного файла

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("can't open file for decoding %w", err)
	}
	defer func() { _ = f.Close() }()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, fmt.Errorf("can't decode page: %w", err)
	}

	return &p, nil

}

func (s Storage) Remove(page *storage.Page) error {
	fileName, err := fileName(page)
	if err != nil {
		return fmt.Errorf("can't find file for remove")
	}

	path := filepath.Join(s.basePath, page.UserName, fileName)

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("can't remove file: %w", err)
	}
	return nil
}

func (s Storage) IsExists(page *storage.Page) (bool, error) {
	fileName, err := fileName(page)
	if err != nil {
		return false, fmt.Errorf("can't find file")
	}

	path := filepath.Join(s.basePath, page.UserName, fileName)
	//проверка на наличие файла
	switch _, err := os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file %s exists", path)

		return false, fmt.Errorf(msg, err)

	}
	return true, nil
}
