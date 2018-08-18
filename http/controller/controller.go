package controller

import (
	httprouter "github.com/yang-zzhong/go-httprouter"
	"log"
)

type Controller struct {
	w *httprouter.ResponseWriter
}

func NewController(responseWriter *httprouter.ResponseWriter) *Controller {
	controller := new(Controller)
	controller.w = responseWriter
	return controller
}

func (controller *Controller) ResponseWriter() *httprouter.ResponseWriter {
	return controller.w
}

func (controller *Controller) String(msg string, statusCode int) {
	controller.w.WithStatusCode(statusCode).String(msg)
}

func (controller *Controller) Json(obj interface{}, statusCode int) {
	controller.w.WithStatusCode(statusCode).WithHeader("Content-Type", "application/json").Json(obj)
}

func (controller *Controller) InternalError(err error) {
	log.Fatal(err)
	controller.w.InternalError(err)
}
