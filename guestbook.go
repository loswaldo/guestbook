package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Guestbook struct {
	Count int
	Signatures []string
}

func getStrings(filename string) []string{
	var lines []string
	file, err :=os.Open(filename)
	if err != nil {
		return nil
	}
	check(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	check(scanner.Err())
	return lines
}

func check(err error){
	if err != nil{
		log.Fatal(err)
	}
}

func newHandler(writer http.ResponseWriter, request *http.Request) {
	html, err := template.ParseFiles("new.html")
	check(err)
	err = html.Execute(writer, nil)
	check(err)
}

func createHandler(writer http.ResponseWriter, request *http.Request) {
		signature := request.FormValue("signature")
		options := os.O_WRONLY | os.O_APPEND | os.O_CREATE
		file, err := os.OpenFile("signatures.txt", options, os.FileMode(0600))
		check(err)
		_, err = fmt.Fprintln(file, signature)
		check(err)
		err = file.Close()
		check(err)
		http.Redirect(writer, request, "/guestbook", http.StatusFound)
		//_, err := writer.Write([]byte(signature))
		//check(err)

}

func viewHandler(writer http.ResponseWriter, request *http.Request){
	signatures := getStrings("signatures.txt")
	fmt.Printf("%#v\n", signatures)
	guestbook := Guestbook{
		Count: len(signatures),
		Signatures: signatures,
	}
	html, err := template.ParseFiles("view.html")
	check(err)
	err = html.Execute(writer, guestbook)
	//placeholder := []byte("signature list goes here")
	//_, err := writer.Write(placeholder)
	//check(err)
}

func main() {
	http.HandleFunc("/guestbook", viewHandler)
	http.HandleFunc("/guestbook/new", newHandler)
	http.HandleFunc("/guestbook/create", createHandler)
	err := http.ListenAndServe("localhost:8080", nil)
	log.Fatal(err)
}
