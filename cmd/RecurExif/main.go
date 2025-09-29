package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/varppi/RecurExif/internal/parser"
	"github.com/varppi/RecurExif/internal/progress"
)

func info(msg string) {
	if !viper.GetBool("quiet") {
		fmt.Fprintln(os.Stderr, msg)
	}
}

func main() {
	fmt.Fprint(os.Stderr, `
______________________________  ______________________  __________________
___  __ \__  ____/_  ____/_  / / /__  __ \__  ____/_  |/ /___  _/__  ____/
__  /_/ /_  __/  _  /    _  / / /__  /_/ /_  __/  __    / __  / __  /_    
_  _, _/_  /___  / /___  / /_/ / _  _, _/_  /___  _    | __/ /  _  __/    
/_/ |_| /_____/  \____/  \____/  /_/ |_| /_____/  /_/|_| /___/  /_/       

`)
	pflag.String("extensions", "", "Only parses files with specific extensions. Example: pdf,png,jpg")
	pflag.String("search", "", "Finds all files that contain metadata containing the search keyword")
	pflag.Bool("no-progress", false, "Disables progress bar")
	pflag.Bool("quiet", false, "Doesn't show any extra info")
	pflag.Bool("json", false, "Outputs everything in JSON")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	args := pflag.Args()
	if len(args) == 0 {
		info("target not specified. please provide target directory/file with: ./program <path>")
		os.Exit(1)
	}
	viper.Set("target", args[len(args)-1])

	target := viper.GetString("target")
	targetInfo, err := os.Stat(target)
	if err != nil {
		info("couldn't open file/dir")
		os.Exit(1)
	}

	fastLookupTable := make(map[string]bool)
	for _, extension := range strings.Split(viper.GetString("extensions"), ",") {
		fastLookupTable[strings.TrimSpace(extension)] = true
	}
	var files []string
	if targetInfo.IsDir() {
		filepath.Walk(target, func(path string, info fs.FileInfo, err error) error {
			fmt.Fprintf(os.Stderr, "searching... %d files found\r", len(files))
			if info.IsDir() {
				return nil
			}
			if viper.GetString("extensions") != "" {
				ext := strings.Split(path, ".")
				if _, ok := fastLookupTable[ext[len(ext)-1]]; !ok {
					return nil
				}
			}
			files = append(files, path)
			return nil
		})
		fmt.Fprintln(os.Stderr, strings.Repeat(" ", 50))
	} else {
		files = []string{target}
	}

	info("starting...")
	progressBar := &progress.ProgressBar{Max: len(files), Disabled: viper.GetBool("no-progress")}
	progressBar.StartLoading()
	for _, file := range files {
		progressBar.Progress()
		metadata, err := parser.ParseFileMetadata(file)
		if err != nil {
			info(fmt.Sprintf("%s: %s", file, err.Error()))
			continue
		}

		var combined string
		if metadata.Exif != nil {
			combined += fmt.Sprintf("EXIF: %s\n", metadata.ExifJson)
		}
		if metadata.Pdf != nil {
			combined += fmt.Sprintf("PDF: %s\n", metadata.PdfJson)
		}
		cond1 := strings.Contains(strings.ToLower(combined), strings.ToLower(viper.GetString("search")))
		cond2 := metadata.Exif != nil || metadata.Pdf != nil
		if cond1 && cond2 && !viper.GetBool("json") {
			fmt.Printf("%s\n%s:\n%s", strings.Repeat(" ", 75), file, combined)
		}
		if cond1 && cond2 && viper.GetBool("json") {
			jsonVersion, err := json.Marshal(struct {
				File string            `json:"file"`
				Exif map[string]any    `json:"exif,omitempty"`
				Pdf  map[string]string `json:"pdf,omitempty"`
			}{
				File: file,
				Exif: metadata.Exif,
				Pdf:  metadata.Pdf,
			})
			if err != nil {
				info(err.Error())
				os.Exit(1)
			}
			fmt.Println(string(jsonVersion))
		}
	}
	progressBar.StopLoading()
	info("all files processed!")
}
