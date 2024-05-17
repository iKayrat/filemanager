package fmsvc

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// UploadFile Request
type UploadFileRequest struct {
	Filename string `json:"filename"`
	Data     []byte `json:"data"`
}

type UploadFileResponse struct {
	Err error `json:"err,omitempty"`
}

// GetListFile Request
type GetListFileRequest struct {
	SortBy string `json:"sortby"`
	Order  string `json:"order"`
}

type GetListFileResponse struct {
	Files []File `json:"files"`
	Err   error  `json:"err,omitempty"`
}

// GetFilename
type SearchFilenameRequest struct {
	Filename string `json:"filename"`
}

type SearchFilenameResponse struct {
	Files []File `json:"files"`
	Err   error  `json:"err,omitempty"`
}

// Download request
type DownloadFileRequest struct {
	Filename string `json:"filename"`
}

type DownloadFileResponse struct {
	FileData []byte `json:"fileData"`
	Err      error  `json:"err,omitempty"`
}

// error
type errorer interface {
	error() error
}

func decodeUploadFileRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req UploadFileRequest

	r.ParseMultipartForm(10 << 20)
	file, fheader, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
	// 	return nil, e
	// }

	req.Filename = fheader.Filename
	req.Data = fileBytes

	return req, nil
}

func decodeGetListFileRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req GetListFileRequest

	vars := mux.Vars(r)
	sort, ok := vars["sortby"]
	if !ok {
		return nil, errors.New("filename not provided")
	}

	orderby, ok := vars["orderby"]
	if !ok {
		return nil, errors.New("filename not provided")
	}

	req.SortBy = sort
	req.Order = orderby

	// if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
	// 	return nil, err
	// }

	return req, nil
}

func decodeSearchFilenameRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req SearchFilenameRequest
	vars := mux.Vars(r)

	filename, ok := vars["filename"]
	if !ok {
		return nil, errors.New("filename not provided")
	}

	req.Filename = filename

	return req, nil
}

func decodeDownloadRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req DownloadFileRequest

	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}

	return req, nil
}

func encodeDownloadFileResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(DownloadFileResponse)
	if resp.Err != nil {
		http.Error(w, resp.Err.Error(), http.StatusInternalServerError)
		return nil
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+"filename.txt")
	w.Header().Set("Content-Type", "application/octet-stream")
	// w.Write(resp.FileData)

	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Fatal("Cannot use flusher")
	}

	_, err := w.Write(resp.FileData)
	if err != nil {
		return err
	}

	flusher.Flush()

	// for i := 1; i < len(resp.FileData); i++ {

	// 	time.Sleep(5 * time.Second)
	// 	w.Write(resp.FileData)
	// 	flusher.Flush()
	// }

	return nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrAlreadyExists, ErrInconsistentIDs:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
