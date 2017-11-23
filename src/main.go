package main

import (
	"flag"
	"fmt"
	"searchdoc/src/core"
	"searchdoc/src/ui"
)

func main() {

	// get the arguments
	var (
		language string
	)

	flag.StringVar(&language, "language", "", "The language to search")

	download_list := flag.Bool("download_list", false, "Indicate if you want to list the downloadable packages.")
	installed_list := flag.Bool("list", false, "List all installed packages.")
	package_to_download := flag.String("download", "", "Download the specified package.")

	flag.Parse()

	// handle non-query commands
	if *download_list {
		for _, x := range core.ListDownloadableDocsets() {
			fmt.Println(x)
		}
		return
	} else if *installed_list {
		for _, x := range core.ListInstalledDocsets() {
			fmt.Println(x)
		}
		return
	} else if *package_to_download != "" {
		core.DownloadDocset(*package_to_download)
		return
	}

	fmt.Printf("language: %s\n", language)

	ui.SetLanguage(language)
	ui.Init()
}
