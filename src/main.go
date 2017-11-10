package main

import (
    "flag"
    "fmt"
    core "searchdoc/src/core"
)

func main() {

    // get the arguments
    var (
        query    string
        language string
    )

    flag.StringVar(&query, "query", "", "The query to search")
    flag.StringVar(&language, "language", "", "The query to search")

    download_list := flag.Bool("download_list", false, "Indicate if you want to list the downloadable packages.")
    installed_list := flag.Bool("list", false, "List all installed packages.")
    package_to_download := flag.String("download", "", "Download the specified package.")

    flag.Parse()

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

    fmt.Printf("language: %s\nquery: %s\n", language, query)
    // process the command
    // TODO: replace with connection to ui
    data, err := core.Query(query, language)
    if err != nil {
        fmt.Printf("Error: %v\n", err.Error())
        return
    }

    fmt.Print(data)
}
