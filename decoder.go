package form

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type Decoder struct {
}

func NewDecoder() *Decoder {
	d := &Decoder{}
	return d
}

func (d *Decoder) Decode(v interface{}, values url.Values) (err error) {
	result := map[string]interface{}{}

	var mapkeys []string
	for k := range values {
		mapkeys = append(mapkeys, k)
	}

	sort.Slice(mapkeys, func(i, j int) bool {

		is := "0"
		js := "0"
		for e := range mapkeys[i] {
			if mapkeys[i][e] >= '0' && mapkeys[i][e] <= '9' {
				is += string(mapkeys[i][e])
			}
		}
		for e := range mapkeys[j] {
			if mapkeys[j][e] >= '0' && mapkeys[j][e] <= '9' {
				js += string(mapkeys[j][e])
			}
		}

		in, _ := strconv.ParseInt(is, 10, 64)

		jn, _ := strconv.ParseInt(js, 10, 64)

		//log.Println("sort", mapkeys[i], mapkeys[j], is, js, in, jn)
		return in < jn
	})
	//log.Println("mapkeys", mapkeys)

	for _, k := range mapkeys {
		ak, keys := keyToKeys(k)

		var v interface{}
		v = values[k]

		if len(keys) > 0 {
			for e := range v.([]string) {
				if result[ak] == nil {
					result[ak] = generateData(keys, v.([]string)[e])
				} else {
					result[ak] = merge(result[ak], generateData(keys, v.([]string)[e]))
				}
			}
		} else {

			if len(v.([]string)) == 1 {
				result[ak] = v.([]string)[0]
			} else {
				result[ak] = v
			}

		}
	}
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, v)
	return err
}

func toMapStringInterface(v interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	data, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil
	}
	return result
}
func toArrayInterface(v interface{}) []interface{} {
	var result []interface{}
	data, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil
	}
	return result
}

func merge(d, i interface{}) interface{} {

	//log.Println(reflect.TypeOf(i).Kind())
	switch reflect.TypeOf(i).Kind() {
	case reflect.Map:
		im := toMapStringInterface(i)
		dm := toMapStringInterface(d)

		for e := range im {
			if dm[e] != nil {
				switch reflect.TypeOf(dm[e]).Kind() {
				case reflect.Map:
					dm[e] = merge(dm[e], im[e])
					return dm
				case reflect.Slice:
					dm[e] = merge(dm[e], im[e])
					return dm
				default:
					dm[e] = im[e]
					return dm
				}
			} else {
				dm[e] = im[e]
				return dm
			}
		}
	case reflect.Slice:
		ia := toArrayInterface(i)
		da := toArrayInterface(d)
		da = append(da, ia...)
		return da
	}
	return d
}

func generateData(keys []string, value string) interface{} {
	if len(keys) == 1 {
		if keys[0] == "" {
			return []string{value}
		} else {
			return map[string]string{
				keys[0]: value,
			}
		}
	}
	if keys[0] == "" {
		return []interface{}{
			generateData(keys[1:], value),
		}
	} else {
		return map[string]interface{}{
			keys[0]: generateData(keys[1:], value),
		}
	}
}

func keyToKeys(key string) (ak string, keys []string) {
	key = strings.Replace(key, "[", ".", -1)
	keys = strings.Split(key, ".")
	for e := range keys {
		if keys[e] == "]" {
			keys[e] = ""
		} else {
			keys[e] = strings.Replace(keys[e], "]", "", -1)
		}
		i, err := strconv.ParseInt(keys[e], 10, 64)
		if err == nil && fmt.Sprintf("%d", i) == keys[e] {
			keys[e] = ""
		}
	}
	//log.Println("keys", keys)
	return keys[0], keys[1:]
}
