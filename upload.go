package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("METHOD IS " + r.Method + " AND CONTENT-TYPE IS " + r.Header.Get("Content-Type"))
	fmt.Printf("====Upload File ===\n")
	file, _, err := r.FormFile("file")
	fmt.Println("METHOD IS ", file)

	if err != nil {
		fmt.Println("Error ", err)
		return
	}
	defer file.Close()
	tempFile, err := ioutil.TempFile("images", "upload-*.png")
	if err != nil {
		log.Println("error 1", err)
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("error 2", err)
		return
	}

	tempFile.Write(fileBytes)
	data := map[string]string{
		"msg":   "Hello Golang",
		"files": tempFile.Name(),
	}

	json.NewEncoder(w).Encode(&data)
}
func getAll(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir("images")
	if err != nil {
		log.Fatal(err)
	}
	output := []string{}
	for _, f := range files {
		output = append(output, f.Name())
	}
	json.NewEncoder(w).Encode(&output)

}
func getFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	file, err := ioutil.ReadFile("images/" + vars["file"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-type", "image/png")
	w.Write(file)
}

func hello(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"msg": "Hello Golang",
	}

	json.NewEncoder(w).Encode(&data)
}

func setupRoutes() {
	router := mux.NewRouter()
	router.HandleFunc("/", hello).Methods(http.MethodGet)
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})
	router.HandleFunc("/getfile/", getAll).Methods(http.MethodGet)
	router.HandleFunc("/getfile/{file}", getFile).Methods(http.MethodGet)
	router.HandleFunc("/upload/", uploadFile).Methods(http.MethodPost)
	
	handler := cors.Handler(router)
	log.Println("listen port :1234")
	http.ListenAndServe(":1234", handler)
}

func main() {
	setupRoutes()
}
