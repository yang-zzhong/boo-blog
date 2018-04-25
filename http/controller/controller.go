package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Controller struct {
	responseWriter interface{}
}

func NewController(responseWriter interface{}) *Controller {
	controller := new(Controller)
	controller.responseWriter = responseWriter
	return controller
}

func (controller *Controller) ResponseWriter() http.ResponseWriter {
	return controller.responseWriter.(http.ResponseWriter)
}

func (controller *Controller) String(msg string, statusCode int) {
	controller.ResponseWriter().WriteHeader(statusCode)
	io.WriteString(controller.ResponseWriter(), msg)
}

func (controller *Controller) Json(obj interface{}, statusCode int) {
	encoder := json.NewEncoder(controller.ResponseWriter())
	encoder.Encode(obj)
	controller.ResponseWriter().Header().Set("Content-Type", "application/json")
}

func (controller *Controller) InternalError(err error) {
	log.Fatal(err)
	controller.String(err.Error(), 500)
}
