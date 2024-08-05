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

		if controller.PreviousController != nil && controller.PreviousControllerIdMapping != "" {
			previousControllerResourceId, err := strconv.Atoi(r.PathValue(controller.Variable))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, err)
				return
			}

			queries.AddFilter(params, controller.PreviousControllerIdMapping, previousControllerResourceId)
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

		if controller.PreviousController != nil && controller.PreviousControllerIdMapping != "" {
			previousControllerResourceId, err := strconv.Atoi(r.PathValue(controller.Variable))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, err)
				return
			}

			resource[controller.PreviousControllerIdMapping] = previousControllerResourceId
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

		if controller.PreviousController != nil && controller.PreviousControllerIdMapping != "" {
			previousControllerResourceId, err := strconv.Atoi(r.PathValue(controller.Variable))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, err)
				return
			}

			queries.AddFilter(params, controller.PreviousControllerIdMapping, previousControllerResourceId)
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

		if controller.PreviousController != nil && controller.PreviousControllerIdMapping != "" {
			previousControllerResourceId, err := strconv.Atoi(r.PathValue(controller.Variable))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, err)
				return
			}

			queries.AddFilter(params, controller.PreviousControllerIdMapping, previousControllerResourceId)
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

		params, err := queries.ExtractQueryParams(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err)
			return
		}

		if controller.PreviousController != nil && controller.PreviousControllerIdMapping != "" {
			previousControllerResourceId, err := strconv.Atoi(r.PathValue(controller.Variable))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, err)
				return
			}

			queries.AddFilter(params, controller.PreviousControllerIdMapping, previousControllerResourceId)
		}

		queries.AddFilter(params, "id", resourceId)

		err = controller.Resource.DeleteResource(resourceId, params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

func (controller BaseController) AddToRouter(router *http.ServeMux) {
	var root string = ""
	if controller.PreviousController != nil {
		root = fmt.Sprintf("/%s/{%s}", controller.PreviousController.Resource.Table, controller.PreviousController.Variable)
	}

	router.Handle(fmt.Sprintf("GET %s/%s", root, controller.Resource.Table), controller.GetAll())
	router.Handle(fmt.Sprintf("POST %s/%s", root, controller.Resource.Table), controller.Post())
	router.Handle(fmt.Sprintf("GET %s/%s/{%s}", root, controller.Resource.Table, controller.Variable), controller.Get())
	router.Handle(fmt.Sprintf("PATCH %s/%s/{%s}", root, controller.Resource.Table, controller.Variable), controller.Patch())
	router.Handle(fmt.Sprintf("DELETE %s/%s/{%s}", root, controller.Resource.Table, controller.Variable), controller.Delete())
}

func New(model models.BaseResource, previous_controller *BaseController, previous_controller_id_mapping string) BaseController {
	return BaseController{
		Resource:                    model,
		PreviousController:          previous_controller,
		PreviousControllerIdMapping: previous_controller_id_mapping,
		Variable:                    fmt.Sprintf("%sId", strings.TrimSuffix(model.Table, "s")),
	}
}
