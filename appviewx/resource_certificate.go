package appviewx

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-provider-appviewx/appviewx/config"
	"terraform-provider-appviewx/appviewx/constants"
)

func ResourceCertificateServer() *schema.Resource {
	//fmt.Println("****************** SSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSS")
	return &schema.Resource{
		Create: resourceCertificateServerCreate,
		Read:   resourceCertificateServerRead,
		Update: resourceCertificateServerUpdate,
		Delete: resourceCertificateServerDelete,

		Schema: map[string]*schema.Schema{
			constants.COMMON_NAME: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			constants.HASH_FUNCTION: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			constants.KEY_TYPE: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			constants.BIT_LENGTH: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			constants.DNS_NAMES: &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			constants.CUSTOM_FIELDS: &schema.Schema{
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			constants.VENDOR_SPECIFIC_FIELDS: &schema.Schema{
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			constants.CERTIFICATE_AUTHORITY: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			constants.CA_SETTING_NAME: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			constants.VALIDITY: &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			constants.IS_SYNC: &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			constants.CERTIFICATE_DOWNLOAD_PATH: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			constants.CERTIFICATE_DOWNLOAD_FORMAT: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			constants.CERTIFICATE_DOWNLOAD_PASSWORD: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			constants.CERTIFICATE_CHAIN_REQUIRED: &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourceCertificateImport,
		},
	}
}

func resourceCertificateImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	id := d.Id()

	parameters := strings.Split(id, ",")

	fmt.Println("parameters = ", parameters)

	return []*schema.ResourceData{d}, nil
}

func resourceCertificateServerRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[INFO]  **************** GET OPERATION NOT SUPPORTED FOR THIS RESOURCE **************** ")
	// Since the resource is for stateless operation, only nil returned
	return nil
}

func resourceCertificateServerUpdate(resourceData *schema.ResourceData, m interface{}) error {
	log.Println("[INFO]  **************** UPDATE OPERATION NOT SUPPORTED FOR THIS RESOURCE **************** ")
	//Update implementation is empty since this resource is for the stateless generic api invocation
	return errors.New("Update not supported")
}

func resourceCertificateServerDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("[INFO]  **************** DELETE OPERATION NOT SUPPORTED FOR THIS RESOURCE **************** ")
	// Delete implementation is empty since this resoruce is for the stateless generic api invocation
	return nil
}

// TODO: cleanup to be done
func resourceCertificateServerCreate(resourceData *schema.ResourceData, m interface{}) error {

	fmt.Println("****************** Resource Certificate Server Create ******************")
	configAppViewXEnvironment := m.(*config.AppViewXEnvironment)

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
		return nil
	}

	result := createCertificate(resourceData, configAppViewXEnvironment, appviewxSessionID)
	if result.Response["resourceId"] == "" {
		log.Println("[ERROR] Resource ID is not obtained to proceed with certificate download")
		return errors.New("[ERROR] Resource ID is not obtained to proceed with certificate download")
	}
	resourceID := result.Response["resourceId"]
	if resourceData.Get(constants.IS_SYNC) == nil || !resourceData.Get(constants.IS_SYNC).(bool) {
		log.Println("[INFO] Certificate is created in ASYNC mode so download can be done once the certificate is issued.")
		log.Println("[INFO] ***** Use this resource ID to download the certificate", resourceID)
		resourceData.SetId(strconv.Itoa(rand.Int()))
		return nil
	} else {
		result := downloadCertificateForSyncFlow(resourceData, resourceID, appviewxSessionID, configAppViewXEnvironment)
		if result {
			log.Println("[INFO] Certificate downloaded successfully in the specified path")
			resourceData.SetId(strconv.Itoa(rand.Int()))
		} else {
			log.Println("[ERROR] Certificate was not downloaded in the specified path")
			return errors.New("[ERROR] Certificate was not downloaded in the specified path")
		}
	}
	return nil
}

