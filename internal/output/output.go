package output

import (
	"encoding/json"
	"fmt"
	"os"
)

// Formatter handles CLI output in text or JSON format.
type Formatter struct {
	JSON    bool
	Quiet   bool
	NoColor bool
}

func (f *Formatter) Print(data any) {
	if f.JSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(data)
		return
	}
	fmt.Println(data)
}
