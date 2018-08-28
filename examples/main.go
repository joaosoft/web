package examples

import (
	service "dependency/service"
)

func main() {
	dependency, err := service.NewDependency()
	if err != nil {
		panic(err)
	}

	if err := dependency.Get(); err != nil {
		panic(err)
	}
}
