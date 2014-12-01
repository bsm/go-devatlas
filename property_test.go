package devatlas

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Property", func() {

	It("should decode from JSON", func() {
		var ps []Property
		Expect(json.Unmarshal([]byte(`["iheight", "smodel", "bisOld"]`), &ps)).NotTo(HaveOccurred())
		Expect(ps).To(Equal([]Property{
			Property{KIND_INT, "height"},
			Property{KIND_STRING, "model"},
			Property{KIND_BOOL, "isOld"},
		}))
	})

})
