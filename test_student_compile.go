package main

import (
	"fmt"
	"go/build"
	"log"
)

func main() {
	fmt.Println("Testing student module compilation...")
	
	// Test if the student module package can be imported
	pkg, err := build.Import("academictoken/x/student", ".", build.FindOnly)
	if err != nil {
		log.Printf("Error importing student module: %v", err)
		return
	}
	
	fmt.Printf("Student module found at: %s\n", pkg.Dir)
	fmt.Println("Module should be ready to compile!")
}
