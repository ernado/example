package typealiastwo

import "github.com/ernado/example/internal/otelifacegen/internal/moq/testpackages/typealiastwo/internal/typealiasinternal"

type AliasType = typealiasinternal.MyInternalType

type GenericAliasType = typealiasinternal.MyGenericType[int]
