package main
import (
	"os"
	"fmt"
	"flag"
	"path/filepath"
	"github.com/gophergala2016/inspector_gopher"
	"io/ioutil"
)

var fileRevisions = []inspector.FileRevision{}

func visit(path string, f os.FileInfo, err error) error {
	if filepath.Ext(path) == ".go" {
		data, _ := ioutil.ReadFile(path)
		fileRevisions = append(fileRevisions, inspector.ParseFileContents(path, string(data)))
	}

	return nil
}

func main() {
	flag.Parse()
	root := flag.Arg(0)
	filepath.Walk(root, visit)
	fmt.Println("", len(fileRevisions))

	for _, val := range fileRevisions {
		for _, v := range val.Units {
			fmt.Printf("%d:%s:%d:%d\n", v.Type, v.Name, v.LineStart, v.LineEnd)
		}
	}
}
