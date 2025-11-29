package apigate

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/saravanan611/log"
)

type ReqStruct struct {
	RealIP      string
	ForwardedIP string
	Method      string
	Path        string
	Host        string
	RemoteAddr  string
	Header      Header
	Body        string
	EndPoint    string
	RequestType string
}

// --------------------------------------------------------------------
// get request header details
// --------------------------------------------------------------------
func GetHeaderDetails(pLog *log.LogStruct, r *http.Request) string {
	pLog.Info("GetHeaderDetails (+)")
	value1 := ""
	for name, values := range r.Header {
		// Loop over all values for the name.
		for _, value := range values {
			value1 = value1 + " " + name + "-" + value
		}
	}
	pLog.Info("GetHeaderDetails (-)")
	return value1
}

type Header map[string][]string

func (pHdr Header) String() string {
	lData, _ := json.Marshal(pHdr)
	return string(lData)
}

// --------------------------------------------------------------------
// function reads the API requestor details and send return them
// as structure to the caller
// --------------------------------------------------------------------
func GetRequestorDetail(pLog *log.LogStruct, r *http.Request) ReqStruct {
	pLog.Info("GetRequestorDetail (+)")

	var lReqRec ReqStruct
	lReqRec.RealIP = r.Header.Get("Referer")
	lReqRec.ForwardedIP = r.Header.Get("X-Forwarded-For")
	lReqRec.Method = r.Method
	lReqRec.Path = r.URL.Path + "?" + r.URL.RawQuery
	lReqRec.Host = r.Host
	lReqRec.RemoteAddr = r.RemoteAddr
	if strings.Contains(r.URL.Path, "/order/placeorder/") {
		lReqRec.EndPoint = r.URL.Path[:len("/order/placeorder/")]
	} else if strings.Contains(r.URL.Path, "/deals/count/") {
		lReqRec.EndPoint = r.URL.Path[:len("/deals/count/")]
	} else {
		lReqRec.EndPoint = r.URL.Path
	}
	lReqRec.RequestType = r.Header.Get("Content-Type")

	lReqRec.Header = Header(r.Header)
	lBody, _ := io.ReadAll(r.Body)
	lReqRec.Body = string(lBody)

	r.Body = io.NopCloser(bytes.NewBuffer([]byte(lBody)))

	pLog.Info("GetRequestorDetail (-)")

	return lReqRec
}
