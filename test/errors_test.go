package test

import (
	"fmt"
	"github.com/feimingxliu/quicksearch/pkg/errors"
	"testing"
)

func TestErrors(t *testing.T) {
	c := func() {
		err := fmt.Errorf("err occurs in c")
		fmt.Println(errors.WithStack(err))
	}
	b := func() {
		err := fmt.Errorf("err occurs in b")
		fmt.Println(errors.WithStack(err))
		c()
	}
	a := func() {
		err := fmt.Errorf("err occurs in a")
		fmt.Println(errors.WithStack(err))
		b()
	}
	err := fmt.Errorf("err occurs in main")
	fmt.Println(errors.WithStack(err))
	a()
}
