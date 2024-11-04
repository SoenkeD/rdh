package core_test

import (
	"time"

	"github.com/SoenkeD/rdh/src/core"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RDH - Create Spec", func() {

	var rdh core.ResourceDefinitionHandler

	BeforeEach(func() {
		rdh = *core.NewResourceDefinitionHandler()
	})

	It("Creates a spec", func() {
		Expect(rdh.CreateSpec("id-1", "demo", core.CreationSpecDefinition{
			Specs: "val",
		})).To(Succeed())

		Expect(rdh.Specs).To(SatisfyAll(
			HaveLen(1),
			HaveKeyWithValue("id-1", SatisfyAll(
				HaveField("Kind", Equal("demo")),
				HaveField("MutableSpecDefinition.CreationSpecDefinition.Specs", Equal("val")),
				HaveField("MutableSpecDefinition.NextReconcile", Not(BeNil())),
			)),
		))
	})

	It("Fails when double id", func() {

		rdh.Specs["id-1"] = core.SpecDefinition{
			Kind: "demo",
		}

		Expect(rdh.CreateSpec("id-1", "demo", core.CreationSpecDefinition{
			Specs: "val",
		})).ToNot(Succeed())
	})
})

var _ = Describe("RDH - Get Spec", func() {

	var rdh core.ResourceDefinitionHandler

	BeforeEach(func() {
		rdh = *core.NewResourceDefinitionHandler()
	})

	It("Gets a spec", func() {

		rdh.Specs["id-1"] = core.SpecDefinition{
			Kind: "demo",
		}

		spec, err := rdh.GetSpec("id-1")
		Expect(err).To(BeNil())

		Expect(spec.Kind).To(Equal("demo"))
	})

	It("Fails when the spec does not exists", func() {
		_, err := rdh.GetSpec("id-1")
		Expect(err).ToNot(BeNil())
	})
})

var _ = Describe("RDH - Set Spec", func() {

	var rdh core.ResourceDefinitionHandler

	BeforeEach(func() {
		rdh = *core.NewResourceDefinitionHandler()
	})

	It("Sets a spec", func() {

		rdh.Specs["id-1"] = core.SpecDefinition{
			Kind: "demo",
		}

		Expect(rdh.SetSpec("id-1", core.MutableSpecDefinition{
			CreationSpecDefinition: core.CreationSpecDefinition{
				Specs: "new_val",
			},
		})).To(Succeed())

		Expect(rdh.Specs).To(SatisfyAll(
			HaveLen(1),
			HaveKeyWithValue("id-1", SatisfyAll(
				HaveField("Kind", Equal("demo")),
				HaveField("MutableSpecDefinition.CreationSpecDefinition.Specs", Equal("new_val")),
			)),
		))
	})

	It("Fails when spec does not exists", func() {
		Expect(rdh.SetSpec("id-1", core.MutableSpecDefinition{})).ToNot(Succeed())
	})
})

var _ = Describe("RDH - Get Next Spec", func() {

	var rdh core.ResourceDefinitionHandler

	BeforeEach(func() {
		rdh = *core.NewResourceDefinitionHandler()
	})

	It("The next spec is the only spec", func() {

		nextRec := time.Now().Add(-2 * time.Second)

		rdh.Specs["id-1"] = core.SpecDefinition{
			CreatedAt: time.Now().Add(-5 * time.Second),
			MutableSpecDefinition: core.MutableSpecDefinition{
				NextReconcile: &nextRec,
			},
		}

		id, _, err := rdh.GetNext()
		Expect(err).To(BeNil())
		Expect(id).To(Equal("id-1"))
	})

	It("The next spec is the latest spec", func() {

		nextRec := time.Now().Add(-2 * time.Second)
		beforeNext := nextRec.Add(-1 * time.Second)

		rdh.Specs["id-1"] = core.SpecDefinition{
			CreatedAt: time.Now().Add(-5 * time.Second),
			MutableSpecDefinition: core.MutableSpecDefinition{
				NextReconcile: &nextRec,
			},
		}
		rdh.Specs["id-2"] = core.SpecDefinition{
			CreatedAt: time.Now().Add(-5 * time.Second),
			MutableSpecDefinition: core.MutableSpecDefinition{
				NextReconcile: &beforeNext,
			},
		}

		id, _, err := rdh.GetNext()
		Expect(err).To(BeNil())
		Expect(id).To(Equal("id-2"))
	})

	It("No item needs a reconcile", func() {

		nextRec := time.Now().Add(5 * time.Second)

		rdh.Specs["id-1"] = core.SpecDefinition{
			CreatedAt: time.Now().Add(-5 * time.Second),
			MutableSpecDefinition: core.MutableSpecDefinition{
				NextReconcile: &nextRec,
			},
		}

		_, _, err := rdh.GetNext()
		Expect(err).ToNot(BeNil())
	})

	It("There is no spec", func() {
		_, _, err := rdh.GetNext()
		Expect(err).ToNot(BeNil())
	})
})
