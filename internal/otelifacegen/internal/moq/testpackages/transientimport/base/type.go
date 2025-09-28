package base

import (
	four "github.com/ernado/example/internal/otelifacegen/internal/moq/testpackages/transientimport/four/app/v1"
	one "github.com/ernado/example/internal/otelifacegen/internal/moq/testpackages/transientimport/one/v1"
	"github.com/ernado/example/internal/otelifacegen/internal/moq/testpackages/transientimport/onev1"
	three "github.com/ernado/example/internal/otelifacegen/internal/moq/testpackages/transientimport/three/v1"
	two "github.com/ernado/example/internal/otelifacegen/internal/moq/testpackages/transientimport/two/app/v1"
)

// Transient is a test interface.
type Transient interface {
	DoSomething(onev1.Zero, one.One, two.Two, three.Three, four.Four)
}
