package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

var scripts = []string{"testing_Nios", "trigger_hotfix_manifest", "inject_http_redirection", "remove_http_redirection", "collect_grid_data"}

type TestConf struct {
	LogLevel string `json:"log.level"`
}

// confOverride fetches log level configuration from file
func confOverride() (string, error) {
	var value TestConf
	jsonFile, err := os.Open(viper.GetString("test.conf.file"))
	if err != nil {
		// File is not present, so no conf override
		return "", err
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		err = fmt.Errorf("failed to read file : %v", err)
		fmt.Println(err)
		return "", err
	}
	err = json.Unmarshal(byteValue, &value)
	if err != nil {
		err = fmt.Errorf("unmarshall failed: %v", err)
		fmt.Println(err)
		return "", err
	}
	return value.LogLevel, nil
}

func init() {
	flag.Parse()

	viper.BindPFlags(flag.CommandLine)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// confOverride should be enabled only in dev environment
	logLevel, err := confOverride()
	if err != nil {
		logLevel = viper.GetString("log.level")
	}

	log.Setup(
		viper.GetString("log.format"),
		logLevel,
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
	go supervisorThreadManager()

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

// SupervisorThreadManager will ensure to bring up and bring down component manager as per mandatory checks
func supervisorThreadManager() {
	logger.Infoln("Started supervisor thread manager")

	for {
		// Supervisor checks the manager and worker thread for every 5 minutes
		time.Sleep(5 * time.Minute)
		logger.Infoln("Running supervisor thread manager  ")
		for _, c := range components {
			if !c.IsUP() && c.Check() {
				logger.Warningf("%v is down, starting the component", c.Name())
				c.Stop()
				time.Sleep(1 * time.Minute)
				wg.Add(1)
				err := c.Init(&wg)
				if err != nil {
					logger.Errorf("Failed to Initialize : %v component Manager, Error: %v", c.Name(), err)
					continue
				} else {
					logger.Infof("Successfully Initialized :%v component Manager...", c.Name())
				}
				err = c.Start()
				if err != nil {
					logger.Errorf("Failed to Start : %v component Manager, Error: %v", c.Name(), err)
					continue
				} else {
					logger.Warningf("Successfully Started: %v component Manager...", c.Name())
				}

			} else if !c.Check() && c.IsUP() {
				c.Stop()
				logger.Infof("Successfully Stopped %v component Manager due to failure of mandatory status", c.Name())
			}
		}
	}
}
