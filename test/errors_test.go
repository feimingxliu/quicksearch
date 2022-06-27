package test

import (
	"fmt"
	"github.com/pkg/errors"
	"testing"
)

func TestErrors(t *testing.T) {
	err := fmt.Errorf("err occurs in main")
	fmt.Printf("%+v\n", errors.WithStack(err))
	a()
}

func a() {
	err := fmt.Errorf("err occurs in a")
	fmt.Printf("%+v\n", errors.WithStack(err))
	b()
}

func b() {
	err := fmt.Errorf("err occurs in b")
	fmt.Printf("%+v\n", errors.WithStack(err))
	c()
}

func c() {
	err := fmt.Errorf("err occurs in c")
	fmt.Printf("%+v\n", errors.WithStack(err))
}
