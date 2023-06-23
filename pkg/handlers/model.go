package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"withNats/pkg/model"
	"withNats/pkg/errorsForProject"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ModelHandler struct {
	Logger         *zap.SugaredLogger
	ModelRepo       model.NatsDataRepo
}

func (h *ModelHandler) Find(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nData, err := h.ModelRepo.FindNatsData(vars["ID"])


	if err != nil {
		JsonError(w, http.StatusBadRequest, "Find: "+err.Error(), h.Logger)
		return
	}

	SendRequest(w, "Find: ", nData, http.StatusOK, h.Logger)

	h.Logger.Infof("Find: %v", nData.OrderUID)
	return
}


func SendRequest(w http.ResponseWriter, errStr string, nData *model.NatsData, status int, Logger *zap.SugaredLogger) {
	resp, err := json.Marshal(nData)
	if err != nil {
		JsonError(w, http.StatusBadRequest, errStr+errorsforproject.ErrCantMarshal.Error(), Logger)
		return
	}
	w.WriteHeader(status)
	w.Write(resp)
}

func JsonError(w io.Writer, status int, msg string, Logger *zap.SugaredLogger) {
	resp, err := json.Marshal(map[string]interface{}{
		"status": status,
		"error":  msg,
	})

	if err != nil {
		w.Write([]byte("bad request"))
		return
	}

	w.Write(resp)
}