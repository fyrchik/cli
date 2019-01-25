package cli

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestContext_Parse(t *testing.T) {
	var (
		c   *Context
		err error
		g   = NewGomegaWithT(t)
	)

	c = NewContext()
	err = c.AddFlags(
		BoolFlag("bf", "-b", "--enable-some-shit"),
		IntFlag("count", "--count"),
		StringFlag("strf", "-s", "--test_flag"),
		StringSliceFlag("multi", "-m", "-ms", "--multi-string"),
	)
	g.Expect(err).NotTo(HaveOccurred())

	err = c.Parse([]string{"-s", "test_value"})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(c.Named).To(HaveKeyWithValue("strf", "test_value"))

	c.clear()
	err = c.Parse([]string{"--test_flag"})
	g.Expect(err).NotTo(BeNil())

	c.clear()
	err = c.Parse([]string{"--test_flag", "value", "--count", "123", "--test_flag", "value2"})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(c.Named).To(HaveKeyWithValue("strf", "value2"))
	g.Expect(c.Named).To(HaveKeyWithValue("count", int64(123)))

	c.clear()
	err = c.Parse([]string{"-b"})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(c.Named).To(HaveKeyWithValue("bf", true))

	c.clear()
	err = c.Parse([]string{"-b", "-s", "123"})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(c.Named).To(HaveKeyWithValue("bf", true))
	g.Expect(c.Named).To(HaveKeyWithValue("strf", "123"))

	c.clear()
	err = c.Parse([]string{"-m", "first", "-s", "string", "--multi-string", "second", "-ms", "third"})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(c.Named).To(HaveKeyWithValue("strf", "string"))
	g.Expect(c.Named).To(HaveKeyWithValue("multi", []string{"first", "second", "third"}))
}
