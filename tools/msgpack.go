package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/Kratos40-sba/urlshort/shorturl"
	"github.com/vmihailenco/msgpack"
)

func httpPort() string {
	port := "8000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
}
func main() {
	addres := fmt.Sprintf("http://localhost%s", httpPort())
	redirect := shorturl.Redirect{}
	redirect.URL = "github.com/Kratos40-sba?tab=reposotories"
	body, err := msgpack.Marshal(&redirect)
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := http.Post(addres, "application/x-msgpack", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	msgpack.Unmarshal(body, &redirect)
	log.Printf("%v \n", redirect)
}
