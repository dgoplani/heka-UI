package main

import (
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/Infoblox-CTO/heka-ui/pkg/backend"
	"github.com/Infoblox-CTO/heka-ui/pkg/log"
	"github.com/Infoblox-CTO/heka-ui/utils"
	logger "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var wg sync.WaitGroup
var components = []Handle{ServerFactory()}

// Handle is the Interface used to manage heka sub components
// All Heka sub component will follow the algorithm defined using this interface
type Handle interface {
	// Name() returns the name of the component to be displayed in the logs
	Name() string

	// Init() takes care of initializing the component
	Init(wg *sync.WaitGroup) error

	// Start() runs the Manager thread of the component
	Start() error

	// Check() checks for the criterian to run the component in this node
	Check() bool

	// Stop() terminates the Manager thread and all Worker threads component gracefully
	Stop()

	// IsUP returns if the module is up and running
	IsUP() bool
}

var scripts = []string{"trigger_hotfix_manifest", "cloud_manifest.json", "inject_http_redirection", "remove_http_redirection", "nios_data_collection"}

func init() {
	flag.Parse()

	viper.BindPFlags(flag.CommandLine)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	log.Setup(
		viper.GetString("log.format"),
		viper.GetString("log.level"),
	)
	if err := backend.ReadConfigFile(viper.GetString("param.file")); err != nil {
		logger.Fatalf("Failed to read params into viper config", err)
	}
	optin := viper.GetBool("labels.heka_optin")
	if !optin {
		logger.Warn("Expected the heka optin value should be true ")
		backend.HekaOptIn = optin
	}

	hekaDir := viper.GetString("hekaui.dir")

	for _, script := range scripts {
		dst := path.Join(hekaDir, "/", script)
		logger.Debugf("script %v", dst)
		err := utils.CopyFile(script, dst)
		if err != nil {
			logger.Errorf("Failed to copy script:%v, Err %v", script, err)
			return
		}

		err = os.Chmod(dst, 0755)
		if err != nil {
			logger.Errorf("Failed to change file permissions of the script:%v, Err %v", script, err)
			return
		}
	}
}

func main() {
	logger.Infoln("Starting heka-ui backend service")

	time.Sleep(10 * time.Second)

	// HTTP injection is one time process, no need to add as component
	if viper.GetBool("labels.heka_optin") {
		if ok := backend.TriggerRedirection(60, viper.GetString("inject.redirect.script"), []string{}); !ok {
			logger.Errorln("No point on running heka-ui container if http redirection injection is failed")
			return
		}
	}

	backend.HandleHekaOptin()

	wg.Add(len(components))
	for _, c := range components {
		if c.Check() {
			err := c.Init(&wg)
			if err != nil {
				continue
			}
			err = c.Start()
			if err != nil {
				continue
			}
		} else {
			logger.Errorf("Mandatory check failed to start %v component", c.Name())
		}
	}
	startNginx()
	go utils.TriggerControldScript()
	wg.Wait()
}

func startNginx() {
	cmd := exec.Command("nginx", "-g", "daemon off;")
	cmd.Stdout = os.Stdout

	err := cmd.Start()

	if err != nil {
		logger.Fatalf("Failed to start nginx", err)
	}
	logger.Infoln("started nginx with pid", cmd.Process.Pid)
}
