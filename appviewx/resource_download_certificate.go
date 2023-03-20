package appviewx

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-provider-appviewx/appviewx/config"
	"terraform-provider-appviewx/appviewx/constants"
)

func ResourceDownloadCertificateServer() *schema.Resource {
	return &schema.Resource{
		Create: resourcedownloadServer,
		Read:   resourceCertificateServerRead,
		Update: resourceCertificateServerUpdate,
		Delete: resourceCertificateServerDelete,

		Schema: map[string]*schema.Schema{
			constants.COMMON_NAME: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			constants.SERIAL_NUMBER: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			constants.RESOURCE_ID: &schema.Schema{
				Type:     schema.TypeString,
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
	}
}

// TODO: cleanup to be done
func resourcedownloadServer(resourceData *schema.ResourceData, m interface{}) error {

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

	commonName := resourceData.Get(constants.COMMON_NAME).(string)
	serialNumber := resourceData.Get(constants.SERIAL_NUMBER).(string)
	var downloadPath, downloadFormat, downloadPassword string
	var isChainRequired bool
	if resourceData.Get(constants.CERTIFICATE_DOWNLOAD_FORMAT) == nil {
		downloadFormat = "CRT"
	} else {
		downloadFormat = resourceData.Get(constants.CERTIFICATE_DOWNLOAD_FORMAT).(string)
	}
	if resourceData.Get(constants.CERTIFICATE_DOWNLOAD_PATH) == nil && resourceData.Get(constants.COMMON_NAME) != nil {
		downloadPath = "/tmp/" + resourceData.Get(constants.COMMON_NAME).(string) + "." + strings.ToLower(downloadFormat)
	} else {
		downloadPath = resourceData.Get(constants.CERTIFICATE_DOWNLOAD_PATH).(string)
		if _, err := os.Stat(downloadPath); os.IsNotExist(err) {
			os.MkdirAll(downloadPath, 0777)
		}
		downloadPath += "/" + resourceData.Get(constants.COMMON_NAME).(string) + "." + strings.ToLower(downloadFormat)
	}
	if resourceData.Get(constants.CERTIFICATE_DOWNLOAD_PASSWORD) != nil && (downloadFormat == "PFX" || downloadFormat == "JKS" || downloadFormat == "P12") {
		downloadPassword = resourceData.Get(constants.CERTIFICATE_DOWNLOAD_PASSWORD).(string)
	} else if (resourceData.Get(constants.CERTIFICATE_DOWNLOAD_PASSWORD) == nil || resourceData.Get(constants.CERTIFICATE_DOWNLOAD_PASSWORD) == "") && (downloadFormat == "PFX" || downloadFormat == "JKS" || downloadFormat == "P12") {
		log.Println("[ERROR] Password not found for the specified download format")
		return errors.New("[ERROR] Password not found for the specified download format")
	}
	if resourceData.Get(constants.CERTIFICATE_CHAIN_REQUIRED) != nil {
		isChainRequired = resourceData.Get(constants.CERTIFICATE_CHAIN_REQUIRED).(bool)
	} else {
		isChainRequired = false
	}

	var result bool
	var resourceId string
	if commonName != "" && serialNumber != "" {
		result = downloadCertificateFromAppviewx(commonName, serialNumber, downloadFormat, downloadPassword, downloadPath, isChainRequired, appviewxSessionID, configAppViewXEnvironment)
		goto ReturnResponse
	}
	resourceId = resourceData.Get(constants.RESOURCE_ID).(string)
	if resourceId == "" {
		log.Println("[ERROR] Resource ID is not obtained to proceed with certificate download")
		return errors.New("[ERROR] Resource ID is not obtained to proceed with certificate download")
	}
	result = downloadCertificateForSyncFlow(resourceData, resourceId, appviewxSessionID, configAppViewXEnvironment)
ReturnResponse:
	if result {
		log.Println("[INFO] Certificate downloaded successfully in the specified path")
		resourceData.SetId(strconv.Itoa(rand.Int()))
	} else {
		log.Println("[ERROR] Certificate was not downloaded in the specified path")
		return errors.New("[ERROR] Certificate was not downloaded in the specified path")
	}
	return nil
}
