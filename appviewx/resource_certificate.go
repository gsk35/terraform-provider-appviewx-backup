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
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-provider-appviewx/appviewx/config"
	"terraform-provider-appviewx/appviewx/constants"
)

func ResourceCertificateServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceCertificateServerCreate,
		Read:   resourceCertificateServerRead,
		Update: resourceCertificateServerUpdate,
		Delete: resourceCertificateServerDelete,

		Schema: map[string]*schema.Schema{
			constants.APPVIEWX_ACTION_ID: {
				Type:     schema.TypeString,
				Required: true,
			},
			constants.PAYLOAD: {
				Type:     schema.TypeString,
				Optional: true,
			},
			constants.TYPE: {
				Type:     schema.TypeString,
				Required: true,
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
func resourceCertificateServerRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[INFO]  **************** GET OPERATION NOT SUPPORTED FOR THIS RESOURCE **************** ")
	// Since the resource is for stateless operation, only nil returned
	return nil
}

func resourceCertificateServerUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("[INFO]  **************** UPDATE OPERATION NOT SUPPORTED FOR THIS RESOURCE **************** ")
	//Update implementation is empty since this resource is for the stateless generic api invocation
	return errors.New("Update not supported")
}

func resourceCertificateServerDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("[INFO]  **************** DELETE OPERATION NOT SUPPORTED FOR THIS RESOURCE **************** ")
	// Delete implementation is empty since this resoruce is for the stateless generic api invocation
	return nil
}

//TODO: cleanup to be done
func resourceCertificateServerCreate(d *schema.ResourceData, m interface{}) error {

	configAppViewXEnvironment := m.(*config.AppViewXEnvironment)

	log.Println("[INFO] *********************** Request received to create")
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
		appviewxGwSource, appviewxEnvironmentIsHTTPS)
	if err != nil {
		log.Println("[ERROR] Error in getting the session : ", err)
		return err
	}

	types := strings.ToUpper(d.Get(constants.TYPE).(string))
	if types == constants.POST || types == constants.PUT || types == constants.DELETE || types == constants.GET {

		actionID := d.Get(constants.APPVIEWX_ACTION_ID).(string)
		payloadString := d.Get(constants.PAYLOAD).(string)

		var masterPayloadFileName = d.Get(constants.MASTER_PAYLOAD).(string)
		if d.Get(constants.MASTER_PAYLOAD) == "" {
			masterPayloadFileName = "./payload.json"
		}

		log.Println("[DEBG] Input minimal payload : ", payloadString)

		payloadMinimal := make(map[string]interface{})
		err = json.Unmarshal([]byte(payloadString), &payloadMinimal)
		if err != nil {
			log.Println("[ERROR] error in unmarshalling the payloadString", payloadString, err)
			return err
		}

		masterPayload := GetMasterPayloadApplyingMinimalPayload(masterPayloadFileName, payloadMinimal)
		log.Println("[DEBG] masterPayload : ", masterPayload)

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
			log.Println("[ERROR] error in Marshalling the payload ", masterPayload, err)
			return err
		}

		printRequest(types, url, headers, requestBody)

		req, err := http.NewRequest(types, url, bytes.NewBuffer(requestBody))
		if err != nil {
			log.Println("[ERROR] error in creating new Request", err)
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
			log.Println("[ERROR] error in http request", err)
			return err
		} else {
			log.Println("[DEBG] Request success : url :", url)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("[ERROR] error in reading the response Body", err)
			return err
		}

		downloadFilePath := d.Get(constants.DOWNLOAD_FILE_PATH).(string)
		if downloadFilePath != "" {
			log.Println("[DEBG] downloadFilePath : ", downloadFilePath)
			err = ioutil.WriteFile(downloadFilePath, body, 0777)
			if err != nil {
				log.Println("[ERROR] error in writing the contents to file", err)
				return err
			}
		} else {
			log.Println("[DEBG] downloadFilePath is empty")
		}

		log.Println(string(body))

		log.Println("[DEBG] API ionvoke success")
		d.SetId(strconv.Itoa(rand.Int()))
		return resourceCertificateServerRead(d, m)
	}
	return nil
}
