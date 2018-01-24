package main

import (
	"fmt"
	"flag"
	"os"
	mg "github.com/moris351/mergbigfolder"
	
)

var cp = flag.Bool("cp",false,"copy file ?") 

func main() {
     // Ensure the file argument is passed
    if len(os.Args) != 3 {
        fmt.Println("Please use this syntax: prog src dst",len(os.Args))
        return
    }

	if err:=mg.FindDiff(os.Args[1],os.Args[2]);err!=nil{
		fmt.Println("FindDiff err,", err)
	}


}

