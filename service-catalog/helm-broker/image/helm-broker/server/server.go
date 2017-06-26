package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kubernetes-incubator/service-catalog/pkg/brokerapi"
	"github.com/kubernetes-incubator/service-catalog/pkg/util"
	"helm-broker/controller"
	"log"
	"net/http"
	"strconv"
)

type server struct {
	controller controller.Controller
}

// CreateHandler creates Broker HTTP handler based on an implementation
// of a controller.Controller interface.
func CreateHandler(c controller.Controller) http.Handler {
	s := server{
		controller: c,
	}

	var router = mux.NewRouter()

	router.HandleFunc("/v2/catalog", s.catalog).Methods("GET")
	router.HandleFunc("/v2/service_instances/{instance_id}/last_operation", s.getServiceInstance).Methods("GET")
	router.HandleFunc("/v2/service_instances/{instance_id}", s.createServiceInstance).Methods("PUT")
	router.HandleFunc("/v2/service_instances/{instance_id}", s.removeServiceInstance).Methods("DELETE")
	router.HandleFunc("/v2/service_instances/{instance_id}/service_bindings/{binding_id}", s.bind).Methods("PUT")
	router.HandleFunc("/v2/service_instances/{instance_id}/service_bindings/{binding_id}", s.unBind).Methods("DELETE")

	return router
}

// Start creates the HTTP handler based on an implementation of a
// controller.Controller interface, and begins to listen on the specified port.
func Start(serverPort int, c controller.Controller) {
	log.Printf("Starting server on %d\n", serverPort)
	http.Handle("/", CreateHandler(c))
	if err := http.ListenAndServe(":"+strconv.Itoa(serverPort), nil); err != nil {
		panic(err)
	}
}

func (s *server) catalog(w http.ResponseWriter, r *http.Request) {
	log.Printf("Get Service Broker Catalog...")

	if result, err := s.controller.Catalog(); err == nil {
		util.WriteResponse(w, http.StatusOK, result)
	} else {
		util.WriteResponse(w, http.StatusBadRequest, err)
	}
}

func (s *server) getServiceInstance(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["instance_id"]
	log.Printf("GetServiceInstance ... %s\n", id)

	result, err := s.controller.GetServiceInstance(id)
	log.Println(result.State)
	if err != nil {
		util.WriteResponse(w, http.StatusBadRequest, err)
	}
	if result.State == brokerapi.StateInProgress {
		util.WriteResponse(w, http.StatusAccepted, result)
	} else if result.State == brokerapi.StateSucceeded {
		util.WriteResponse(w, http.StatusOK, result)
	} else {
		util.WriteResponse(w, http.StatusBadRequest, err)
	}
}

func (s *server) createServiceInstance(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["instance_id"]
	log.Printf("CreateServiceInstance %s...\n", id)

	var req brokerapi.CreateServiceInstanceRequest
	if err := util.BodyToObject(r, &req); err != nil {
		log.Fatalf("error unmarshalling: %v", err)
		util.WriteResponse(w, http.StatusBadRequest, err)
		return
	}
	if req.Parameters == nil {
		req.Parameters = make(map[string]interface{})
	}

	if result, err := s.controller.CreateServiceInstance(id, &req); err == nil {
		util.WriteResponse(w, http.StatusAccepted, result)
	} else {
		util.WriteResponse(w, http.StatusBadRequest, err)
	}
}

func (s *server) removeServiceInstance(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["instance_id"]
	log.Printf("RemoveServiceInstance %s...\n", id)

	if result, err := s.controller.RemoveServiceInstance(id); err == nil {
		util.WriteResponse(w, http.StatusOK, result)
	} else {
		util.WriteResponse(w, http.StatusBadRequest, err)
	}
}

func (s *server) bind(w http.ResponseWriter, r *http.Request) {
	bindingID := mux.Vars(r)["binding_id"]
	instanceID := mux.Vars(r)["instance_id"]

	log.Printf("Bind binding_id=%s, instance_id=%s\n", bindingID, instanceID)

	var req brokerapi.BindingRequest

	if err := util.BodyToObject(r, &req); err != nil {
		log.Printf("Failed to unmarshall request: %v", err)
		util.WriteResponse(w, http.StatusBadRequest, err)
		return
	}
	if req.Parameters == nil {
		req.Parameters = make(map[string]interface{})
	}

	// Pass in the instanceId to the template.
	req.Parameters["instanceId"] = instanceID

	if result, err := s.controller.Bind(instanceID, bindingID, &req); err == nil {
		util.WriteResponse(w, http.StatusOK, result)
	} else {
		util.WriteResponse(w, http.StatusBadRequest, err)
	}
}

func (s *server) unBind(w http.ResponseWriter, r *http.Request) {
	instanceID := mux.Vars(r)["instance_id"]
	bindingID := mux.Vars(r)["binding_id"]
	log.Printf("UnBind: Service instance guid: %s:%s", bindingID, instanceID)

	if err := s.controller.UnBind(instanceID, bindingID); err == nil {
		w.WriteHeader(http.StatusOK)
		fmt.Print(w, "{}") //id)
	} else {
		util.WriteResponse(w, http.StatusBadRequest, err)
	}
}
