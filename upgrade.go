package main

import (
    "fmt"
    "os"

    "github.com/jak3kaj/ePICmon/ePIC"
)

func main() {
	args := os.Args
    
    if ePIC.UpgradeFirmware(args[1], args[2]) {
        fmt.Printf("%s Upgraded Successfully\n", args[1])
    } else {
        fmt.Printf("%s Upgrade Failed\n", args[1])
    }
}
