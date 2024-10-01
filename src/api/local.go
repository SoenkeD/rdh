package api

import (
	"fmt"
	"log"
	"time"

	"github.com/SoenkeD/rdh/src/core"
)

type ResourceKindHandler func(id string, specs core.SpecDefinition) (core.MutableSpecDefinition, error)

type Local struct {
	Rdhs           map[string]ResourceKindHandler
	IdleWait       time.Duration
	ReadyReconcile time.Duration
}

func (loc *Local) Init(backend *core.ResourceDefinitionHandler, done chan bool) error {

	go func() {
		for len(done) == 0 {
			err := loc.HandleNext(backend)
			if err != nil {
				panic(err)
			}
		}
	}()

	return nil
}

func (loc *Local) HandleNext(backend *core.ResourceDefinitionHandler) error {

	id, specs, err := backend.GetNext()
	if err != nil {
		// nothing to do - check again later
		time.Sleep(loc.IdleWait)

		return nil
	}

	rkHandler, ok := loc.Rdhs[specs.Kind]
	if !ok {
		return fmt.Errorf("cannot find a resource kind handler for %s", specs.Kind)
	}

	setSpecs, err := rkHandler(id, specs)
	if err != nil {
		log.Println("resource kind handler error", err)

		return err
	}

	if setSpecs.NextReconcile == nil {
		nextReconcile := time.Now()

		if setSpecs.Status == string(core.StateReady) {
			nextReconcile = time.Now().Add(loc.ReadyReconcile)
		}

		setSpecs.NextReconcile = &nextReconcile
	}

	err = backend.SetSpec(id, setSpecs)
	if err != nil {
		return err
	}

	return nil
}
