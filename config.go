package weso

import (
	"flag"
	"io"
)

// Config is configuration file for CLI.
type Config struct {
	URL      string
	Origin   string
	Template string
	Debug    bool
}

// NewConfig create Config and output error to `errOut`.
func NewConfig(args []string, errOut io.Writer) (*Config, bool) {
	flags := flag.NewFlagSet("weso", flag.ContinueOnError)
	template := flags.String("template", "", "")
	origin := flags.String("origin", "", "")
	debug := flags.Bool("debug", false, "")

	usage := `Usage: weso [options] url

options:
	-template path
		path to template file.
	-origin url
		oritin url for websocket header.
	-debug
		enable debug mode (print expanded template).
`

	flags.Usage = func() {
		errOut.Write([]byte(usage))
	}

	if err := flags.Parse(args); err != nil {
		return nil, false
	}

	if len(flags.Args()) != 1 {
		flags.Usage()
		return nil, false
	}

	c := &Config{
		flags.Args()[0],
		*origin,
		*template,
		*debug,
	}

	return c, true
}
