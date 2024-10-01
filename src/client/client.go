package client

import (
	"github.com/SoenkeD/rdh/src/core"
)

type Client interface {
}

type InputClient interface {
	LoadSpecs(rdhBe *core.ResourceDefinitionHandler) error
}
