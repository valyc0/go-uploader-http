package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Directory di upload e porta configurabili
var uploadDir string
var port string

// Elenca i file nella directory configurabile e permette di scaricarli cliccando
func listFilesHandler(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(uploadDir)
	if err != nil {
		http.Error(w, "Unable to read directory", http.StatusInternalServerError)
		return
	}

	// HTML per elencare i file e il form per l'upload
	fmt.Fprintf(w, "<html><body>")
	fmt.Fprintf(w, "<h2>File List:</h2><ul>")
	for _, file := range files {
		if !file.IsDir() {
			// Crea il link per scaricare ogni file
			fmt.Fprintf(w, `<li><a href="/download?file=%s">%s</a></li>`, file.Name(), file.Name())
		}
	}
	fmt.Fprintf(w, "</ul>")

	// Form per l'upload dei file
	fmt.Fprintf(w, `<h2>Upload File</h2>
		<form enctype="multipart/form-data" action="/upload" method="post">
		<input type="file" name="file" />
		<input type="submit" value="Upload" />
		</form>`)

	fmt.Fprintf(w, "</body></html>")
}

// Gestisce il download dei file
func downloadFileHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	filePath := filepath.Join(uploadDir, fileName)

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found.", http.StatusNotFound)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Error serving file.", http.StatusInternalServerError)
	}
}

// Gestisce l'upload dei file, senza limiti di grandezza
func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Lettura del file senza limite di grandezza
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	destFile, err := os.Create(filepath.Join(uploadDir, header.Filename))
	if err != nil {
		http.Error(w, "Unable to create file on server", http.StatusInternalServerError)
		return
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, file)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Mostra esempio di utilizzo se nessun parametro è passato
func usageExample() {
	fmt.Println("Usage: go run main.go <upload-directory> <port>")
	fmt.Println("Example: go run main.go ./uploads 8080")
}

func main() {
	// Controlla i parametri passati
	if len(os.Args) < 3 {
		usageExample()
		return
	}

	uploadDir = os.Args[1]
	port = os.Args[2]

	// Verifica se la directory esiste, altrimenti la crea
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, 0755)
	}

	// Handlers per elencare, scaricare e caricare file
	http.HandleFunc("/", listFilesHandler)
	http.HandleFunc("/download", downloadFileHandler)
	http.HandleFunc("/upload", uploadFileHandler)

	fmt.Printf("Server running at http://localhost:%s/\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
