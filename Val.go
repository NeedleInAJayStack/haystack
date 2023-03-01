package haystack

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

// Val represents a haystack tag value.
type Val interface {
	ToZinc() string
	json.Marshaler
	MarshalHayson() ([]byte, error)
}

// ValFromJSON converts a standard unmarshalled JSON object into the correct Haystack Val object
func ValFromJSON(jsonObj interface{}) (Val, error) {
	switch typedObj := jsonObj.(type) {
	case nil:
		return NewNull(), nil
	case bool:
		return boolFromJSON(typedObj)
	case string:
		if strings.HasPrefix(typedObj, "b:") {
			return binFromJSON(typedObj)
		} else if strings.HasPrefix(typedObj, "c:") {
			return coordFromJSON(typedObj)
		} else if strings.HasPrefix(typedObj, "d:") {
			return dateFromJSON(typedObj)
		} else if strings.HasPrefix(typedObj, "t:") {
			return dateTimeFromJSON(typedObj)
		} else if strings.HasPrefix(typedObj, "m:") {
			return markerFromJSON(typedObj)
		} else if strings.HasPrefix(typedObj, "-:") {
			return removeFromJSON(typedObj)
		} else if strings.HasPrefix(typedObj, "z:") {
			return naFromJSON(typedObj)
		} else if strings.HasPrefix(typedObj, "n:") {
			return numberFromJSON(typedObj)
		} else if strings.HasPrefix(typedObj, "r:") {
			return refFromJSON(typedObj)
		} else if strings.HasPrefix(typedObj, "h:") {
			return timeFromJSON(typedObj)
		} else if strings.HasPrefix(typedObj, "u:") {
			return uriFromJSON(typedObj)
		} else if strings.HasPrefix(typedObj, "x:") {
			return xStrFromJSON(typedObj)
		} else {
			return strFromJSON(typedObj)
		}
	case map[string]interface{}:
		if typedObj["meta"] != nil && typedObj["cols"] != nil && typedObj["rows"] != nil {
			return gridFromJSON(typedObj)
		} else {
			return dictFromJSON(typedObj)
		}
	case []interface{}:
		return listFromJSON(typedObj)
	default:
		return nil, errors.New("JSON type doesn't correlate to any haystack type: " + reflect.TypeOf(typedObj).Name())
	}
}
