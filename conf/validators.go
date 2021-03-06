package conf

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/zhanbei/static-server/helpers/terminator"
)

func ValidateArgAddressOrExit(address string) (string, int) {
	port, err := strconv.Atoi(address)
	if err != nil {
		return address, -1
		// Check the address.
	}
	// The address is only a port.
	if port < 1 || 65535 < port {
		terminator.ExitWithConfigError(nil, "ERROR: unavailable port["+strconv.Itoa(port)+"]; make sure http port is number and is limited to <0-65535>.")
	}
	if port <= 1024 {
		fmt.Println("WARNING: the port[" + strconv.Itoa(port) + "] specified is not bigger than 1024; root privileges may be needed!")
	}
	address = ":" + strconv.Itoa(port)
	return address, port
}

func ValidateArgRootDirOrExit(rootDir string) string {
	rootDir, err := filepath.Abs(rootDir)
	if err != nil {
		fmt.Println("ERROR: The specified www-root-directory is invalid:" + rootDir)
		terminator.ExitWithConfigError(err, "ERROR: The specified www-root-directory is invalid: "+rootDir)
	}
	return rootDir
}
