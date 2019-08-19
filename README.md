# prep - PRE-Processor

A simple file pre-processor using [go templates](https://golang.org/pkg/text/template/#pkg-overview) and optional [yaml](https://en.wikipedia.org/wiki/YAML) files.

## Usage

Processes `my_file.conf.tmpl`, using variables defined in the variables file `my_vars.yaml` and defining/overriding the variable `username` and outputting the result to the terminal:

`prep -var-file my_vars.yaml -var username=jbloggs my_file.conf.tmpl`

... the same command, but writes specifically to the named file, `my_file.conf`:

`prep -var-file my_vars.yaml -var username=jbloggs -write-file my_file.conf -force my_file.conf.tmpl`

A useful command to check for any go templating related issues before you start adding any go template syntax:

`prep -pre-check my_file.conf.tmpl`

The `-pre-check` argument interprets your input file as a go template with no variables and prints any syntax errors or differences in the output with that of the original template. Both checks point to parts of your being incorrectly interpreted as go template code which should be fixed before you being adding any go templating code.

## Installation

Either [download the latest release](https://github.com/sparkcodeuk/prep/releases) for your OS and stick the `prep` binary (making sure it's executable) in your `$PATH`, or build the release yourself.

To build the release, ensure [glide](https://glide.sh) is installed (required to manage the library dependencies). Clone the repo and run the `.../prep/util/build-release.sh` shell script. This will run some tests and then build for all supported OS's their respective binaries. You can then grab the one you need. (You can also run `go build`, but you'll forgo the tests that `build-release.sh` runs).

## Future improvements

* Additional template utility functions
* Add support for `-template name=/path/to/sub-template.tmpl` definitions
* .ini `-var-file` support
* Functional tests
