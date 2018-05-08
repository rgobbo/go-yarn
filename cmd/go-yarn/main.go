package main

import (
	"flag"
	"fmt"
	"log"

	goyarn "github.com/rgobbo/go-yarn"
)

func main() {
	var showHelp bool
	flag.BoolVar(&showHelp, "h", false, "show help")
	var configFile string
	var destPath string

	flag.StringVar(&configFile, "c", "yarn.json", "yarn.json configuration file")
	flag.StringVar(&destPath, "f", "./vendor", "destination path")

	flag.Parse()

	usage := `
Install dependencies from a yarn.json file into a given path
example
  go-yarn -c yarn.json -f ./vendor
`
	// if no args or -h flag
	// print usage and return

	if showHelp  {
		fmt.Println(usage)
		return
	}

	err := goyarn.YarnInstall(configFile, destPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Libraries downloaded !!!")

}