package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/uveaceslav/affise/model"
	"github.com/uveaceslav/affise/model/validator"
	"github.com/uveaceslav/affise/parser"
)

type mainController struct {
	logger                 *log.Logger
	requestEntityValidator validator.RequestEntityValidator
	urlProcessor           parser.URLProcessor
}

func NewMainController(
	logger *log.Logger,
	requestEntityValidator validator.RequestEntityValidator,
	urlProcessor parser.URLProcessor,
) MainController {
	return &mainController{
		logger:                 logger,
		requestEntityValidator: requestEntityValidator,
		urlProcessor:           urlProcessor,
	}
}

func (mc *mainController) Index(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ToDo: Change this code. The following code is deprecated
	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		mc.logger.Println("The client closed the connection prematurely.")
		cancel()
	}()

	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var requestEntity model.RequestEntity
	if err := json.NewDecoder(r.Body).Decode(&requestEntity); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := mc.requestEntityValidator.Validate(requestEntity); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	processedURLs, err := mc.urlProcessor.Process(ctx, cancel, requestEntity.URLs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(model.ResponseEntity{Data: processedURLs})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
