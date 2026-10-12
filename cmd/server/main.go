package main
import (
	"log"
	"net/http"
	"inkbottle/conf"
)

func main() {
	fs := http.FileServer(http.Dir(conf.DocsDir))
	http.Handle("/", fs)
	log.Println("preview: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
