package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	http.HandleFunc("/upload", fileUploadHandler)
	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func fileUploadHandler(w http.ResponseWriter, r *http.Request) {
	// Limit file size to 10MB
	r.ParseMultipartForm(10 << 20)

	// Retrieve the file from form data
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read the file into a byte slice to validate its type
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}

	// Validate file type
	if !isValidFileType(fileBytes) {
		http.Error(w, "Invalid file type", http.StatusUnsupportedMediaType)
		return
	}

	fmt.Fprintf(w, "Uploaded File: %s\n", handler.Filename)
	fmt.Fprintf(w, "File Size: %d\n", handler.Size)
	fmt.Fprintf(w, "MIME Header: %v\n", handler.Header)

	// Now let’s save it locally
	dst, err := createFile(handler.Filename)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Write the file bytes to the destination file
	if _, err := dst.Write(fileBytes); err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "File uploaded successfully!")
}

func createFile(filename string) (*os.File, error) {
	// Create an uploads directory if it doesn’t exist
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", 0755)
	}

	// Build the file path and create it
	dst, err := os.Create(filepath.Join("uploads", filename))
	if err != nil {
		return nil, err
	}

	return dst, nil
}

func isValidFileType(file []byte) bool {
	fileType := http.DetectContentType(file)
	return strings.HasPrefix(fileType, "image/") // Only allow images
}
