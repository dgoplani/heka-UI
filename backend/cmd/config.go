package main

import (
	"time"

	flag "github.com/spf13/pflag"
)

const (
	defaultHekaUiDir          = "/etc/onprem.d/scripts/"
	defaultParamFile          = "/etc/onprem.d/params.json"
	defaultLogLevel           = "info"
	defaultLogFormat          = "json"
	defaultNiosAddress        = "0.0.0.0"
	defaultNiosPort           = 999
	defaultInjectScript       = "inject_http_redirection"
	defaultInjectRemoveScript = "remove_http_redirection"
	defaultServerAddress      = "127.0.0.1:26750"
	defaultWapiVersion        = "v2.12"
	defaultHotfixInterval     = 5 * time.Minute
	defaultGridUuid           = ""
	defaultNiosVersion        = ""
)

var (
	flagLogLevel      = flag.String("log.level", defaultLogLevel, "set the logging level: debug, info, warning, error, fatal, panic")
	flagLogFormat     = flag.String("log.format", defaultLogFormat, "set the logging format: json or text")
	flagHekaDir       = flag.String("hekaui.dir", defaultHekaUiDir, "Heka directory where nios files are stored")
	flagParamFile     = flag.String("param.file", defaultParamFile, "params configuration file to get auth")
	flagcniosAddress  = flag.String("cnios.address", defaultNiosAddress, "cnios server address")
	flagcniosPort     = flag.Int("cnios.port", defaultNiosPort, "cnios server port")
	flagInjectscript  = flag.String("inject.redirect.script", defaultInjectScript, "params configuration file to get auth")
	flagRemoveInject  = flag.String("remove.redirect.script", defaultInjectRemoveScript, "cnios server address")
	flagHealthAddress = flag.String("server.address", defaultServerAddress, "Address and port to bind the healthcheck listener to")
	flagWapiVersion   = flag.String("wapi.version", defaultWapiVersion, "Wapi version to be used to make NIOS APi call")
	flagPollInterval  = flag.Duration("hotfix.interval", defaultHotfixInterval, "sleep time untill that hotfix data will be collected")
	flagGridUUid      = flag.String("grid_uuid", defaultGridUuid, "Wapi version to be used to make NIOS APi call")
	flagNiosVersion   = flag.String("nios_version", defaultNiosVersion, "sleep time untill that hotfix data will be collected")
)
