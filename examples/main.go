package examples

import (
	"dependency"
)

func main() {
	dependency := dependency.NewDependency()
	if err := dependency.Get(); err != nil {
		panic(err)
	}
}
