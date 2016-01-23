package main
import (
	"github.com/gophergala2016/inspector_gopher"
	"os"
	"log"
)

func main() {
	arg := os.Args[1]

	log.Println(arg)

	coordinator := inspector.NewCoordinator(arg)

	coordinator.Heatmap()
}