func createCertificate(resourceData *schema.ResourceData, configAppViewXEnvironment *config.AppViewXEnvironment, appviewxSessionID string) config.AppviewxCreateCertResponse {
	var result config.AppviewxCreateCertResponse
	httpMethod := config.HTTPMethodPost
	appviewxEnvironmentIP := configAppViewXEnvironment.AppViewXEnvironmentIP
	appviewxEnvironmentPort := configAppViewXEnvironment.AppViewXEnvironmentPort
	appviewxEnvironmentIsHTTPS := configAppViewXEnvironment.AppViewXIsHTTPS
	queryParams := frameQueryParams()
	if resourceData.Get(constants.IS_SYNC) != nil {
		isSync := resourceData.Get(constants.IS_SYNC).(bool)
		queryParams["isSync"] = strconv.FormatBool(isSync)
	}
	headers := frameHeaders()
	url := GetURL(appviewxEnvironmentIP, appviewxEnvironmentPort, config.CreateCertificateActionId, queryParams, appviewxEnvironmentIsHTTPS)
	payload := frameCertificatePayload(resourceData)
	requestBody, err := json.Marshal(payload)
	if err != nil {
		log.Println("[ERROR] error in Marshalling the payload ", payload, err)
		return result
	}
	client := &http.Client{Transport: HTTPTransport()}

	printRequest(httpMethod, url, headers, requestBody)

	req, err := http.NewRequest(httpMethod, url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("[ERROR] error in creating new Request", err)
		return result
	}

	for key, value := range headers {
		value1 := fmt.Sprintf("%v", value)
		key1 := fmt.Sprintf("%v", key)
		req.Header.Add(key1, value1)
	}
	req.Header.Add(constants.SESSION_ID, appviewxSessionID)

	httpResponse, err := client.Do(req)
	if err != nil {
		log.Println("[ERROR] error in http request", err)
		return result
	} else {
		log.Println("[DEBUG] Request success : url :", url)
	}
	if !(httpResponse.StatusCode == 200 || httpResponse.StatusCode == 201 || httpResponse.StatusCode == 202) {
		log.Println("[ERROR] Response.Status : ", httpResponse.Status)
		log.Printf("[ERROR] Error in making http client request with status code : %d\n", httpResponse.StatusCode)
	}
	responseByte, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		log.Println(err)
	} else {
		err = json.Unmarshal(responseByte, &result)
		if err != nil {
			log.Println("[ERROR] ", err)
		} else {
			log.Println("[INFO] Response for Certificate creation = ", result)
		}
	}
	return result

}

func downloadCertificateForSyncFlow(resourceData *schema.ResourceData, appviewxResourceId, appviewxSessionID string, configAppViewXEnvironment *config.AppViewXEnvironment) bool {
	var downloadPath, downloadFormat, downloadPassword, commonName string
	var isChainRequired bool
	if resourceData.Get(constants.COMMON_NAME) != nil && resourceData.Get(constants.COMMON_NAME).(string) != "" {
		commonName = resourceData.Get(constants.COMMON_NAME).(string)
	} else {
		log.Println("[INFO] Commona name is not specified to download the certificate so downloading with name enrolledCertificate")
		commonName = "enrolledCertificate"
	}

	if resourceData.Get(constants.CERTIFICATE_DOWNLOAD_FORMAT) == nil {
		downloadFormat = "CRT"
	} else {
		downloadFormat = resourceData.Get(constants.CERTIFICATE_DOWNLOAD_FORMAT).(string)
	}
	if resourceData.Get(constants.CERTIFICATE_DOWNLOAD_PATH) == nil && resourceData.Get(constants.COMMON_NAME) != nil {
		downloadPath = "/tmp/" + commonName + "." + strings.ToLower(downloadFormat)
	} else {
		downloadPath = resourceData.Get(constants.CERTIFICATE_DOWNLOAD_PATH).(string)
		if _, err := os.Stat(downloadPath); os.IsNotExist(err) {
			os.MkdirAll(downloadPath, 0777)
		}
		downloadPath += "/" + commonName + "." + strings.ToLower(downloadFormat)
	}
	if resourceData.Get(constants.CERTIFICATE_DOWNLOAD_PASSWORD) != nil && (downloadFormat == "PFX" || downloadFormat == "JKS" || downloadFormat == "P12") {
		downloadPassword = resourceData.Get(constants.CERTIFICATE_DOWNLOAD_PASSWORD).(string)
	} else if (resourceData.Get(constants.CERTIFICATE_DOWNLOAD_PASSWORD) == nil || resourceData.Get(constants.CERTIFICATE_DOWNLOAD_PASSWORD) == "") && (downloadFormat == "PFX" || downloadFormat == "JKS" || downloadFormat == "P12") {
		log.Println("[ERROR] Password not found for the specified download format")
		return false
	}
	if resourceData.Get(constants.CERTIFICATE_CHAIN_REQUIRED) != nil {
		isChainRequired = resourceData.Get(constants.CERTIFICATE_CHAIN_REQUIRED).(bool)
	} else {
		isChainRequired = false
	}

	searchResponse := searchCertificate(appviewxResourceId, appviewxSessionID, configAppViewXEnvironment)

	if searchResponse.AppviewxResponse.ResponseObject.Objects[0].CommonName == "" || searchResponse.AppviewxResponse.ResponseObject.Objects[0].SerialNumber == "" {
		log.Println("[ERROR] Cannot find the common name and serial number for the resource id " + appviewxResourceId + " to proceed with certificate download")
		return false
	}
	commonName = searchResponse.AppviewxResponse.ResponseObject.Objects[0].CommonName
	serialNumber := searchResponse.AppviewxResponse.ResponseObject.Objects[0].SerialNumber
	return downloadCertificateFromAppviewx(commonName, serialNumber, downloadFormat, downloadPassword, downloadPath, isChainRequired, appviewxSessionID, configAppViewXEnvironment)
}

