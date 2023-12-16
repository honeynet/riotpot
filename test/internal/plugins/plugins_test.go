package main

import (
	"testing"

	lr "github.com/riotpot/pkg/logger"
	"github.com/riotpot/pkg/plugins"
	"github.com/riotpot/pkg/service"
	"github.com/stretchr/testify/assert"
)

var (
	pluginPath = "../../../bin/plugins/*.so"
)

func TestLoadPlugins(t *testing.T) {
	pgs, err := plugins.GetPluginServices(pluginPath)
	if err != nil {
		lr.Log.Fatal().Err(err).Msgf("One or more services could not be found")
	}

	assert.Equal(t, 1, len(pgs))

	plg := pgs[0]
	i, ok := plg.(service.PluginService)
	if !ok {
		lr.Log.Fatal().Err(err).Msgf("Service is not a plugin")
	}
	go i.Run()
}

func TestNewPrivateKey(t *testing.T) {
	key := plugins.NewPrivateKey(plugins.DefaultKey)
	pem := key.GetPEM()

	assert.Equal(t, 1, len(pem))
}
