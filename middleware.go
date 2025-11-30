package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/saravanan611/log"
)

type ResponseCaptureWriter struct {
	http.ResponseWriter
	status int
	body   []byte
}

func (rw *ResponseCaptureWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *ResponseCaptureWriter) Write(body []byte) (int, error) {
	rw.body = append(rw.body, body...)
	return rw.ResponseWriter.Write(body)
}

func (rw *ResponseCaptureWriter) Status() int {
	if rw.status == 0 {
		return http.StatusOK
	}
	return rw.status
}

func (rw *ResponseCaptureWriter) Body() []byte {
	return rw.body
}

var (
	allowOrigin     = "*"
	allowCredential = false
	allowHeader     = []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "credentials"}
	methods         = []string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodConnect, http.MethodOptions, http.MethodTrace}
)

/* set the header for your request  */
func SetHeader(pHeader ...string) {
	if len(pHeader) > 0 {
		allowHeader = append(allowHeader, pHeader...)
	}
}

/* set the origin for your request ,default "*"" */
func SetOrigin(pOrigin string) {
	allowOrigin = pOrigin
}

/* enable Credential to true for cookie_set,some other browser side operation dun by golang ,default "false" */
func EnableCredential() {
	allowCredential = true
}

type ownStr string

const GateKey ownStr = "G-key"

// Middleware to log requests and route based on API version
func logMiddleware(pNext http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Initialize the logger
		(w).Header().Set("Access-Control-Allow-Origin", allowOrigin)
		(w).Header().Set("Access-Control-Allow-Credentials", fmt.Sprint(allowCredential))
		(w).Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
		(w).Header().Set("Access-Control-Allow-Headers", strings.Join(allowHeader, ","))

		log := log.Init()
		log.Info("LogMiddleware (+)")
		// Check if it is an OPTIONS request

		requestorDetail := GetRequestorDetail(log, r)

		ctx := context.WithValue(r.Context(), GateKey, log.Uid)
		r = r.WithContext(ctx)

		contentLength := r.Header.Get("Content-Length")
		if contentLength != "" {
			length, lErr := strconv.Atoi(contentLength)
			if lErr != nil {
				log.Err(lErr)
			}

			// if length >= 2<<20 { // 2 MB = 2 * 1024 * 1024 bytes
			if length >= 1<<20 { // 1 MB = 1 * 1024 * 1024 bytes
				requestorDetail.Body = "File Data"
			}
		}

		log.Info("Req Info :", requestorDetail)

		// Move the logging of request after setting the context
		captureWriter := &ResponseCaptureWriter{ResponseWriter: w}

		captureWriter.Header().Set("", log.Uid)

		pNext.ServeHTTP(captureWriter, r)

		log.Info("Resp Info :", string(captureWriter.Body()))

		log.Info("LogMiddleware (-)")
	})

}

/* set up your server to execuate  */
func SetServer(pRuterFunc func(pRouterInfo *mux.Router), pReadTimeout, pWriteTimeout, pIdleTimeout, pPortAdrs int) error {

	if pReadTimeout == 0 {
		pReadTimeout = 30
	}
	if pWriteTimeout == 0 {
		pWriteTimeout = 30
	}
	if pIdleTimeout == 0 {
		pIdleTimeout = 120
	}

	lRouter := mux.NewRouter()
	pRuterFunc(lRouter)

	lRouter.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "Optional Call Success")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, `{"status":"E","error": "Method %s not allowed on %s"}`, r.Method, r.URL.Path)
	})

	lHandler := logMiddleware(lRouter)
	lSrv := &http.Server{
		ReadTimeout:  time.Duration(pReadTimeout) * time.Second,
		WriteTimeout: time.Duration(pWriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(pIdleTimeout) * time.Second,
		Handler:      lHandler,
		Addr:         fmt.Sprintf(":%d", pPortAdrs),
	}

	fmt.Printf("server start on :%d ....", pPortAdrs)
	if lErr := lSrv.ListenAndServe(); lErr != nil {
		return log.Error(lErr)
	}

	return nil
}
