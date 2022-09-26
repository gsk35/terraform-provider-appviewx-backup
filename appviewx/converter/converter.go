package converter

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

func GenerateNewMapUsingMasterAndUserInputMaps(masterMap, userInputMap map[string]interface{}) map[string]interface{} {

	for k, v := range userInputMap {
		kSlice := strings.Split(k, ".")
		var requiredMap interface{} = masterMap
		for i := 0; i < len(kSlice)-1; i++ {
			n, err := strconv.Atoi(kSlice[i])
			if err == nil && i != 0 {
				requiredMap = requiredMap.([]interface{})[n]
			} else {
				requiredMap = requiredMap.(map[string]interface{})[kSlice[i]]
			}
		}
		requiredMap.(map[string]interface{})[kSlice[len(kSlice)-1]] = v
	}
	printMap(masterMap)
	return masterMap
}

func printMap(input map[string]interface{}) {
	inputContents, err := json.Marshal(input)
	if err != nil {
		log.Println("Error in marshalling the input ", err)
	}
	log.Println(string(inputContents))
}

func GenerateNewMapUsingMasterAndUserInputMapsWithOutDot(masterMap, userInputMap *map[string]interface{}) *map[string]interface{} {
	if masterMap == nil || userInputMap == nil {
		return masterMap
	}
	for k, v := range *userInputMap {
		log.Println("***************** k :", k)
		masterMapValue, ok := (*masterMap)[k]
		if !ok {
			(*masterMap)[k] = v
			continue
		}

		log.Println("[DEBUG] reflect.ValueOf(v).String() : ", reflect.ValueOf(v).Type())
		log.Println("[DEBUG] reflect.ValueOf(masterMapValue).String() : ", reflect.ValueOf(masterMapValue).Type())

		if fmt.Sprintf("%s", reflect.ValueOf(v).Type()) == "map[string]interface {}" &&
			fmt.Sprintf("%s", reflect.ValueOf(masterMapValue).Type()) == "map[string]interface {}" {
			//Assignign to a variable due to limitation in go
			masterValueMap := masterMapValue.(map[string]interface{})
			vValueMap := v.(map[string]interface{})
			masterValueMapPointer := &masterValueMap
			vValueMapPointer := &vValueMap

			GenerateNewMapUsingMasterAndUserInputMapsWithOutDot(masterValueMapPointer, vValueMapPointer)
		} else {
			(*masterMap)[k] = v
		}
	}
	return masterMap
}
