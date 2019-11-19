package conf

import (
	"github.com/zhanbei/static-server/helpers/terminator"
	"github.com/zhanbei/static-server/utils"
)

// In strict mode, configures will be validated saturantly, like:
// - ignoring the #Enabled flag;
// - validate the logical warnings;
func strictMode() bool {
	return true
}

func doPass(enabled bool) bool {
	// Pass the validation if disabled & in the loose mode.
	return !strictMode() && !enabled
}

type Configure struct {
	// The www-root-dir path.
	RootDir string `json:"rootDir"`
	// The address or port for the server.
	Address string `json:"address"`

	Server *ServerOptions `json:"server"`

	Loggers *OptionLoggers `json:"loggers"`

	MongoDbOptions *MongoDbOptions `json:"mongo"`

	GorillaOptions *OptionLoggerGorilla `json:"mongo"`
}

var has = utils.NotEmpty
var exist = utils.NotEmpty

// Validate the syntax.
func (m *Configure) IsValid() bool {
	if !has(m.RootDir) {
		terminator.ExitWithConfigError(nil, "Please specify an address( or at least a port) in your configuration file!")
		return false
	}
	m.Address, _ = ValidateArgAddressOrExit(m.Address)
	m.RootDir = ValidateArgRootDirOrExit(m.RootDir)

	if m.Server == nil {
		m.Server = NewDefaultServerOptions()
	}

	if m.Loggers != nil {
		for _, logger := range *m.Loggers {
			if !logger.IsValid() {
				return false
			}
		}
	}

	if m.MongoDbOptions != nil && !m.MongoDbOptions.IsValid() {
		return false
	}

	if m.GorillaOptions != nil && !m.GorillaOptions.IsValid() {
		return false
	}

	return m.Server.IsValid()
}

// Validate the required resources.
func (m *Configure) ValidateFile() error {

	return nil
}
