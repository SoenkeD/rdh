package client

import (
	"github.com/SoenkeD/rdh/src/core"
)

type InputClient interface {
	LoadSpecs(rdhBe *core.ResourceDefinitionHandler) error
}
