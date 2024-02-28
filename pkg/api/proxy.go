package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/riotpot/pkg/proxy"
	srvs "github.com/riotpot/pkg/service"
	"github.com/riotpot/pkg/utils"
	"github.com/riotpot/pkg/validators"
)

// Structures used to serialize data:
type GetProxy struct {
	ID      string      `json:"id" binding:"required" gorm:"primary_key"`
	Port    int         `json:"port"`
	Network string      `json:"network"`
	Status  string      `json:"status"`
	Service *GetService `json:"service"`
}

type PatchProxy struct {
	Port    int         `json:"port"`
	Network string      `json:"network"`
	Status  string      `json:"status"`
	Service *GetService `json:"service"`
}

type CreateProxy struct {
	Port    int    `json:"port" binding:"required"`
	Network string `json:"network" binding:"required"`
}

type ChangeProxyStatus struct {
	Status string `json:"status" binding:"required"`
}

type ChangeProxyPort struct {
	Port int `json:"port" binding:"required"`
}

// Routes
var (

	// General routes for proxies
	proxiesRoutes = []Route{
		// GET and POST proxies
		NewRoute("", "GET", getProxies),
		NewRoute("", "POST", createProxy),
	}

	// Routes to manipulate a proxy
	proxyRoutes = []Route{
		// CRUD operations for each proxy
		NewRoute("", "GET", getProxy),
		NewRoute("", "PATCH", patchProxy),
		NewRoute("", "DELETE", delProxy),
		NewRoute("/port", "POST", changeProxyPort),
		NewRoute("/status", "POST", changeProxyStatus),
	}
)

// Routers
var (
	// Proxies
	ProxiesRouter = NewRouter("proxies/", proxiesRoutes, []Router{ProxyRouter})
	ProxyRouter   = NewRouter(":id/", proxyRoutes, []Router{ServiceRouter})
)

func NewProxy(px proxy.Proxy) *GetProxy {
	serv := NewService(px.GetService())

	return &GetProxy{
		ID:      px.GetID(),
		Port:    px.GetPort(),
		Network: px.GetNetwork().String(),
		Status:  px.IsRunning().String(),
		Service: serv,
	}
}

// GET proxies registered
// Contains a filter to get proxies by port, id, status
func getProxies(ctx *gin.Context) {
	id := ctx.Query("id")
	portStr := ctx.Query("port")
	status := ctx.Query("status")

	switch {
	case id != "":
		getProxiesByID(ctx, id)

	case portStr != "":
		getProxiesByPort(ctx, portStr)

	case status != "":
		getProxiesByStatus(ctx, status)

	default:
		getAllProxies(ctx)
	}
}

func getProxiesByID(ctx *gin.Context, id string) {
	pr, err := proxy.Proxies.GetProxy(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Serialize the proxy and send it as a response
	px := NewProxy(pr)
	ctx.JSON(http.StatusOK, px)
}

func getProxiesByPort(ctx *gin.Context, portStr string) {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid port number"})
		return
	}

	pr, err := proxy.Proxies.GetProxyByPort(port)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Serialize the proxy and send it as a response
	px := NewProxy(pr)
	ctx.JSON(http.StatusOK, px)
}

func getProxiesByStatus(ctx *gin.Context, status string) {
	casted := make([]GetProxy, 0)
	for _, pr := range proxy.Proxies.GetProxyByStatus(status) {
		// Serialize the proxy 
		px := NewProxy(pr)
		casted = append(casted, *px)
	}

	ctx.JSON(http.StatusOK, casted)
}

func getAllProxies(ctx *gin.Context) {
	casted := make([]GetProxy, 0)
	for _, px := range proxy.Proxies.GetProxies() {
		// Serialize the proxy
		pr := NewProxy(px)
		//Append the proxy to the casted
		casted = append(casted, *pr)
	}
	// Set the header and transform the struct to JSON format
	ctx.JSON(http.StatusOK, casted)
}

// POST a proxy by port ":port"
func createProxy(ctx *gin.Context) {
	// Validate the post request to create a new proxy
	var input CreateProxy
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "msg": input})
		return
	}

	nt, err := utils.ParseNetwork(input.Network)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a new proxy
	pe, err := proxy.Proxies.CreateProxy(nt, input.Port)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Serialize the new proxy and return it as a response
	pr := NewProxy(pe)
	ctx.JSON(http.StatusOK, pr)
}

func getProxy(ctx *gin.Context) {
	id := ctx.Param("id")
	pe, err := proxy.Proxies.GetProxy(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Serialize the proxy and send it as a response
	pr := NewProxy(pe)
	ctx.JSON(http.StatusOK, pr)
}

// Can update:
// port, status and service
func patchProxy(ctx *gin.Context) {
	var errors []error

	// Validate the post request to patch the proxy
	var input PatchProxy
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the proxy to update
	id := ctx.Param("id")
	pe, err := proxy.Proxies.GetProxy(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the port and the service
	validServ, err := srvs.Services.GetService(input.Service.ID)
	if err != nil {
		errors = append(errors, err)
	}

	/*
		[9/5/2022] TODO: Find a way to update the proxy using a buffer copy, and update every
		field slowly.

		validPort, err := validators.ValidatePort(input.Port)
		if err != nil {
			errors = append(errors, err)
		}
		// Update the Port
		pe.SetPort(validPort)
	*/

	// If there are errors in the list, send a message to the client and return
	if len(errors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Update the service
	pe.SetService(validServ)

	// Serialize the proxy and send it as a response
	pr := NewProxy(pe)
	ctx.JSON(http.StatusOK, pr)
}

// DELETE registered proxy
func delProxy(ctx *gin.Context) {
	id := ctx.Param("id")

	err := proxy.Proxies.DeleteProxy(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Serialize the proxy and send it as a response
	ctx.JSON(http.StatusOK, gin.H{"success": "Proxy deleted"})
}

// POST request to change the status of the proxy
func changeProxyStatus(ctx *gin.Context) {
	// Validate the post request to patch the proxy
	var input ChangeProxyStatus
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the proxy to update
	id := ctx.Param("id")
	pe, err := proxy.Proxies.GetProxy(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Attempt to change the status
	status, err := utils.ParseStatus(input.Status)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch status {
	case utils.RunningStatus:
		err = pe.Start()
	case utils.StoppedStatus:
		pe.Stop()
	default:
		err = fmt.Errorf("status not allowed")
	}

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Serialize the status and send it as the response
	ctx.JSON(http.StatusOK, gin.H{"status": pe.IsRunning().String()})
}

// POST request to change the port of the proxy
func changeProxyPort(ctx *gin.Context) {
	// Validate the post request to update the port
	var input ChangeProxyPort
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the proxy to update
	id := ctx.Param("id")
	pe, err := proxy.Proxies.GetProxy(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validPort, err := validators.ValidatePort(input.Port)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the Port
	pe.SetPort(validPort)

	// Serialize the proxy and send it as a response
	pr := NewProxy(pe)
	ctx.JSON(http.StatusOK, pr)
}