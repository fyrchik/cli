package cli

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestFlag_Set(t *testing.T) {
	const (
		combine = iota
		parse
		parseMany
		postValidate
		validate
	)
	var (
		f       = &Flag{}
		invoked = map[int]struct{}{}
		g       = NewGomegaWithT(t)
	)

	f = f.
		SetCombine(func(a, b interface{}) interface{} {
			invoked[combine] = struct{}{}
			return nil
		}).
		SetDefault(123).
		SetParse(func(arg string) (interface{}, error) {
			invoked[parse] = struct{}{}
			return nil, nil
		}).
		SetParseMany(func(args []string) (interface{}, int, error) {
			invoked[parseMany] = struct{}{}
			return nil, 0, nil
		}).
		SetValidate(func(v interface{}) error {
			invoked[validate] = struct{}{}
			return nil
		}).
		SetPostValidate(func(v interface{}) error {
			invoked[postValidate] = struct{}{}
			return nil
		})

	g.Expect(f.Combine).NotTo(BeNil())
	f.Combine(nil, nil)
	g.Expect(invoked).To(HaveKey(combine))

	g.Expect(f.Default).To(Equal(123))

	g.Expect(f.Parse).NotTo(BeNil())
	f.Parse("")
	g.Expect(invoked).To(HaveKey(parse))

	g.Expect(f.ParseMany).NotTo(BeNil())
	f.ParseMany(nil)
	g.Expect(invoked).To(HaveKey(parseMany))

	g.Expect(f.PostValidate).NotTo(BeNil())
	f.PostValidate(nil)
	g.Expect(invoked).To(HaveKey(postValidate))

	g.Expect(f.Validate).NotTo(BeNil())
	f.Validate(nil)
	g.Expect(invoked).To(HaveKey(validate))
}
