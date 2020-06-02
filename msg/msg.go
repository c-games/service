package msg

import (
	"encoding/json"
	"github.com/c-games/common/fail"
	"github.com/c-games/service"
)

type totalAndData struct {
	Total int64         `json:"total"`
	Data json.RawMessage `json:"data"`
}

func PackCgResMsg_TotalAndData(cgMessage *service.CGMessage, errorCode int, total int64, data json.RawMessage) service.CGResMessage {
	jsonData, err := json.Marshal(totalAndData{
		Total: total,
		Data:  data,
	})
	fail.FailedOnError(err, "marshal failed")
	return service.CGResMessage{
		Serial:    cgMessage.Serial,
		Command:   cgMessage.Command,
		ErrorCode: errorCode,
		Data:      jsonData,
	}
}

func PackCgResponseMessageMany2(cgMessage *service.CGMessage, errorCode int, data []interface{}) service.CGResMessage {
	jsonData, err := json.Marshal(data)
	fail.FailedOnError(err, "marshal failed")
	return service.CGResMessage{
		Serial:    cgMessage.Serial,
		Command:   cgMessage.Command,
		ErrorCode: errorCode,
		Data:      jsonData,
	}
}

func PackCgResponseError(cgMsg *service.CGMessage, errcode int, errmsg string) service.CGResMessage {
	empty_data, _ := json.Marshal([]interface{}{})

	return service.CGResMessage{
		Serial:       cgMsg.Serial,
		Command:      cgMsg.Command,
		ErrorCode:    errcode,
		ErrorMessage: errmsg,
		Data: empty_data,
	}
}

func PackCgResponseMessage(cgMessage *service.CGMessage, errorCode int, data interface{}) service.CGResMessage {
	empty_data, _ := json.Marshal([]interface{}{})

	if data == nil {
		return service.CGResMessage{
			Serial:    cgMessage.Serial,
			Command:   cgMessage.Command,
			ErrorCode: errorCode,
			Data: empty_data,
		}
	} else {
		var resData []interface{}
		resData = append(resData, data)
		jsonRes, err := json.Marshal(resData)
		fail.FailedOnError(err, "marshal failed")
		return service.CGResMessage{
			Serial:    cgMessage.Serial,
			Command:   cgMessage.Command,
			ErrorCode: errorCode,
			Data:      jsonRes,
		}
	}

}

func ToJson(d interface{}) string {
	return string(ToByteArray(d))
}

func ToByteArray(d interface{}) []byte {
	json, err := json.Marshal(d)
	fail.FailedOnError(err, "parse json failed")
	return json
}
