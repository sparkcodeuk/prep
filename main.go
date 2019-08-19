package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"text/template"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// Handle the -force argument if file exists
func writeForceCheck(path string, force bool) {
	if _, err := os.Stat(path); !os.IsNotExist(err) && !force {
		fmt.Printf("Error: file exists [%s], use -force to override\n", path)
		os.Exit(1)
	}
}

func main() {
	// Command line argument parsing
	args := ParseArgs()

	// Init input/output
	inputFileRegexp := regexp.MustCompile(`^(.+)\.tmpl$`)
	inputFileMatches := inputFileRegexp.FindStringSubmatch(path.Base(args.InputFile))

	if len(inputFileMatches) != 2 {
		log.Fatalln("Invalid input filename, expecting a go template with a .tmpl file extension (e.g., myfile.conf.tmpl)")
	}

	var inputFileData []byte
	var err error

	inputFileData, err = ioutil.ReadFile(args.InputFile)
	if err != nil {
		log.Fatal(err)
	}

	var outputFile *os.File
	defer outputFile.Close()

	if args.WriteFile != "" {
		writeForceCheck(args.WriteFile, args.Force)

		outputFile, err = os.Create(args.WriteFile)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		outputFile = os.Stdout
	}

	// Perform pre-check if requested
	if args.PreCheck {
		var tmplParsedData bytes.Buffer

		// A simple parse of the template
		// (to see if existing content manifests as go template syntax errors)
		tmpl, err := template.New("root").Parse(string(inputFileData))
		if err != nil {
			fmt.Println("Failed parsing your file as a go template. This failure points to content being accidentally interpreted as go template code, please fix before editing your file:")
			log.Fatal(err)
		}

		err = tmpl.Execute(&tmplParsedData, "")
		if err != nil {
			fmt.Println("Failed executing your file as a go template. This failure points to content being accidentally interpreted as go template code, please fix before editing your file:")
			log.Fatal(err)
		}

		// Diff the parsed template against the file to check for differences
		if string(inputFileData) != tmplParsedData.String() {
			fmt.Println("The following difference(s) were detected in your file when parsed as a go template. These failure(s) point to content being accidentally interpreted as go template code, please fix before editing your file:")
			dmp := diffmatchpatch.New()
			diffs := dmp.DiffMain(string(inputFileData), tmplParsedData.String(), false)
			fmt.Println(dmp.DiffPrettyText(diffs))

			os.Exit(1)
		} else {
			fmt.Println("File passed; no issues detected.")
		}

		os.Exit(0)
	}

	// Gather var data and merge in order declared
	var tmplValues = make(VarValues)

	varIndex := 0
	varFileIndex := 0
	skipNext := false
	for _, v := range os.Args[1:] {
		if skipNext {
			skipNext = false
			continue
		}

		switch v {
		case "-var":
			tmplValues.Merge(args.Vars[varIndex].Values())
			varIndex++
			skipNext = true
		case "-var-file":
			tmplValues.Merge(args.VarFiles[varFileIndex].Values())
			varFileIndex++
			skipNext = true
		}
	}

	// Process template
	tmpl, err := template.New("root").Parse(string(inputFileData))
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(outputFile, &tmplValues)
	if err != nil {
		log.Fatal(err)
	}
}
