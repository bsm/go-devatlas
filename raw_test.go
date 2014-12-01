package devatlas

import (
	"encoding/json"
	"regexp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("indexMap", func() {

	It("should unmarshal", func() {
		var v indexMap
		err := json.Unmarshal([]byte(`{"7":12314,"8":63662}`), &v)
		Expect(err).NotTo(HaveOccurred())
		Expect(v).To(Equal(indexMap{7: 12314, 8: 63662}))
	})

})

var _ = Describe("mixedSlice", func() {

	It("should unmarshal", func() {
		var v mixedSlice
		err := json.Unmarshal([]byte(`["a",1,2,7,"b",2538591,"x"]`), &v)
		Expect(err).NotTo(HaveOccurred())
		Expect(v).To(Equal(mixedSlice{
			"a", uint32(1), uint32(2), uint32(7), "b", uint32(2538591), "x",
		}))
	})

})

var _ = Describe("regexpSliceV", func() {

	It("should unmarshal", func() {
		var v regexpSliceV
		err := json.Unmarshal([]byte(`{"5":["A","B"], "6":["C","D"]}`), &v)
		Expect(err).NotTo(HaveOccurred())
		Expect(v).To(Equal(regexpSliceV{
			regexp.MustCompile(`C`),
			regexp.MustCompile(`D`),
		}))
	})

})

var _ = Describe("regexpSliceO", func() {

	It("should unmarshal", func() {
		var s regexpSliceO
		err := json.Unmarshal([]byte(`{"d":["A", "B", "C"], "6":{"1":"D"}}`), &s)
		Expect(err).NotTo(HaveOccurred())
		Expect(s).To(Equal(regexpSliceO{
			regexp.MustCompile(`A`),
			regexp.MustCompile(`D`),
			regexp.MustCompile(`C`),
		}))
	})

})
