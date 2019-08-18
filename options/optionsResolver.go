package options

import "fmt"

type Options map[string]interface{}

func ResolveValueInt(key string, o Options, defaultValue int, isRequired bool) (int, error) {
	val, ok := o[key]
	if !ok {
		if isRequired {
			return defaultValue, fmt.Errorf("Option '%s' is required", key)
		}
		return defaultValue, nil
	}

	valToReturn, ok := val.(int)
	if !ok {
		return defaultValue, fmt.Errorf("Option '%s' is not a valid number", key)
	}

	return valToReturn, nil
}

func ResolveValueStr(key string, o Options, defaultValue string, isRequired bool) (string, error) {
	val, ok := o[key]
	if !ok {
		if isRequired {
			return defaultValue, fmt.Errorf("Option '%s' is required", key)
		}
		return defaultValue, nil
	}

	valToReturn, ok := val.(string)
	if !ok {
		return defaultValue, fmt.Errorf("Option '%s' is not a valid string", key)
	}

	return valToReturn, nil
}

func ResolveValueStrSlice(key string, o Options, defaultValue []string, isRequired bool) ([]string, error) {
	val, keyOk := o[key]
	if !keyOk {
		if isRequired {
			return defaultValue, fmt.Errorf("Option '%s' is required", key)
		}
		return defaultValue, nil
	}

	valSlice, typeOk := val.([]string)
	if !typeOk {
		return defaultValue, fmt.Errorf("Invalid value type %v for key %s, expected type is []string", val, key)
	}

	return valSlice, nil
}