func searchCertificate(resourceID, appviewxSessionID string, configAppViewXEnvironment *config.AppViewXEnvironment) config.AppviewxSearchCertResponse {
	var response config.AppviewxSearchCertResponse
	httpMethod := config.HTTPMethodPost
	appviewxEnvironmentIP := configAppViewXEnvironment.AppViewXEnvironmentIP
	appviewxEnvironmentPort := configAppViewXEnvironment.AppViewXEnvironmentPort
	appviewxEnvironmentIsHTTPS := configAppViewXEnvironment.AppViewXIsHTTPS
	headers := frameHeaders()
	url := GetURL(appviewxEnvironmentIP, appviewxEnvironmentPort, config.SearchCertificateActionId, frameQueryParams(), appviewxEnvironmentIsHTTPS)
	payload := frameSearchCertificatePayload(resourceID)
	requestBody, err := json.Marshal(payload)
	if err != nil {
		log.Println("[ERROR] error in Marshalling the payload ", payload, err)
		return response
	}
	client := &http.Client{Transport: HTTPTransport()}

	printRequest(httpMethod, url, headers, requestBody)

	req, err := http.NewRequest(httpMethod, url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("[ERROR] error in creating new Request", err)
		return response
	}

	for key, value := range headers {
		value1 := fmt.Sprintf("%v", value)
		key1 := fmt.Sprintf("%v", key)
		req.Header.Add(key1, value1)
	}
	req.Header.Add(constants.SESSION_ID, appviewxSessionID)

	httpResponse, err := client.Do(req)
	if err != nil {
		log.Println("[ERROR] error in http request", err)
		return response
	} else {
		log.Println("[DEBUG] Request success : url :", url)
	}
	if !(httpResponse.StatusCode == 200 || httpResponse.StatusCode == 201 || httpResponse.StatusCode == 202) {
		log.Println("[ERROR] Response.Status : ", httpResponse.Status)
		log.Printf("[ERROR] Error in making http client request with status code : %d\n", httpResponse.StatusCode)
	}
	responseByte, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		log.Println(err)
	} else {
		err = json.Unmarshal(responseByte, &response)
		if err != nil {
			log.Println("[ERROR] ", err)
		} else {
			log.Println("[INFO] Response for Certificate Search = ", response)
		}
	}
	return response
}

func downloadCertificateFromAppviewx(commonName, serialNumber, downloadFormat, downloadPassword, downloadPath string, isChainRequired bool, appviewxSessionID string, configAppViewXEnvironment *config.AppViewXEnvironment) bool {
	httpMethod := config.HTTPMethodPost
	appviewxEnvironmentIP := configAppViewXEnvironment.AppViewXEnvironmentIP
	appviewxEnvironmentPort := configAppViewXEnvironment.AppViewXEnvironmentPort
	appviewxEnvironmentIsHTTPS := configAppViewXEnvironment.AppViewXIsHTTPS
	headers := frameHeaders()
	url := GetURL(appviewxEnvironmentIP, appviewxEnvironmentPort, config.DownloadCertificateActionId, frameQueryParams(), appviewxEnvironmentIsHTTPS)
	payload := frameDownloadCertificatePayload(commonName, serialNumber, downloadFormat, downloadPassword, isChainRequired)
	requestBody, err := json.Marshal(payload)
	if err != nil {
		log.Println("[ERROR] error in Marshalling the payload ", payload, err)
		return false
	}
	client := &http.Client{Transport: HTTPTransport()}

	printRequest(httpMethod, url, headers, requestBody)

	req, err := http.NewRequest(httpMethod, url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("[ERROR] error in creating new Request", err)
		return false
	}

	for key, value := range headers {
		value1 := fmt.Sprintf("%v", value)
		key1 := fmt.Sprintf("%v", key)
		req.Header.Add(key1, value1)
	}
	req.Header.Add(constants.SESSION_ID, appviewxSessionID)

	httpResponse, err := client.Do(req)
	if err != nil {
		log.Println("[ERROR] error in http request", err)
		return false
	} else {
		log.Println("[DEBUG] Request success : url :", url)
	}
	if !(httpResponse.StatusCode == 200 || httpResponse.StatusCode == 201 || httpResponse.StatusCode == 202) {
		log.Println("[ERROR] Response.Status : ", httpResponse.Status)
		log.Printf("[ERROR] Error in making http client request with status code : %d\n", httpResponse.StatusCode)
	}
	responseByte, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		log.Println("[ERROR] ", err)
		return false
	} else {
		err = os.WriteFile(downloadPath, responseByte, 0777)
		if err != nil {
			log.Println("[ERROR] Error while writting certificate file content in ", downloadPath, " due to : ", err)
			return false
		} else {
			log.Println("[INFO] Downloaded certificate file content in ", downloadPath)
			return true
		}
	}

}

