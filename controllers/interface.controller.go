package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Eckle/TheFramework/db/queries"
	"github.com/Eckle/TheFramework/httpcodec"
	"github.com/Eckle/TheFramework/models"
)

type Controllers interface {
	GetAll() http.Handler
	Post() http.Handler
	Get() http.Handler
	Patch() http.Handler
	Delete() http.Handler

	AddToRouter(router *http.ServeMux)
}

type BaseControllers struct {
	Resource models.BaseResource

	// Used for nested controllers
	PreviousController          *BaseControllers
	PreviousControllerIdMapping string // A string containing the name of the field in the current controller's resource that relates it to the previous controller's resource

	Variable string
}

func (controller BaseControllers) GetAll() http.Handler {
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

func (controller BaseControllers) Post() http.Handler {
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

func (controller BaseControllers) Get() http.Handler {
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

func (controller BaseControllers) Patch() http.Handler {
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

func (controller BaseControllers) Delete() http.Handler {
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

func (controller BaseControllers) AddToRouter(router *http.ServeMux) {
	router.Handle(fmt.Sprintf("GET /%s", controller.Resource.Table), controller.GetAll())
	router.Handle(fmt.Sprintf("POST /%s", controller.Resource.Table), controller.Post())
	router.Handle(fmt.Sprintf("GET /%s/{%s}", controller.Resource.Table, controller.Variable), controller.Get())
	router.Handle(fmt.Sprintf("PATCH /%s/{%s}", controller.Resource.Table, controller.Variable), controller.Patch())
	router.Handle(fmt.Sprintf("DELETE /%s/{%s}", controller.Resource.Table, controller.Variable), controller.Delete())
}
