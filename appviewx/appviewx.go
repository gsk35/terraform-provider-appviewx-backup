package appviewx

import (
	"terraform-provider-appviewx/appviewx/converter"
	"terraform-provider-appviewx/appviewx/fileops"
)

func GetURL(ip, port, actionId string, queryParams map[string]string, isHttps bool) (output string) {
	if isHttps {
		output += "https://"
	} else {
		output += "http://"
	}
	output += ip + ":" + port + "/avxapi/" + actionId + "?"

	if queryParams != nil {
		for k, v := range queryParams {
			output += (k + "=" + v + "&")
		}
	}
	output = output[:len(output)-1]
	return
}

func GetMasterPayloadApplyingMinimalPayload(masterPayloadFileName string, payloadMinimal map[string]interface{}) map[string]interface{} {
	masterPayload := fileops.GetFileContentsInMap(masterPayloadFileName)
	return *converter.GenerateNewMapUsingMasterAndUserInputMapsWithOutDot(&masterPayload, &payloadMinimal)
}
