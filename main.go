package main

import (
	"fmt"

	"github.com/l3lackShark/gosumemory/patterns"
)

func main() {
	fmt.Println("owo")
	// for i := 0; i < 500; i++ {
	// 	fmt.Println(patterns.ResolveOsuStatus())
	// }
	err := patterns.InitBase()
	if err != nil {
		fmt.Println("Error has occured! ", err)
	}
	fmt.Println(patterns.OsuStaticAddresses)

}
