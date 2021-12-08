package assets

import "embed"

// F assets for windows or darwin
//go:embed res
var F embed.FS

func demo() {
	//F.ReadDir()
}
