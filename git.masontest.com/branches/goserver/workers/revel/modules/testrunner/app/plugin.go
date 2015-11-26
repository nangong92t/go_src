package app

import (
	"fmt"
	"git.masontest.com/branches/goserver/workers/revel"
)

func init() {
	revel.OnAppStart(func() {
		fmt.Println("Go to /@tests to run the tests.")
	})
}
