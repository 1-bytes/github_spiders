package utils

import "encoding/json"

type JsonBody map[string]interface{}

type SliceBody []JsonBody

// JsonUnmarshalBody Json 反序列化.
func JsonUnmarshalBody(jsonData []byte) (SliceBody, error) {
	body := make(SliceBody, 0)
	err := json.Unmarshal(jsonData, &body)
	if err != nil {
		body = append(body, make(JsonBody, 0))
		err = json.Unmarshal(jsonData, &body[0])
		if err != nil {
			return nil, err
		}
	}
	return body, nil
}
