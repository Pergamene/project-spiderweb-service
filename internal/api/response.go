package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type responseFormat struct {
	Result interface{} `json:"result,omitempty"`
	Meta   struct {
		HTTPStatus string `json:"httpStatus"`
		Message    string `json:"message,omitempty"`
	} `json:"meta"`
}

// RespondWith responds to the given request with the given responsewriter.
// It also logs information regarding the request and response.
func RespondWith(r *http.Request, w http.ResponseWriter, status int, responseData interface{}, errToLog error) {
	dataWrapper := responseFormat{}
	// @TODO: add in transaction/request ids.
	dataWrapper.Meta.HTTPStatus = fmt.Sprintf("%v - %v", status, http.StatusText(status))
	if errMsg, ok := responseData.(error); ok {
		dataWrapper.Meta.Message = errMsg.Error()
	} else {
		dataWrapper.Result = responseData
	}
	if errToLog != nil {
		// @TODO: we need a better logging paradigm.
		logger, _ := zap.NewProduction()
		defer logger.Sync()
		logger.Info("Response error",
			zap.String("err", errToLog.Error()),
		)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(dataWrapper)
}
