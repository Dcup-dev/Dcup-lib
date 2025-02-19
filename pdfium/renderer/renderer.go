package renderer

import (
	"log"
	"time"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/multi_threaded"
)

// Be sure to close pools/instances when you're done with them.
var pool pdfium.Pool
var  Instance pdfium.Pdfium

func init() {
	// Init the PDFium library and return the instance to open documents.
	// You can tweak these configs to your need. Be aware that workers can use quite some memory.
	pool = multi_threaded.Init(multi_threaded.Config{
		MinIdle:  1, // Makes sure that at least x workers are always available
		MaxIdle:  2, // Makes sure that at most x workers are ever available
		MaxTotal: 5, // Maxium amount of workers in total, allows the amount of workers to grow when needed, items between total max and idle max are automatically cleaned up, while idle workers are kept alive so they can be used directly.
		Command: multi_threaded.Command{
			BinPath: "go",                                     // Only do this while developing, on production put the actual binary path in here. You should not want the Go runtime on production.
			Args:    []string{"run", "pdfium/worker/main.go"}, // This is a reference to the worker package, this can be left empty when using a direct binary path.
		},
	})

	var err error
	Instance, err = pool.GetInstance(time.Second * 30)
	if err != nil {
		log.Fatal(err)
	}
}
