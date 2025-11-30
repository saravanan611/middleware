package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/saravanan611/log"
)

const (
	Success = "S"
)

type RespStruct struct {
	Status   string `json:"status,omitempty"`
	ErrCode  string `json:"code,omitempty"`
	Msg      string `json:"msg,omitempty"`
	RespInfo any    `json:"info,omitempty"`
}

/* standard error responce structure for api */

func ErrorSender(w http.ResponseWriter, pLog *log.LogStruct, pErrCode string, pErr error) {
	pLog.Err(pErr)
	pLog.Info("ErrorSender (+)")
	w.WriteHeader(http.StatusInternalServerError)
	if _, lErr := fmt.Fprintf(w, "Error: << %s >>. Please refer to this code for developer fast support: (%s).", pErr.Error(), pErrCode); lErr != nil {
		pLog.Err(lErr)
	}
	pLog.Info("ErrorSender (-)")
}

func MsgSender(w http.ResponseWriter, pLog *log.LogStruct, pInfo any) {
	pLog.Info("MsgSender (+)")
	pLog.Info(pInfo)
	if lErr := json.NewEncoder(w).Encode(RespStruct{Status: Success, RespInfo: pInfo}); lErr != nil {
		pLog.Err(lErr)
	}
	pLog.Info("MsgSender (-)")
}
