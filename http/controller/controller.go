package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Controller struct {
	Rwriter interface{}
}

func (controller Controller) Writer() http.ResponseWriter {
	return controller.Rwriter.(http.ResponseWriter)
}

func (controller Controller) String(msg string, statusCode int) {
	controller.Writer().WriteHeader(statusCode)
	io.WriteString(controller.Writer(), msg)
}

func (controller Controller) Json(obj interface{}, statusCode int) {
	encoder := json.NewEncoder(controller.Writer())
	encoder.Encode(obj)
	controller.Writer().Header().Set("Content-Type", "application/json")
}

func (controller Controller) InternalError(err error) {
	log.Fatal(err)
	controller.String(err.Error(), 500)
}
