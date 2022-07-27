package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

type authKeys struct {
	Jointoken     string `json:"jointoken"`
	WeakAuthToken string `json:"weak.auth.token"`
}

// checkStrongAuthKeys checks for api.key/jointoken in params.json file
func checkStrongAuthKeys(paramData authKeys) bool {
	return paramData.Jointoken != ""
}

/*CopyFile ... This function copys a file from src to dst*/
func CopyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()
	destination, err := os.Create(dst)
	if err != nil {
		return err
	}

	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}

	err = destination.Close()
	if err != nil {
		return err
	}

	return nil
}

//MakePostRequest will do HTTP POST request to connect to NIOS
func MakePostRequest(sec time.Duration, script string, args []string) (int, []byte, error) {
	url := "http://" + viper.GetString("cnios.address") + ":" + viper.GetString("cnios.port") + "/nios/api/v1.0/nios/1/exec"

	// create request
	input := struct {
		Name string   `json:"name"`
		Args []string `json:"args"`
	}{script, args}

	jsonStr, err := json.Marshal(input)
	if err != nil {
		logger.Errorf("Error while creating json input for API call: %v", err)
		return 0, nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		logger.Errorf("Error while creating the request: %v", err)

		return 0, nil, err
	}

	//set request headers
	req.Header.Add("Content-Type", "application/json")

	cli := http.Client{Timeout: sec * time.Second}
	//send request
	res, err := cli.Do(req)
	if err != nil {
		logger.Errorf("Request failed - %v", err)

		return 0, nil, err
	}
	defer res.Body.Close()

	resDump, _ := httputil.DumpResponse(res, true)
	logger.Debugf("Response dump - %s", string(resDump))

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Errorf("Failed to read POST heka_data status response body - %v", err)

	}
	return res.StatusCode, data, err
}

func TriggerControldScript() {
	for {
		collectHotfixData("--trigger")
		wait := viper.GetDuration("hotfix.interval")
		logger.Infoln("Collecing hotfix data from all nodes for every ", wait)
		time.Sleep(wait)
	}
}

func collectHotfixData(action string) {
	// Call a python script hotfix_manifest
	responseInfo := struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{}
	var (
		args []string
		// Duration in sec
		timeout time.Duration = 30
		script  string        = "trigger_hotfix_manifest"
	)
	args = append(args, action)

	code, response, err := MakePostRequest(timeout, script, args)

	if err != nil {
		logger.Errorf("NdAction HTTP POST request failed : %v ", err)
		return
	}

	if code != http.StatusOK {
		logger.Errorf("NdAction HTTP POST returned bad response : %v", code)
		return
	}

	err = json.Unmarshal(response, &responseInfo)
	if err != nil {
		logger.Errorf("unmarshall error %v", err)
		return
	}

	if responseInfo.Code != 0 {
		logger.Errorf("NdAction HTTP POST returned with exception: %v", responseInfo.Message)
		return
	}

	logger.Infof("NdAction completed successfully %v", action)
	return

}
