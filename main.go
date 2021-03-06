package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/foomo/gocontentful/erm"
)

const VERSION = "v0.9.5"

var Usage = func() {
	fmt.Printf("\nSYNOPSIS\n")
	fmt.Printf("     gocontentful -spaceid SpaceID -cmakey CMAKey [-contenttypes firsttype,secondtype...lasttype] path/to/target/package\n\n")
	flag.PrintDefaults()
	fmt.Printf("\nNote: The last segment of the path/to/target/package will be used as package name\n\n")
}

func usageError(comment string) {
	fmt.Println("ERROR:", comment)
	Usage()
	os.Exit(1)
}

func fatal(infos ...interface{}) {
	fmt.Println(infos...)
	os.Exit(1)
}

func main() {
	// Get parameters from cmd line flags
	flagSpaceID := flag.String("spaceid", "", "Contentful space ID")
	flagCMAKey := flag.String("cmakey", "", "Contentful CMA key")
	flagContentTypes := flag.String("contenttypes", "", "[Optional] Content type IDs to parse, comma separated")
	flagVersion := flag.Bool("version", false, "Print version and exit")
	flagHelp := flag.Bool("help", false, "Print version and exit")
	flag.Parse()

	if *flagVersion {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	if *flagHelp {
		Usage()
		os.Exit(0)
	}

	if *flagSpaceID == "" || *flagCMAKey == "" {
		usageError("Please specify the Contentful space ID and access Key")
	}

	if len(flag.Args()) != 1 {
		usageError("Missing arg path/to/target/package")
	}

	path := flag.Arg(0)
	packageName := filepath.Base(path)

	matched, err := regexp.MatchString(`[a-z].{2,}`, packageName)
	if !matched {
		usageError("Please specify the package name correctly (only small caps letters)")
	}

	fmt.Printf("Contentful API Generator starting...\n\n")

	var flagContentTypesSlice []string
	if *flagContentTypes != "" {
		for _, contentType := range strings.Split(*flagContentTypes, ",") {
			flagContentTypesSlice = append(flagContentTypesSlice, strings.TrimSpace(contentType))
		}
	}

	err = erm.GenerateAPI(filepath.Dir(path), packageName, *flagSpaceID, *flagCMAKey, flagContentTypesSlice)
	if err != nil {
		fatal("Something went horribly wrong...", err)
	}
	fmt.Println("ALL DONE!")

}
