/*
 * File:    main.go
 * Date:    October 13, 2024
 * Author:  J.
 * Email:   jaime.gomez@usach.cl
 * Project: goHttpUploadServer - Drone FrameWork
 *
 * Build Win:
 * 		go build -ldflags "-s -w" main.go
 *
 * Build Linux:
 * 		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o main main.go
 *
 * Do not forget to create c:/uploads/ under the host
 *
 * Docker Linux with WSL:
 *		docker build -t uploadserver .
 * 		docker run --name uploadserver -p 3333:3333 -v /mnt/c/uploads:/uploadserver/uploads uploadserver
 *
 * Docker testing after build:
 * 		docker run -it --entrypoint /bin/sh uploadserver
 *		ls -l
 *		./main
 *
 */

package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// const uploadDir = `C:\uploads`
const uploadDir = `./uploads`

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./public")))
	mux.HandleFunc("/hello", getHello)
	mux.HandleFunc("/upload", uploadFiles)

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		fmt.Printf("Folder %s does not exist. It will be created\n", uploadDir)
		errF := os.Mkdir(uploadDir, 0755)
		if errF != nil {
			log.Fatal(errF)
		}
	}

	fmt.Println("UploadServer started at :3333")
	err := http.ListenAndServe(":3333", mux)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
		os.Exit(1)
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}

}

func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s\n", r.Method)
	fmt.Printf("got /hello request\n")
	_, _ = io.WriteString(w, "Hello http!\n")
}

func uploadFiles(w http.ResponseWriter, r *http.Request) {

	var finalDir string

	if r.Method != "POST" {
		fmt.Printf("Invalid method: %s\n", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	folderName := r.FormValue("folderName")
	fmt.Println("Creating folder: ", folderName)
	if folderName != "" {
		finalDir = filepath.Join(uploadDir, folderName)
		os.Mkdir(finalDir, os.ModePerm)
	} else {
		finalDir = uploadDir
	}

	fmt.Println("Handling file upload...")

	// 10MB max upload size >> moved to front-end in JS
	//err := r.ParseMultipartForm(10 << 20)
	//if err != nil {
	//	fmt.Printf("Error parsing form: %v\n", err)
	//	http.Error(w, "Error parsing form data", http.StatusInternalServerError)
	//	return
	//}

	// Retrieve files from the "photos" field
	files := r.MultipartForm.File["photos"]
	if files == nil {
		fmt.Println("No files found in request")
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received %d files\n", len(files))

	// Process each file
	for _, fileHeader := range files {
		fmt.Printf("Processing file: %s\n", fileHeader.Filename)

		// Open the file
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Create the destination file
		dst, err := os.Create(filepath.Join(finalDir, fileHeader.Filename))
		if err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Copy the uploaded file data to the destination file
		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Error writing file", http.StatusInternalServerError)
			return
		}
	}

	// Respond with a success message
	fmt.Fprintf(w, "Files uploaded successfully!")
}
