package appviewx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-provider-appviewx/appviewx/config"
	"terraform-provider-appviewx/appviewx/constants"
)

func ResourceAutomationServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAutomationServerCreate,
		Read:   resourceAutomationServerRead,
		Update: resourceAutomationServerUpdate,
		Delete: resourceAutomationServerDelete,

		Schema: map[string]*schema.Schema{
			constants.APPVIEWX_ACTION_ID: {
				Type:     schema.TypeString,
				Required: true,
			},
			constants.PAYLOAD: {
				Type:     schema.TypeString,
				Optional: true,
			},
			constants.HEADERS: {
				Type:     schema.TypeMap,
				Optional: true,
			},
			constants.MASTER_PAYLOAD: {
				Type:     schema.TypeString,
				Optional: true,
			},
			constants.QUERY_PARAMS: {
				Type:     schema.TypeMap,
				Optional: true,
			},
			constants.DOWNLOAD_FILE_PATH: {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAutomationServerRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[INFO] **************** GET OPERATION NOT SUPPORTED FOR THIS RESOURCE **************** ")
	// Since the resource is for stateless operation, only nil returned
	return nil
}

func resourceAutomationServerUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("[INFO]  **************** UPDATE OPERATION NOT SUPPORTED FOR THIS RESOURCE **************** ")
	//Update implementation is empty since this resource is for the stateless generic api invocation
	return errors.New("Update not supported")
}

func resourceAutomationServerDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("[INFO]  **************** DELETE OPERATION NOT SUPPORTED FOR THIS RESOURCE **************** ")
	// Delete implementation is empty since this resoruce is for the stateless generic api invocation
	return nil
}

//TODO: cleanup to be done
func resourceAutomationServerCreate(d *schema.ResourceData, m interface{}) error {

	configAppViewXEnvironment := m.(*config.AppViewXEnvironment)

	log.Println("[DEBUG] *********************** Request received to create")
	appviewxUserName := configAppViewXEnvironment.AppViewXUserName
	appviewxPassword := configAppViewXEnvironment.AppViewXPassword
	appviewxEnvironmentIP := configAppViewXEnvironment.AppViewXEnvironmentIP
	appviewxEnvironmentPort := configAppViewXEnvironment.AppViewXEnvironmentPort
	appviewxEnvironmentIsHTTPS := configAppViewXEnvironment.AppViewXIsHTTPS
	appviewxGwSource := "WEB"

	appviewxSessionID, err := GetSession(appviewxUserName,
		appviewxPassword,
		appviewxEnvironmentIP,
		appviewxEnvironmentPort,
		appviewxGwSource,
		appviewxEnvironmentIsHTTPS)
	if err != nil {
		log.Println("[ERROR] Error in getting the session : ", err)
		return err
	}

	types := constants.POST

	actionID := d.Get(constants.APPVIEWX_ACTION_ID).(string)
	payloadString := d.Get(constants.PAYLOAD).(string)

	var masterPayloadFileName = d.Get(constants.MASTER_PAYLOAD).(string)
	if d.Get(constants.MASTER_PAYLOAD) == "" {
		masterPayloadFileName = "./payload.json"
	}

	log.Println("[DEBUG] Input minimal payload : ", payloadString)

	payloadMinimal := make(map[string]interface{})
	err = json.Unmarshal([]byte(payloadString), &payloadMinimal)
	if err != nil {
		log.Println("[ERROR] error in unmarshalling the payloadString", payloadString)
		return err
	}

	masterPayload := GetMasterPayloadApplyingMinimalPayload(masterPayloadFileName, payloadMinimal)
	log.Println("[DEBUG] masterPayload : ", masterPayload)

	queryParams := make(map[string]string)
	queryParams[constants.GW_SOURCE] = appviewxGwSource

	var queryParamReceived = d.Get(constants.QUERY_PARAMS).(map[string]interface{})
	for k, v := range queryParamReceived {
		queryParams[k] = v.(string)
	}

	url := GetURL(appviewxEnvironmentIP, appviewxEnvironmentPort, actionID, queryParams, appviewxEnvironmentIsHTTPS)

	var headers = d.Get(constants.HEADERS).(map[string]interface{})
	if len(headers) == 0 {
		headers["Content-Type"] = "application/json"
		headers["Accept"] = "application/json"
	}

	client := &http.Client{Transport: HTTPTransport()}
	requestBody, err := json.Marshal(masterPayload)
	if err != nil {
		log.Println("[ERROR] error in Marshalling the masterPayload", masterPayload, err)
		return err
	}

	printRequest(types, url, headers, requestBody)

	req, err := http.NewRequest(types, url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("[ERROR] error in creating the new request ", err)
		return err
	}

	for key, value := range headers {
		value1 := fmt.Sprintf("%v", value)
		key1 := fmt.Sprintf("%v", key)
		req.Header.Add(key1, value1)
	}
	req.Header.Add(constants.SESSION_ID, appviewxSessionID)

	resp, err := client.Do(req)
	if err != nil {
		log.Println("[ERROR] error in making http request", err)
		return err
	} else {
		log.Println("[DEBUG] Request success : url :", url)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[ERROR] error in Reading the body", err)
		return err
	}

	downloadFilePath := d.Get(constants.DOWNLOAD_FILE_PATH).(string)
	if downloadFilePath != "" {
		log.Println("downloadFilePath : ", downloadFilePath)
		err = ioutil.WriteFile(downloadFilePath, body, 0777)
		if err != nil {
			log.Println("[ERROR] error in writing the contents to file", err)
			return err
		}
	} else {
		log.Println("[DEBUG] downloadFilePath is empty")
	}

	log.Println(string(body))

	log.Println("[DEBUG] API ionvoke success")
	d.SetId(strconv.Itoa(rand.Int()))
	return resourceAutomationServerRead(d, m)
}
