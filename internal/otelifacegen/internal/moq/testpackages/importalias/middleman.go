package importalias

import (
	srcclient "github.com/ernado/example/internal/otelifacegen/internal/moq/testpackages/importalias/source/client"
	tgtclient "github.com/ernado/example/internal/otelifacegen/internal/moq/testpackages/importalias/target/client"
)

// MiddleMan is a test interface.
type MiddleMan interface {
	Connect(src srcclient.Client, tgt tgtclient.Client)
}
