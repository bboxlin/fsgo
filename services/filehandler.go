package services

import (
	"fmt"
	"io"
	"net/http"

	"github.com/unibytes/fsgo/utils"
)

type StreamIterator struct {
	reader io.Reader
	buffer []byte
}

func NewStreamIterator(reader io.Reader, bufferSize int) *StreamIterator {
	return &StreamIterator{
		reader: reader,
		buffer: make([]byte, bufferSize),
	}
}

func (i *StreamIterator) NextChunk() ([]byte, error) {
	n, err := i.reader.Read(i.buffer)
	if err != nil && err != io.EOF {
		return nil, err
	}
	if n == 0 {
		return nil, io.EOF
	}
	return i.buffer[:n], nil
}

type FileHandler struct {
	basePath string
}

func NewFileHandler(basePath string) *FileHandler {
	return &FileHandler{
		basePath: basePath,
	}
}

func (h *FileHandler) Upload(w http.ResponseWriter, r *http.Request) {
	// Get the file from the request body
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create the "static" folder if it doesn't exist
	if err := utils.CreateStaticFolder(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create the file in the "static" folder
	newFile, err := utils.CreateFile(header.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer newFile.Close()

	// Create a stream iterator to read the file in chunks
	iterator := NewStreamIterator(file, 1024*1024) // 1MB buffer

	// Write the data to the new file
	for {
		data, err := iterator.NextChunk()
		if err != nil {
			if err == io.EOF {
				break
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := newFile.Write(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Return a success response
	fmt.Fprintf(w, "File uploaded successfully!")
}