func frameSearchCertificatePayload(resourceId string) config.SearchCertificatePayload {
	var payload config.SearchCertificatePayload
	payload.Filter.SortOrder = "asc"
	payload.Input.ResourceId = resourceId
	return payload
}

func frameDownloadCertificatePayload(commonName, serialNumber, format, password string, isChainRequired bool) config.DownloadCertificatePayload {
	var payload config.DownloadCertificatePayload
	payload.CommonName = commonName
	payload.SerialNumber = serialNumber
	payload.Format = format
	payload.IsChainRequired = isChainRequired
	payload.Password = password
	return payload
}
func frameCertificatePayload(resourceData *schema.ResourceData) config.CreateCertificatePayload {
	var payload config.CreateCertificatePayload
	var csrParams config.CSRParameters
	csrParams.CommonName = resourceData.Get(constants.COMMON_NAME).(string)
	csrParams.HashFunction = resourceData.Get(constants.HASH_FUNCTION).(string)
	csrParams.KeyType = resourceData.Get(constants.KEY_TYPE).(string)
	csrParams.BitLength = resourceData.Get(constants.BIT_LENGTH).(string)
	dnsNames, ok := resourceData.GetOk(constants.DNS_NAMES)
	var enhancedSAN config.EnhancedSANTypes
	if ok {
		dns := dnsNames.([]interface{})
		var dnsValues = make([]string, len(dns))
		for key, value := range dns {
			dnsValues[key] = value.(string)
		}
		enhancedSAN.DNSNames = dnsValues
		csrParams.EnhancedSANTypes = enhancedSAN
	}
	csrParams.CertificateCategories = []string{"Server", "Client"}
	payload.CaConnectorInfo.CSRParameters = csrParams
	payload.CaConnectorInfo.CASettingName = resourceData.Get(constants.CA_SETTING_NAME).(string)
	payload.CaConnectorInfo.CertificateAuthority = resourceData.Get(constants.CERTIFICATE_AUTHORITY).(string)
	payload.CaConnectorInfo.CAConnectorName = payload.CaConnectorInfo.CertificateAuthority + " Connector  Terraform"
	payload.CaConnectorInfo.ValidityInDays = resourceData.Get(constants.VALIDITY).(int)
	customFields, ok := resourceData.GetOk(constants.CUSTOM_FIELDS)
	if ok {
		var customFieldValues = make(map[string]string)
		customFields := customFields.(map[string]interface{})
		for key, values := range customFields {
			customFieldValues[key] = values.(string)
		}
		payload.CaConnectorInfo.CustomAttributes = customFieldValues
	}
	vendorSpecFields, ok := resourceData.GetOk(constants.VENDOR_SPECIFIC_FIELDS)
	if ok {
		var vendorFields = make(map[string]string)
		vendorSpecFieldList := vendorSpecFields.(map[string]interface{})
		for key, values := range vendorSpecFieldList {
			vendorFields[key] = values.(string)
		}
		payload.CaConnectorInfo.VendorSpecificfields = vendorFields
	}
	return payload
}

func frameHeaders() map[string]interface{} {
	var headers = make(map[string]interface{})
	headers["Content-Type"] = "application/json"
	headers["Accept"] = "application/json"
	return headers
}

func frameQueryParams() map[string]string {
	var queryParams = make(map[string]string)
	queryParams[constants.GW_SOURCE] = "WEB"
	return queryParams
}
