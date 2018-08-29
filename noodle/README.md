# Noodle embed static assets in go binary

## Usage

See [example](_examples/embed) [Makefile](_examples/embed/Makefile) and [main.go](_examples/embed/main.go)

Generate from a single folder, using `gommon generat noodle`

````bash
gommon generate noodle --root assets --output gen/noodle.go --pkg gen --name YangChunMian
````

Generate multiple folders into one file, use `gommon.yml`

````yaml
noodles:
- src: "assets"
  dst: "gen/noodle.go"
  name: "Assets"
  package: "gen"
- src: "third_party"
  dst: "gen/noodle.go"
  name: "ThirdParty"
  package: "gen"
````

````go
import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dyweb/gommon/noodle/_examples/embed/gen"
)

func main() {
	mode := "dev"
	if len(os.Args) > 1 {
		if strings.HasPrefix(os.Args[1], "p") {
			mode = "prod"
		}
	}
	var root http.FileSystem
	if mode == "dev" {
		localDir := "assets"
		root = http.Dir(localDir)
	} else {
		bowel1, err := gen.GetNoodleYangChunMian()
		if err != nil {
			log.Fatal(err)
		}
		root = &bowel1
	}
	addr := ":8080"
	fmt.Printf("listen on %s in %s mode\n", addr, mode)
	fmt.Printf("use http://localhost:8080/index.html")
	log.Fatal(http.ListenAndServe(addr, http.FileServer(root)))
}

````