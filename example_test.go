package sen_test

import (
	"fmt"

	"github.com/bongnv/sen"
)

func Example() {
	err := sen.Run()
	fmt.Println(err)

	// Output:
	// <nil>
}
