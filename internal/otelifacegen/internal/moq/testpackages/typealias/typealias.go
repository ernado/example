package typealias

import (
	"github.com/ernado/example/internal/otelifacegen/internal/moq/testpackages/typealiastwo"
)

type Example interface {
	Do(a typealiastwo.AliasType, b typealiastwo.GenericAliasType) error
}
