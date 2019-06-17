package utils

func CopyMap(originalMap map[string]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{})
	for key, values := range originalMap {
		newMap[key] = values
	}
	return newMap
}
