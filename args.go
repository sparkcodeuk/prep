package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type VarValues map[interface{}]interface{}

func (dest *VarValues) Merge(src VarValues) {
	for k, v := range src {
		(*dest)[k] = v
	}
}

type VarValuesable interface {
	Values() VarValues
}

type Var struct {
	Name  string
	Value string
}

func (v Var) Values() VarValues {
	var vv = make(VarValues)

	vv[v.Name] = v.Value

	return vv
}

type Vars []Var

func (v *Vars) Set(value string) error {
	parts := strings.SplitN(value, "=", 2)
	if len(parts) != 2 {
		return errors.New("argument should be in the form: name=value")
	}

	varName := strings.TrimSpace(parts[0])
	if len(varName) == 0 {
		return errors.New("no name defined for template variable")
	}

	*v = append(*v, Var{Name: varName, Value: parts[1]})

	return nil
}

func (v *Vars) String() string {
	var s []string
	for _, i := range *v {
		s = append(s, fmt.Sprintf("%v=%v", i.Name, i.Value))
	}

	return strings.Join(s, ", ")
}

type VarFile string

func (v VarFile) Values() VarValues {
	var vv VarValues

	data, err := ioutil.ReadFile(string(v))
	if err != nil {
		log.Fatal(err)
	}

	yaml.Unmarshal(data, &vv)

	return vv
}

type VarFiles []VarFile

func (v *VarFiles) Set(value string) error {
	if value == "" {
		return errors.New("no filename defined")
	}

	varFilesRegexp := regexp.MustCompile(`^(.+)\.y(a)?ml$`)

	if !varFilesRegexp.Match([]byte(value)) {
		return errors.New("only YAML files are supported")
	}

	*v = append(*v, VarFile(value))

	return nil
}

func (v *VarFiles) String() string {
	var s []string
	for _, i := range *v {
		s = append(s, string(i))
	}

	return strings.Join(s, ", ")
}

type Args struct {
	Help      bool
	Version   bool
	Force     bool
	PreCheck  bool
	InputFile string
	WriteFile string
	Vars      Vars
	VarFiles  VarFiles
}

func fatalError(message string) {
	fmt.Println("Error:", message)
	fmt.Printf("\n---\n\n")
	flag.Usage()
	os.Exit(1)
}

func ParseArgs() Args {
	args := Args{}

	flag.BoolVar(&args.Help, "help", false, "Display this help message")
	flag.BoolVar(&args.Version, "version", false, "Print version and exit")
	flag.BoolVar(&args.Force, "force", false, "Force file overwrite")
	flag.BoolVar(&args.PreCheck, "pre-check", false, "Check a file for issues prior to editing as a go template")
	flag.Var(&args.Vars, "var", "Define a template variable (e.g., -var name=value)")
	flag.Var(&args.VarFiles, "var-file", "Import template variables from a file (e.g., -var-file vars.yaml)")
	flag.StringVar(&args.WriteFile, "write-file", "", "Output the interpolated template to a file")

	// Custom usage message
	flag.Usage = func() {
		PrintVersion()
		fmt.Println(`
Usage:-
  prep [args] <input file>`)
		fmt.Println()
		flag.PrintDefaults()
		fmt.Println(`
Command examples:-

Processes "my_file.conf.tmpl", using variables defined in the variables file
"my_vars.yaml" and defining/overriding the variable "username" and outputting
the result to the terminal:
  prep -var-file my_vars.yaml -var username=jbloggs my_file.conf.tmpl

... the same command, but writes specifically to the named file, "my_file.conf":
  prep -var-file my_vars.yaml -var username=jbloggs -write-file my_file.conf -force my_file.conf.tmpl

A useful command to check for any go templating related issues before you start
adding any go template syntax:
  prep -pre-check my_file.conf.tmpl`)
		fmt.Println()
	}

	flag.Parse()

	if args.Help {
		flag.Usage()
		os.Exit(0)
	}

	if args.Version {
		PrintVersion()
		os.Exit(0)
	}

	// Parse positional arguments
	switch flag.NArg() {
	case 0:
		fatalError("no input file specified")
	case 1:
		args.InputFile = flag.Arg(0)
	default:
		fatalError("this command takes a maximum of one input file")
	}

	return args
}
