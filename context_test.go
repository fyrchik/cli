package cli

import (
	"os"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

func TestContext_Parse(t *testing.T) {
	type (
		mapSS = map[string]string
		pair  = struct {
			key   string
			value string
		}
	)

	var (
		c   *Context
		err error
		g   = NewGomegaWithT(t)
	)

	c = NewContext()
	err = c.Parse([]string{"only", "positional"})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(c.Positional).To(Equal([]string{"only", "positional"}))

	err = c.AddFlags(
		BoolFlag("enable", "-e", "--enable"),
		BoolFlagT("disable", "-d", "--disable"),
		DurationFlag("duration", "-ti", "--time-interval"),
		IntFlag("count", "--count").
			SetEnvKey("ATTEMPTS_COUNT, COUNT"),
		StringFlag("strFlag", "-s", "--test_flag").
			SetEnvKey("STRING_FLAG"),
		StringFlag("type", "-t", "--type").
			SetValidate(OneOf("type1", "type2", "type3")),
		StringSliceFlag("multi", "-m", "-ms", "--multi-string"),
		StringSliceFlag("3pigs", "-p").
			SetPostValidate(func(v interface{}) error {
				if len(v.([]string)) < 3 {
					return errors.New("must have more than 3 elements")
				}
				return nil
			}),
		NewFlag("key-value", "-kv").
			SetParseMany(func(args []string) (interface{}, int, error) {
				if len(args) < 2 {
					return nil, 0, errors.New("must have more than 2 arguments")
				}
				return pair{args[0], args[1]}, 2, nil
			}).
			SetCombine(func(m, kv interface{}) interface{} {
				p := kv.(pair)
				if m == nil {
					return mapSS{p.key: p.value}
				}
				m.(mapSS)[p.key] = p.value
				return m
			}),
	)
	g.Expect(err).NotTo(HaveOccurred())

	err = c.Parse([]string{""})
	g.Expect(err).NotTo(HaveOccurred())

	err = c.AddFlags(IntFlag("enable", "-i"))
	g.Expect(err).To(HaveOccurred())

	err = c.AddFlags(IntFlag("time", "-t"))
	g.Expect(err).To(HaveOccurred())

	err = c.Parse([]string{"-s", "test_value"})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(c.Named).To(HaveKeyWithValue("strFlag", "test_value"))
	g.Expect(c.Named).To(HaveKeyWithValue("enable", false))
	g.Expect(c.Named).To(HaveKeyWithValue("disable", true))

	err = c.Parse([]string{"--test_flag"})
	g.Expect(err).NotTo(BeNil())

	err = os.Setenv("ATTEMPTS_COUNT", "nan")
	g.Expect(err).NotTo(HaveOccurred())
	err = os.Setenv("COUNT", "10")
	g.Expect(err).NotTo(HaveOccurred())
	err = c.Parse([]string{""})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(c.Named).To(HaveKeyWithValue("count", int64(10)))

	err = c.Parse([]string{"--test_flag", "value", "--count", "123", "--test_flag", "value2"})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(c.Named).To(HaveKeyWithValue("strFlag", "value2"))
	g.Expect(c.Named).To(HaveKeyWithValue("count", int64(123)))

	err = c.Parse([]string{"--count", "12s"})
	g.Expect(err).To(HaveOccurred())

	err = os.Setenv("STRING_FLAG", "value")
	g.Expect(err).NotTo(HaveOccurred())

	err = c.Parse([]string{"-e", "-d"})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(c.Named).To(HaveKeyWithValue("enable", true))
	g.Expect(c.Named).To(HaveKeyWithValue("disable", false))
	g.Expect(c.Named).To(HaveKeyWithValue("strFlag", "value"))

	err = c.Parse([]string{"-e", "-s", "123"})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(c.Named).To(HaveKeyWithValue("enable", true))
	g.Expect(c.Named).To(HaveKeyWithValue("strFlag", "123"))

	err = c.Parse([]string{"-m", "first", "-s", "string", "--multi-string", "second", "-ms", "third"})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(c.Named).To(HaveKeyWithValue("strFlag", "string"))
	g.Expect(c.Named).To(HaveKeyWithValue("multi", []string{"first", "second", "third"}))

	err = c.Parse([]string{"-ti", "2m2s"})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(c.Named).To(HaveKeyWithValue("duration", time.Duration(122*time.Second)))

	err = c.Parse([]string{"--type", "type123"})
	g.Expect(err).To(HaveOccurred())

	err = c.Parse([]string{"--type", "type2"})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(c.Named).To(HaveKeyWithValue("type", "type2"))

	err = c.Parse([]string{"-p", "nif-nif", "-p", "nuf-nuf"})
	g.Expect(err).To(HaveOccurred())

	err = c.Parse([]string{"-p", "nif-nif", "-p", "nuf-nuf", "-p", "naf-naf"})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(c.Named).To(HaveKeyWithValue("3pigs", []string{"nif-nif", "nuf-nuf", "naf-naf"}))

	err = c.Parse([]string{"-kv", "onlykey"})
	g.Expect(err).To(HaveOccurred())

	err = c.Parse([]string{"-kv", "mary", "girl", "-kv", "jonnie", "boy"})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(c.Named).To(HaveKeyWithValue(
		"key-value", map[string]string{"mary": "girl", "jonnie": "boy"}),
	)

}
