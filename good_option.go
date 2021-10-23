package good

import "github.com/xqk/good/pkg/application"

type Option = application.Option

type Disable = application.Disable

const (
	DisableParserFlag      Disable = application.DisableParserFlag
	DisableLoadConfig      Disable = application.DisableLoadConfig
	DisableDefaultGovernor Disable = application.DisableDefaultGovernor
)

var WithConfigParser = application.WithConfigParser
var WithDisable = application.WithDisable
