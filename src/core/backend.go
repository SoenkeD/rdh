package core

import (
	"fmt"
	"log"
	"time"
)

type State string

const StateNew State = "New"
const StateReady State = "Ready"

type SpecDefinition struct {
	Kind          string
	UpdatedAt     time.Time
	NextReconcile *time.Time
	Status        string
	Specs         any
}

type ResourceDefinitionHandler struct {
	Specs map[string]SpecDefinition
}

func NewResourceDefinitionHandler() *ResourceDefinitionHandler {
	return &ResourceDefinitionHandler{
		Specs: map[string]SpecDefinition{},
	}
}

func (rdh *ResourceDefinitionHandler) CreateSpec(id, kind string, specs any) error {

	_, ok := rdh.Specs[id]
	if ok {
		return fmt.Errorf("spec definition with id %s already exists", id)
	}

	nextReconcile := time.Now()

	rdh.Specs[id] = SpecDefinition{
		Kind:          kind,
		UpdatedAt:     time.Now(),
		NextReconcile: &nextReconcile,
		Status:        string(StateNew),
		Specs:         specs,
	}

	return nil
}

func (rdh *ResourceDefinitionHandler) SetSpec(id string, spec SpecDefinition) error {

	log.Println("Update spec", id, spec)

	rdh.Specs[id] = spec

	return nil
}

func (rdh *ResourceDefinitionHandler) GetSpec(id string) (SpecDefinition, error) {

	spec, ok := rdh.Specs[id]
	if !ok {
		return SpecDefinition{}, fmt.Errorf("spec with id %s does not exists", id)
	}

	return spec, nil
}

func (rdh *ResourceDefinitionHandler) GetNext() (string, SpecDefinition, error) {

	initTime := time.Now()
	var selectedId string

	for id, specDef := range rdh.Specs {
		if specDef.NextReconcile.Before(initTime) {
			selectedId = id
			initTime = *specDef.NextReconcile
		}
	}

	if selectedId == "" {
		return "", SpecDefinition{}, fmt.Errorf("there is no spec to be reconciled now")
	}

	return selectedId, rdh.Specs[selectedId], nil
}
