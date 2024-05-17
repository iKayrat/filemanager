package fmsvc

import (
	"context"
	"database/sql"
	"errors"
	random "filemanagerService/pkg/fmsvc"

	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
)

type Service interface {
	Upload(ctx context.Context, filename string, data []byte) error
	Download(ctx context.Context, filename string) ([]byte, error)
	SearchByName(ctx context.Context, filename string) ([]File, error)
	List(ctx context.Context, sortBy, order string) ([]File, error)
}

type FManager struct {
	FStore FileStore
}

func New(conn *sql.DB) Service {

	fs := NewFileStorage(conn)

	return &FManager{
		FStore: fs,
	}
}

func (fm *FManager) Upload(ctx context.Context, filename string, data []byte) error {
	//todo

	// random string
	uid := uuid.New().String()
	id, _ := strconv.ParseInt(uid, 10, 64)

	//full filename
	secondName := concatExt(filename)

	filePath := filepath.Join("uploads", secondName)

	file := File{
		Id:         id,
		FileName:   filename,
		SecondName: secondName,
		Path:       filePath,
	}

	err := func() error {
		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0666))
		if err != nil {
			// log.Println("err:", err)
			return err
		}

		_, err = f.Write(data)
		if err1 := f.Close(); err1 != nil && err == nil {
			err = err1
		}
		// log.Println("err:", err)
		return err
	}()
	if err != nil {
		return err
	}

	err = func() error {
		size, _ := os.Stat(filePath)
		file.Size = size.Size()
		file.UploadedAt = size.ModTime()

		// log.Println(" file:", file.Path)
		return fm.FStore.Save(ctx, file)
	}()
	if err != nil {
		return err
	}

	return nil
}

func (fm *FManager) SearchByName(ctx context.Context, filename string) ([]File, error) {

	// get list of files by size from db
	filesFromDB, err := fm.FStore.SerachByName(ctx, filename)
	if err != nil {
		return nil, err
	}

	return filesFromDB, nil
}

func (fm *FManager) List(ctx context.Context, column, orderby string) ([]File, error) {
	//todo
	// filesFromDir := make([]File, 0)

	// // get list of files from local directory
	// err := filepath.Walk("uploads", func(path string, info os.FileInfo, err error) error {
	// 	if !info.IsDir() {
	// 		filesFromDir = append(filesFromDir, File{
	// 			SecondName: info.Name(),
	// 			Size:       info.Size(),
	// 			UploadedAt: info.ModTime(),
	// 		})
	// 	}

	// 	return nil
	// })
	// if err != nil {
	// 	log.Println("list m err", err)
	// 	return nil, err
	// }
	// log.Println("service list...")

	// get list of files by size from db
	filesFromDB, err := fm.FStore.GetListBy(ctx, column, orderby)
	if err != nil {
		return nil, err
	}

	// log.Println("filenames", filesFromDB)
	return filesFromDB, nil
}

func (fm *FManager) Download(ctx context.Context, filename string) ([]byte, error) {
	//todo

	getfile, err := fm.FStore.GetFile(ctx, filename)
	if err != nil {
		return nil, err
	}

	filePath := "uploads/" + getfile.SecondName

	fileData, err := os.ReadFile(filePath)
	if os.IsNotExist(err) {
		return nil, err
	} else if err != nil {
		return nil, err
	}

	return fileData, nil
}

// look for local and
// func getOrigNames(fromlocal []File, fromdb []File) (out []File) {

// 	for _, db := range fromdb {
// 		for _, local := range fromlocal {
// 			if local.SecondName == db.SecondName {
// 				out = append(out, File{
// 					FileName:   db.FileName,
// 					SecondName: db.SecondName,
// 					Path:       db.Path,
// 					Size:       db.Size,
// 					UploadedAt: db.UploadedAt,
// 				})
// 			}
// 		}
// 	}

// 	return
// }

// Concatenation of filename and extension
func concatExt(filename string) string {

	genStr := random.String(filename)

	ext := filepath.Ext(filename)
	var str strings.Builder

	str.WriteString(genStr)
	str.WriteString(ext)

	return str.String()
}
