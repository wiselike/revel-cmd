package main_test

import (
	"os"
	"testing"

	"github.com/wiselike/revel-cmd/model"
	main "github.com/wiselike/revel-cmd/revel"
	"github.com/stretchr/testify/assert"
)

// test the commands.
func TestRevelTest(t *testing.T) {
	a := assert.New(t)
	gopath := setup("revel-test-test", a)

	t.Run("Test", func(t *testing.T) {
		a := assert.New(t)
		c := newApp("test-test", model.NEW, nil, a)
		a.Nil(main.Commands[model.NEW].RunWith(c), "Failed to run test-test")
		c.Index = model.TEST
		c.Test.ImportPath = c.ImportPath
		a.Nil(main.Commands[model.TEST].RunWith(c), "Failed to run test-test")
	})

	if !t.Failed() {
		if err := os.RemoveAll(gopath); err != nil {
			a.Fail("Failed to remove test path")
		}
	}
}
