package backend

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Infoblox-CTO/heka-ui/utils"
	"github.com/spf13/viper"

	"io/ioutil"
	"net/http"
	"os"

	logger "github.com/sirupsen/logrus"

	"net/http/httputil"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "API is up and running")
}

type Grid struct {
	Version     string `json:"version"`
	NiosVersion string `json:"nios_version"`
	GridName    string `json:"grid_name"`
	GridUuid    string `json:"grid_uuid"`
}
type GridResp struct {
	GridName string `json:"name"`
}

type Node struct {
	UniqueId        string       `json:"unique_id"`
	Ip              string       `json:"ip"`
	HostName        string       `json:"hostname"`
	Role            string       `json:"role"`
	Status          string       `json:"status"`
	HaEnable        bool         `json:"ha_enable"`
	MasterCandidate bool         `json:"master_candidate"`
	Hotfixes        []HotFixData `json:"hotfixes"`
}

type HotFixData struct {
	Name      string `json:"name"`
	Timestamp string `json:"timestamp"`
	Status    string `json:"status"`
}

type NiosMemberDetails struct {
	UniqueId          string     `json:"_ref"`
	IsEnableha        bool       `json:"enable_ha"`
	IsMasterCandidate bool       `json:"master_candidate"`
	NodeInfos         []NodeInfo `json:"node_info"`
	Hostname          string     `json:"host_name"`
	VipAddress        Vip        `json:"vip_setting"`
}

type NodeInfo struct {
	HaStatus      string          `json:"ha_status"`
	HardwareId    string          `json:"hwid"`
	MgmtNetwork   Mgmt            `json:"mgmt_network_setting"`
	PhysicalOid   string          `json:"physical_oid"`
	ServiceStatus []ServiceStatus `json:"service_status"`
}

type NodeHa struct {
	HaStatus    string `json:"ha_status"`
	PhysicalOid string `json:"physical_oid"`
}
type Vip struct {
	Address string `json:"address"`
}

type ServiceStatus struct {
	Description string `json:"description"`
	Service     string `json:"service"`
	Role        string `json:"role"`
}
type HaStatus struct {
	HaNodeType   string `json:"ha_node_type"`
	HaNodeStatus string `json:"ha_node_status"`
}

type Mgmt struct {
	Address string `json:"address"`
}

type Manifesta struct {
	Method string    `json:"method"`
	Object string    `json:"object"`
	Args   Arguments `json:"args"`
}

type Arguments struct {
	Returnfield string `json:"_return_fields"`
}

type LoginParams struct {
	Uname    string `json:"username"`
	Password string `json:"password"`
}

