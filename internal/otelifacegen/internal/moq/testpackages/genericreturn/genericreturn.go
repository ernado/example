package genericreturn

import "github.com/ernado/example/internal/otelifacegen/internal/moq/testpackages/genericreturn/otherpackage"

// GenericBar is a test type.
type GenericBar[T any] struct {
	Bar T
}

// IFooBar is a test interface.
type IFooBar interface {
	Foobar() GenericBar[otherpackage.Foo]
}
