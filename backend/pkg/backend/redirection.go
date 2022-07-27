package backend

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Infoblox-CTO/heka-ui/utils"
	"github.com/fsnotify/fsnotify"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// HekaOptin expected to be true by default as in where heka ui container gets started
var HekaOptIn = true

func TriggerRedirection(timeout int, script string, args []string) bool {

	var (
		statusCode int
		response   []byte
		err        error
		retryFlag  = false
		count      = 0
	)

	for {
		statusCode, response, err = utils.MakePostRequest(60, script, args)

		if err != nil || statusCode != http.StatusOK {
			logger.Errorf(" redirection HTTP POST request failed : %v with statusCode: %v", err, statusCode)
			if err != nil {
				logger.Errorf(" redirection HTTP POST request failed with err: %v", err)
			} else {
				logger.Errorf(" redirection HTTP POST request failed with Response code : %v", statusCode)
			}
			if retryFlag {
				logger.Errorf("Retry failed to execute HTTP redirect execution %v", err)
				return false
			}
			time.Sleep(30 * time.Second)
			retryFlag = true
			if count > 5 {
				return false
			}
			count++
			continue

		}
		break
	}

	type Err struct {
		Error  string `json:"error"`
		Status string `json:"status"`
	}
	var check Err

	err = json.Unmarshal(response, &check)
	if err != nil {
		logger.Errorf("Error while Unmarshalling redirection HTTP Response: %v", err)
		return false
	}

	if check.Status != "SUCCESS" {
		logger.Errorf("Failed to execute HTTP rediection in nios confirguration: %v : %v", check.Status, check.Error)
		return false
	}
	logger.Infoln("Successfully executed HTTP redirection configuration in nios", check.Status)
	return true

}

func HandleHekaOptin() {
	go func() {
		for {
			handleConfigChange(10 * time.Minute)
			time.Sleep(2 * time.Second)
		}
	}()

}

func handleConfigChange(duration time.Duration) {
	// waits for heka optin flag change
	// or waits for the specified duration
	timer := time.NewTimer(duration)
	eventChan := make(chan bool)
	defer timer.Stop()
	defer close(eventChan)
	for {
		viper.OnConfigChange(func(event fsnotify.Event) {
			defer viper.OnConfigChange(nil)
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
				logger.Infof("config file has been changed, operation: %v, filename: %s", event.Op, event.Name)
				select {
				case _, ok := <-eventChan:
					if !ok {
						logger.Infoln("file change event channel closed")
					}
				case eventChan <- true:
				}
			}
		})
		select {
		case <-timer.C:
			logger.Infoln("Returned from timer")
			return

		case <-eventChan:
			if isHekaOptInChanged() {
				logger.Infoln("heka optin changed")
				if !viper.GetBool("labels.heka_optin") {
					logger.Infoln("trigger remove redirection")
					TriggerRedirection(60, viper.GetString("remove.redirect.script"), []string{})
				} else {
					logger.Infoln("heka is opted in again trigger injection")
					// inject redirection again
					TriggerRedirection(60, viper.GetString("inject.redirect.script"), []string{})
				}
				return
			}
			logger.Infoln("heka not changed")
		}
	}
}

func ReadConfigFile(filename string) error {
	viper.SetConfigFile(filename)
	err := viper.ReadInConfig()
	// Listen for any updates done to the json file irrespective of the result
	// as the file can be updated later
	viper.WatchConfig()
	if err != nil {
		return err
	}
	return nil
}

func isHekaOptInChanged() bool {
	hFlag := getHekaOptIn()
	if HekaOptIn != hFlag {
		HekaOptIn = hFlag
		return true
	}
	return false
}

func getHekaOptIn() bool {
	labelsMap := viper.GetStringMap("labels")
	if len(labelsMap) > 0 {
		val, ok := labelsMap["heka_optin"]
		if ok {
			return val.(bool)
		}
	}
	return false
}
