package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/riotpot/pkg/api"
	"github.com/riotpot/pkg/utils"
)

func SetupRouter() *gin.Engine {
	// Create a router
	router := gin.Default()
	group := router.Group("/api/")
	// Add the proxy routes
	api.ProxiesRouter.AddToGroup(group)
	api.ServicesRouter.AddToGroup(group)

	return router
}

func TestApiProxy(t *testing.T) {

	expected := &api.CreateProxy{
		Port:    8080,
		Network: utils.TCP.String(),
	}

	router := SetupRouter()
	w := httptest.NewRecorder()

	// POST request to create a new proxy
	body, _ := json.Marshal(expected)
	req, _ := http.NewRequest("POST", "/api/proxies/", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)
	response, _ := ioutil.ReadAll(w.Body)

	// Assert the body of the created proxy is equal to the response
	outputPost := &api.CreateProxy{}
	json.Unmarshal(response, outputPost)
	assert.Equal(t, expected, outputPost)

	// GET all the proxies
	req, _ = http.NewRequest("GET", "/api/proxies/", nil)
	router.ServeHTTP(w, req)
	response, _ = ioutil.ReadAll(w.Body)

	// Assert we got 1 proxy in total
	outputGet := &[]api.CreateProxy{}
	json.Unmarshal(response, outputGet)
	assert.Equal(t, 1, len(*outputGet))
}

func TestApiService(t *testing.T) {

	expected := &api.CreateService{
		Name:    "Test Service",
		Host:    "localhost",
		Port:    8080,
		Network: utils.TCP.String(),
	}

	router := SetupRouter()
	w := httptest.NewRecorder()

	// POST to create a new service
	body, _ := json.Marshal(expected)
	req, _ := http.NewRequest("POST", "/api/services/", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)
	response, _ := ioutil.ReadAll(w.Body)

	// Assert the body of the created service is equal to the response
	outputPost := &api.CreateService{}
	json.Unmarshal(response, outputPost)
	assert.Equal(t, expected, outputPost)

	// Request all services
	req, _ = http.NewRequest("GET", "/api/services/", nil)
	router.ServeHTTP(w, req)
	response, _ = ioutil.ReadAll(w.Body)

	outputGet := &[]api.CreateService{}
	json.Unmarshal(response, outputGet)
	assert.Equal(t, 1, len(*outputGet))
}
