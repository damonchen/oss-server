package handle

import (
	"encoding/json"
	"net/http"

	"github.com/damonchen/oss-server/internal/web/base"
)

// HandleError handle error
func HandleError(w http.ResponseWriter, err error) {
	resp := base.Response{
		Code:    500,
		Message: err.Error(),
		Data:    nil,
	}

	data, _ := json.Marshal(resp)
	w.Write(data)
}

// HandleSuccess handle success
func HandleSuccess(w http.ResponseWriter, data interface{}) {
	resp := base.Response{
		Code:    200,
		Message: "",
		Data:    data,
	}

	buff, _ := json.Marshal(resp)
	w.Write(buff)
}
