package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type Message struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Println("hello")
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func writeDummyFile(w http.ResponseWriter, req *http.Request) {

	fmt.Println("writeDummyFile")
	d1 := []byte("the long and winding road")
	err := ioutil.WriteFile("/tmp/dat1", d1, 0644)
	check(err)

	f, err := os.Create("/tmp/goServer/test1")
	check(err)
	defer f.Close()
	d2 := []byte{115, 111, 109, 101, 10}
	n2, err := f.Write(d2)
	check(err)
	fmt.Printf("wrote %d bytes\n", n2)
	fmt.Println("writeDummyFile")
}

func readReq(w http.ResponseWriter, req *http.Request) {
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var msg Message
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	output, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	f, err := os.Create(fmt.Sprint("/tmp/goServer/test_json", msg.Id))
	check(err)
	defer f.Close()
	//f_xml, err := os.Create(fmt.Sprint("/tmp/goServer/test_xml", msg.Id))
	//check(err)
	//defer f_xml.Close()
	f.Write(output)
	w.Write(output)
	//out, err := xml.MarshalIndent(output, "\t", "\t")
	//fmt.Println(out)
}
func savaFile(w http.ResponseWriter, r *http.Request) {
	// the FormFile function takes in the POST input id file
	file, header, err := r.FormFile("file")
	check(err)
	defer file.Close()

	out, err := os.Create(header.Filename)
	check(err)
	defer out.Close()

	// write the content from POST to the file
	_, err = io.Copy(out, file)
	check(err)
	fmt.Fprintf(w, "File uploaded successfully: ")
	fmt.Fprintf(w, header.Filename)
}
func main() {

	fmt.Println("hey yo")
	http.HandleFunc("/writeDummyFile", writeDummyFile)
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/WriteToFile", readReq)
	http.HandleFunc("/saveFile", savaFile)
	err := http.ListenAndServe(":8090", nil)
	check(err)
}