func NodeData(w http.ResponseWriter, r *http.Request) {

	sess, err := r.Cookie("user-session")
	if err != nil {
		logger.Errorln("Failed to read node data wapi cookie", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	logger.Infoln("node data browser cookie", sess.Value)
	storedSessValue, ok := SessionMap.Load(sess.Value)
	if ok {
		logger.Infoln(" node stored session ", storedSessValue)
	} else {
		logger.Errorln("node session got expired")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	member, err := prepareResponse("node", r, sess.Value, storedSessValue)
	if err != nil {
		logger.Error("Failed to get response", err)
		http.Error(w, err.Error(), 400)
		return
	}

	respNode := getNodeDetails(member[0])

	hotfixData, err := getHotfixData(respNode)
	if err != nil {
		if strings.Contains(err.Error(), "no hotfix data") {

			logger.Errorf("hotfix data file not available:%v", err)
		} else {
			logger.Errorf("Failed to fetch Hotfix data:%v", err)
			http.Error(w, err.Error(), 400)
			return
		}
	}

	if hotfixData != nil {
		respNode.Hotfixes = hotfixData
	} else {
		respNode.Hotfixes = make([]HotFixData, 0)
	}

	logger.Infoln("node data is", respNode)

	jsonResponse, err := json.Marshal(respNode)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func ListNodes(w http.ResponseWriter, r *http.Request) {

	var nodes []Node

	sess, err := r.Cookie("user-session")
	if err != nil {
		logger.Errorln("Failed to read wapi cookie", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	logger.Infoln("list browser cookie", sess.Value)
	storedSessValue, ok := SessionMap.Load(sess.Value)
	if ok {
		logger.Infoln(" list stored session ", storedSessValue)
	} else {
		logger.Errorln("list session got expired")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	members, err := prepareResponse("nodes", r, sess.Value, storedSessValue)
	if err != nil {
		logger.Errorf("Failed to get response", err)
		http.Error(w, err.Error(), 400)
		return
	}

	for _, mem := range members {
		respNode := getNodeDetails(mem)
		nodes = append(nodes, respNode)
	}

	logger.Infoln("list node data is", nodes)

	w.Header().Set("Content-Type", "application/json")

	jsonResponse, err := json.Marshal(nodes)
	if err != nil {
		logger.Errorf("failed marshalling")
		http.Error(w, err.Error(), 400)
		return
	}
	w.Write(jsonResponse)
	w.WriteHeader(http.StatusOK)

}

func getNodeDetails(mem NiosMemberDetails) Node {
	//hardware_id is obtained from params.json file
	hardware_id := viper.GetString("labels.hardware_id")
	var node Node
	if len(mem.NodeInfos) == 1 {
		node.HaEnable = mem.IsEnableha
		node.Status = "OFFLINE"
		node.HostName = mem.Hostname
		node.Ip = mem.VipAddress.Address
		node.MasterCandidate = mem.IsMasterCandidate

		for _, svc := range mem.NodeInfos[0].ServiceStatus {
			if svc.Description == "Running" && svc.Service == "NODE_STATUS" {
				node.Status = "ONLINE"
			}
		}

		if (node.MasterCandidate && (mem.NodeInfos[0].HaStatus == "ACTIVE" || mem.NodeInfos[0].HaStatus == "NOT_CONFIGURED")) && (hardware_id == mem.NodeInfos[0].HardwareId) {
			node.Role = "MASTER"
		} else {
			node.Role = "MEMBER"
		}

		node.UniqueId = mem.UniqueId
		//if its ha pair execute else clause
	} else {
		node.Status = "OFFLINE"
		node.HaEnable = mem.IsEnableha
		node.MasterCandidate = mem.IsMasterCandidate
		for _, x := range mem.NodeInfos {
			//x.HardwareId is obtained from wapi response
			if node.MasterCandidate && x.HaStatus == "ACTIVE" && (hardware_id == x.HardwareId) {
				node.Role = "MASTER"
				break
			} else {
				node.Role = "MEMBER"
			}
		}
		node.HostName = mem.Hostname
		node.Ip = mem.VipAddress.Address
		node.UniqueId = mem.UniqueId
		node.MasterCandidate = mem.IsMasterCandidate
		for _, x := range mem.NodeInfos {
			for _, s := range x.ServiceStatus {
				if s.Description == "Running" && s.Service == "NODE_STATUS" {
					node.Status = "ONLINE"
				}
			}
		}
	}
	return node
}

func getHotfixData(node Node) ([]HotFixData, error) {

	var result map[string]interface{}
	var responses []HotFixData
	var hotfixes []interface{}
	var response HotFixData
	path := fmt.Sprintf(viper.GetString("nios.hotfix.path")+"hotfix_"+"%s", node.HostName)
	if _, err := os.Stat(path); err == nil {
		f, err := os.Open(path)
		if err != nil {
			logger.Errorf("Unable to read the hotfix file")
			return nil, err
		}
		defer f.Close()

		hotfix_data, err := ioutil.ReadAll(f)
		if err != nil {
			logger.Errorf("Ioutil read failed %v", err)
			return nil, err
		}

		err = json.Unmarshal(hotfix_data, &result)
		if err != nil {
			logger.Errorf("Unmarshall failed %v", err)
			return nil, err
		}
		if result != nil {
			hotfixes = result["hotfixes"].([]interface{})
		}
		for _, y := range hotfixes {
			response.Name = y.(map[string]interface{})["version"].(string)
			if y.(map[string]interface{})["status"].(string) == "successfully" {
				response.Status = "SUCCESS"
			} else {
				response.Status = "FAILURE"
			}
			response.Timestamp = y.(map[string]interface{})["timestamp"].(string)
			responses = append(responses, response)
		}
		logger.Infof("hotfix response is %v", responses)

	} else {
		logger.Errorf("File does not exist no hotfix data available")
		return nil, fmt.Errorf("File does not exist no hotfix data available")
	}

	return responses, nil
}

func CloudManifest(w http.ResponseWriter, r *http.Request) {
	var result map[string]interface{}
	var response []byte

	f, err := os.Open(viper.GetString("nios.manifest.file"))
	if err != nil {
		logger.Errorf("unable to read the cloud manifest file", err)
		http.Error(w, err.Error(), 400)
		return
	}
	defer f.Close()
	cloud_data, err := ioutil.ReadAll(f)
	if err != nil {
		logger.Errorf("Read error %v", err)
		http.Error(w, err.Error(), 400)
		return
	}

	err = json.Unmarshal(cloud_data, &result)
	if err != nil {
		logger.Errorf("Unmarshall failed %v", err)
		http.Error(w, err.Error(), 400)
		return
	}
	response, err = json.Marshal(result)
	if err != nil {
		logger.Errorf("Marshall failed %v")
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func GridData(w http.ResponseWriter, r *http.Request) {

	sess, err := r.Cookie("user-session")
	if err != nil {
		logger.Errorln("Failed to read grid wapi cookie", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	logger.Infoln("grid browser cookie", sess.Value)
	storedSessValue, ok := SessionMap.Load(sess.Value)
	if ok {
		logger.Infoln("grid stored session ", storedSessValue)
	} else {
		logger.Errorln("session got expired login again")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	args := "name"
	var gridresp [][]GridResp
	var grid Grid
	code, resp, err := MakePostRequest(2000, "GET", "grid", args, sess.Value, storedSessValue)
	if err != nil {
		logger.Errorf("Error while doing Postrequest %v", err)
		http.Error(w, err.Error(), 400)
		return
	}
	err = json.Unmarshal(resp, &gridresp)
	if err != nil {
		logger.Errorf("Unmarshall failed %v", err, code)
		http.Error(w, err.Error(), 400)
		return
	}
	grid.Version = "1"
	grid.NiosVersion = viper.GetString("nios_version")

	for _, y := range gridresp {
		for _, n := range y {
			grid.GridName = n.GridName

		}
	}
	grid.GridUuid = viper.GetString("grid_uuid")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, err := json.Marshal(grid)
	if err != nil {
		logger.Errorf("Marshal error", err)
		http.Error(w, err.Error(), 400)
		return
	}
	w.Write(jsonResponse)
}

func prepareResponse(collectType string, r *http.Request, uuid string, storedSession interface{}) ([]NiosMemberDetails, error) {

	var nmd [][]NiosMemberDetails
	var grid_response []NiosMemberDetails
	var unique_node_response []NiosMemberDetails
	vars := mux.Vars(r)
	id := vars["uniq_object"] + "/" + vars["uniq_id"]
	args := "config_addr_type,master_candidate,platform,vip_setting,node_info,host_name,enable_ha"
	code, resp, err := MakePostRequest(2000, "GET", "member", args, uuid, storedSession)
	if err != nil {
		logger.Errorf("MakePostRequest failed %v", err)
		return grid_response, err
	}
	err = json.Unmarshal(resp, &nmd)
	if err != nil {
		logger.Errorf("unmarshall failed %v", err, code)
		return nil, err
	}

	for _, y := range nmd {
		for _, n := range y {
			grid_response = append(grid_response, n)
		}
	}

	if collectType == "nodes" {
		return grid_response, nil
	} else if collectType == "node" {
		for _, x := range grid_response {
			{
				if x.UniqueId == id {
					unique_node_response = append(unique_node_response, x)
					return unique_node_response, nil
				}
			}
		}

	}

	return nil, nil
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

	code, response, err := utils.MakePostRequest(timeout, script, args)

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

}

//MakePostRequest will do HTTP POST request to connect to NIOS
func MakePostRequest(sec time.Duration, Method, object string, args string, sesskey string, sessCookie interface{}) (int, []byte, error) {
	url := "https://" + "127.0.0.1" + "/wapi/" + viper.GetString("wapi.version") + "/request"
	var input []Manifesta
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cli := &http.Client{
		Transport: tr,
		Timeout:   sec * time.Second,
	}

	input = []Manifesta{

		{
			Method: Method,
			Object: object,
			Args: Arguments{
				Returnfield: args,
			},
		},
	}

	jsonStr, err := json.Marshal(input)
	if err != nil {
		fmt.Printf("Error while creating json input for API call: %v", err)
		return 0, nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		logger.Errorf("Error while creating the request: %v", err)
		return 0, nil, err
	}
	//set request headers
	cookie := &http.Cookie{
		Name: sessCookie.(*http.Cookie).String(),
	}
	//set request headers
	req.AddCookie(cookie)
	req.Header.Add("Content-Type", "application/json")
	//send request

	res, err := cli.Do(req)
	if err != nil {
		logger.Errorf("Request failed - %v", err)
		return 0, nil, err
	}
	cke := res.Cookies()
	if len(cke) > 0 {
		logger.Infoln("New response Cookie: ", cke)
		SessionMap.Store(sesskey, cke[0])
		logger.Infoln("Updated new cookie in map")
	} else {
		logger.Warn("empty cookie returned from nios")
	}
	defer res.Body.Close()
	resDump, _ := httputil.DumpResponse(res, true)
	logger.Debugf("Response dump - %s", string(resDump))
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Failed to read POST heka_data status response body - %v", err)
		return 0, nil, err
	}

	return res.StatusCode, data, err
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func Login(w http.ResponseWriter, r *http.Request) {

	var (
		login          LoginParams
		adminGroup     string
		result         []byte
		res            *http.Response
		adminUser      []map[string]interface{}
		err            error
		maxRetryPeriod = 2 * time.Minute
		retryInterval  = 30 * time.Second
	)
	// TODO check creds or cookie
	// Let's handle multiple tab seesion management later
	urlprofile := "https://" + "127.0.0.1" + "/wapi/" + viper.GetString("wapi.version") + "/userprofile?_return_fields=name,admin_group"

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&login); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	for {
		result, res, err = validateUser(urlprofile, login)
		if err != nil {
			if retryInterval > maxRetryPeriod {
				logger.Errorf("validate user failed for userprofile %v", err)
				continue
			} else {
				time.Sleep(retryInterval)
				retryInterval = retryInterval * 2
				logger.Errorf("Retrying  to validate user profile for %v", r)
			}
		}
		break
	}
	if res.StatusCode != http.StatusOK {
		http.Error(w, "Invalid status code", res.StatusCode)
		logger.Errorf("Unauthorized error:%v", res.StatusCode)
		return
	}

	err = json.Unmarshal(result, &adminUser)
	if err != nil {
		logger.Errorf("Failed to unmarshall during validate user for Userprofile:%v", err)
		return
	}

	for _, y := range adminUser {
		if y["admin_group"] != nil {
			adminGroup = y["admin_group"].(string)
			break
		}
	}

	urladminGroup := "https://" + "127.0.0.1" + "/wapi/" + viper.GetString("wapi.version") + "/admingroup?name=" + adminGroup + "&_return_fields=superuser"

	for {
		result, res, err = validateUser(urladminGroup, login)
		if err != nil {

			if retryInterval > maxRetryPeriod {
				logger.Errorf("validate user failed for userprofile %v", err)
				continue
			} else {
				time.Sleep(retryInterval)
				retryInterval = retryInterval * 2
				logger.Errorf("Retrying  to validate user profile for %v", r)
			}

		}
		break
	}

	if res.StatusCode != http.StatusOK {
		http.Error(w, "Invalid status code", res.StatusCode)
		logger.Errorf("Invalid status code:%v", res.StatusCode)
		return
	}

	err = json.Unmarshal(result, &adminUser)
	if err != nil {
		http.Error(w, "Failed to unmarshall", http.StatusBadRequest)
		logger.Errorf("Failed to unmarshall")
		return
	}
	var superUser bool
	for _, y := range adminUser {
		if y["superuser"] != nil {
			superUser = y["superuser"].(bool)
			if superUser == true {
				logger.Infoln("its valid usecase pls proceed")
				break

			} else {
				respondWithError(w, http.StatusBadRequest, "Only superuser can be logged in")
				logger.Errorf("Only superuser can be logged in")
				return
			}
		}
	}
	logger.Infof("superuser flag is %v  if true proceed to login", superUser)
	retryInterval = 30 * time.Second
	for {
		res, err = callWapiObject(w, r, "GET", "grid", "name", login)
		if err != nil {

			if retryInterval > maxRetryPeriod {
				logger.Errorf("callWapiObject failed %v", err)
				continue
			} else {
				time.Sleep(retryInterval)
				retryInterval = retryInterval * 2
				logger.Errorf("Retrying callWapiObject for %v", r)
			}
		}
		break
	}
	okresponse := res.StatusCode
	logger.Infoln("okresponse", okresponse)

	if res.StatusCode != 201 {
		respondWithError(w, http.StatusUnauthorized, "Login Failed")
		return
	}
	logger.Infoln("Response Cookie: ", res.Cookies())
	sess := cookieManager(res.Cookies())

	if sess != "" {
		cookie := &http.Cookie{
			Name:  "user-session",
			Value: sess,
		}
		http.SetCookie(w, cookie)
	} else {
		logger.Errorln("cookiemanager session is empty")
		//TODO handle
	}

	logger.Infoln("Login successfull")
	payload := struct {
		message string
	}{
		message: "login successfull",
	}
	respondWithJSON(w, http.StatusOK, payload)
}

func callWapiObject(w http.ResponseWriter, r *http.Request, method, object, returnField string, login LoginParams) (*http.Response, error) {
	url := "https://" + "127.0.0.1" + "/wapi/" + viper.GetString("wapi.version") + "/request"
	var input []Manifesta
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	cli := &http.Client{
		Transport: tr,
		Timeout:   2600 * time.Second,
	}
	// create request
	input = []Manifesta{
		{
			Method: method,
			Object: object,
			Args: Arguments{
				Returnfield: returnField,
			},
		},
	}

	jsonStr, err := json.Marshal(input)
	if err != nil {
		fmt.Printf("Error while creating json input for API call: %v", err)
		respondWithError(w, http.StatusInternalServerError, "json marshal error")
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Printf("Error while creating the request: %v", err)
		respondWithError(w, http.StatusInternalServerError, "http request error")
	}
	//set request headers
	req.SetBasicAuth(login.Uname, login.Password)
	req.Header.Add("Content-Type", "application/json")
	res, err := cli.Do(req)
	if err != nil {
		logger.Errorf("Request failed - %v", err)
		respondWithError(w, http.StatusInternalServerError, "server request error")
		return res, err
	}
	defer res.Body.Close()

	return res, nil
}

// validateUser validates whteher user is superuser or non super user
func validateUser(url string, login LoginParams) ([]byte, *http.Response, error) {

	//var input []Manifesta
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	cli := &http.Client{
		Transport: tr,
		Timeout:   2600 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("Error while creating the request: %v", err)
		//respondWithError(w, http.StatusInternalServerError, "http request error")
	}
	//set request headers
	req.SetBasicAuth(login.Uname, login.Password)
	req.Header.Add("Content-Type", "application/json")
	res, err := cli.Do(req)
	if err != nil {
		logger.Errorf("Request failed - %v", err)
		//respondWithError(w, http.StatusInternalServerError, "server request error")
		return nil, res, err
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Errorf("failed to read response body %v", err)
		return nil, res, err
	}
	return data, res, err

}

func cookieManager(cke []*http.Cookie) string {
	id, err := uuid.NewRandom()
	if err != nil {
		logger.Errorf("Failed to generate random number", err)
		return ""
	}

	uuid := id.String()
	logger.Infoln("random uuid :", uuid)
	if len(cke) != 0 {
		SessionMap.Store(uuid, cke[0])
	} else {
		logger.Errorln("Failure to decode session from cookie")
		return ""
	}
	return uuid
}

func Logout(w http.ResponseWriter, r *http.Request) {
	logger.Infoln("triggered logout request")
	sess, err := r.Cookie("user-session")
	if err != nil {
		logger.Errorln("logout failed to read wapi cookie", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	logger.Infoln("logout client cookie", sess.Value)
	actVal, ok := SessionMap.Load(sess.Value)
	if ok {
		logger.Infoln("logout session stored", actVal)
		SessionMap.Delete(sess.Value)
	} else {
		logger.Errorln("logout session already got expired")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}
