package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Eckle/TheFramework/db/queries"
	"github.com/Eckle/TheFramework/httpcodec"
	"github.com/Eckle/TheFramework/models"
)

type Controller interface {
	GetAll() http.Handler
	Post() http.Handler
	Get() http.Handler
	Patch() http.Handler
	Delete() http.Handler

	AddToRouter(router *http.ServeMux)
}

type BaseController struct {
	Resource models.BaseResource

	// Used for nested controllers
	PreviousController          *BaseController
	PreviousControllerIdMapping string // A string containing the name of the field in the current controller's resource that relates it to the previous controller's resource

	Variable string
}

func (controller BaseController) GetAll() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params, err := queries.ExtractQueryParams(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err)
			return
		}

		resourceList, err := controller.Resource.GetResource(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}

		httpcodec.Encode(w, r, resourceList)
	})
}

func (controller BaseController) Post() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resource, err := httpcodec.Decode(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err)
			return
		}

		err = controller.Resource.CreateResource(&resource)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err)
			return
		}

		httpcodec.Encode(w, r, resource)
	})
}

func (controller BaseController) Get() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resourceId, err := strconv.Atoi(r.PathValue(controller.Variable))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err)
			return
		}

		params, err := queries.ExtractQueryParams(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err)
			return
		}

		queries.AddFilter(params, "id", resourceId)

		resourceList, err := controller.Resource.GetResource(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}

		httpcodec.Encode(w, r, resourceList)
	})
}

func (controller BaseController) Patch() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resourceId, err := strconv.Atoi(r.PathValue(controller.Variable))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err)
			return
		}

		resource, err := httpcodec.Decode(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err)
			return
		}

		params, err := queries.ExtractQueryParams(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err)
			return
		}

		err = controller.Resource.UpdateResource(resourceId, params, &resource)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err)
			return
		}

		httpcodec.Encode(w, r, resource)
	})
}

func (controller BaseController) Delete() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resourceId, err := strconv.Atoi(r.PathValue(controller.Variable))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err)
			return
		}

		err = controller.Resource.DeleteResource(resourceId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

func (controller BaseController) AddToRouter(router *http.ServeMux) {
	router.Handle(fmt.Sprintf("GET /%s", controller.Resource.Table), controller.GetAll())
	router.Handle(fmt.Sprintf("POST /%s", controller.Resource.Table), controller.Post())
	router.Handle(fmt.Sprintf("GET /%s/{%s}", controller.Resource.Table, controller.Variable), controller.Get())
	router.Handle(fmt.Sprintf("PATCH /%s/{%s}", controller.Resource.Table, controller.Variable), controller.Patch())
	router.Handle(fmt.Sprintf("DELETE /%s/{%s}", controller.Resource.Table, controller.Variable), controller.Delete())
}

func New(resource models.BaseResource, previous_controller *BaseController, previous_controller_id_mapping string) BaseController {
	return BaseController{
		Resource:                    resource,
		PreviousController:          previous_controller,
		PreviousControllerIdMapping: previous_controller_id_mapping,
		Variable:                    fmt.Sprintf("%sId", strings.TrimSuffix(resource.Table, "s")),
	}
}
