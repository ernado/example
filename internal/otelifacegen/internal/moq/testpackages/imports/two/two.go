package two

import (
	"github.com/ernado/example/internal/otelifacegen/internal/moq/testpackages/imports/one"
)

// DoSomething does something.
type DoSomething interface {
	Do(thing one.Thing) error
	Another(thing one.Thing) error
}
