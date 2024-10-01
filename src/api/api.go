package api

import (
	"github.com/SoenkeD/rdh/src/client"
	"github.com/SoenkeD/rdh/src/core"
)

type ApiClient interface {
	Init(backend *core.ResourceDefinitionHandler, done chan bool) error
}

type ApiHandler struct {
	be     *core.ResourceDefinitionHandler
	inputs []client.InputClient
	client ApiClient
	done   chan bool
}

type ApiHandlerInput struct {
	Backend *core.ResourceDefinitionHandler
	Inputs  []client.InputClient
	Client  ApiClient
}

func NewApiHandler(input *ApiHandlerInput) *ApiHandler {
	return &ApiHandler{
		be:     input.Backend,
		inputs: input.Inputs,
		client: input.Client,
		done:   make(chan bool, 1),
	}
}

func (ah *ApiHandler) Input() error {
	for _, in := range ah.inputs {
		err := in.LoadSpecs(ah.be)
		if err != nil {

			return err
		}
	}

	return nil
}

func (ah *ApiHandler) Init() error {

	err := ah.client.Init(ah.be, ah.done)
	if err != nil {

		return err
	}

	return nil
}

func (ah *ApiHandler) Stop() {
	ah.done <- true
}
