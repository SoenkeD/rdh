package core

import (
	"fmt"
	"time"
)

type State string

const StateNew State = "New"
const StateReady State = "Ready"
const StateError State = "Error"

type CreationSpecDefinition struct {
	Annotations map[string]string
	Labels      map[string]string
	Specs       any
}

type MutableSpecDefinition struct {
	NextReconcile *time.Time
	Status        string

	CreationSpecDefinition
}

type ControlledSpecDefinition struct {
	UpdatedAt time.Time
}

type SpecDefinition struct {
	Kind      string
	CreatedAt time.Time

	ControlledSpecDefinition
	MutableSpecDefinition
}

type ResourceDefinitionHandler struct {
	Specs map[string]SpecDefinition
}

func NewResourceDefinitionHandler() *ResourceDefinitionHandler {
	return &ResourceDefinitionHandler{
		Specs: map[string]SpecDefinition{},
	}
}

func (rdh *ResourceDefinitionHandler) CreateSpec(id, kind string, specs CreationSpecDefinition) error {

	_, ok := rdh.Specs[id]
	if ok {
		return fmt.Errorf("spec definition with id %s already exists", id)
	}

	nextReconcile := time.Now()

	rdh.Specs[id] = SpecDefinition{
		Kind:      kind,
		CreatedAt: time.Now(),
		ControlledSpecDefinition: ControlledSpecDefinition{
			UpdatedAt: time.Now(),
		},
		MutableSpecDefinition: MutableSpecDefinition{
			NextReconcile:          &nextReconcile,
			Status:                 string(StateNew),
			CreationSpecDefinition: specs,
		},
	}

	return nil
}

func (rdh *ResourceDefinitionHandler) SetSpec(id string, mSpec MutableSpecDefinition) error {

	spec, err := rdh.GetSpec(id)
	if err != nil {
		return err
	}

	spec.MutableSpecDefinition = mSpec
	spec.ControlledSpecDefinition.UpdatedAt = time.Now()

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
