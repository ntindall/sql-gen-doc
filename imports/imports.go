package imports

import (
	// required for go mod to be happy when it does tree shaking so that we can
	// build these binaries out of the vendor directory despite the fact that
	// they are not used in the library.
	_ "github.com/pressly/goose/cmd/goose"
	_ "golang.org/x/lint/golint"
)
