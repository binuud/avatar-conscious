package mcpManagerServer_test

import (
	"context"
	"testing"

	pbMcp "github.com/binuud/avatar-conscious/gen/go/v1/mcpManager"
	"github.com/binuud/avatar-conscious/pkg/servers/mcpManagerServer"
)

const (
	TOTAL_MCP_ADDED = 3
)

func Test_MCP_Manager(t *testing.T) {
	// Step 1: Setup
	t.Run("Step=Register", func(t *testing.T) {
		_register(t)
	})

	// Step 2: Depends on Step 1
	t.Run("Step=Check", func(t *testing.T) {
		_check(t)
	})

	// Step 3: Depends on Step 2
	t.Run("Step=List", func(t *testing.T) {
		_list(t)
	})

	// Step 3: Depends on Step 3
	t.Run("Step=Get", func(t *testing.T) {
		_get(t)
	})

	t.Run("Step=Delete", func(t *testing.T) {
		_delete(t)
	})

}

func _register(t *testing.T) {
	s := mcpManagerServer.NewMCPManagerGrpcServer()

	got, err := s.Register(context.Background(), &pbMcp.RegisterRequest{
		Config: &pbMcp.MCPServerConfig{
			Name:      "get-weather",
			Transport: pbMcp.MCPTransportType_TRANSPORT_TYPE_STDIO,
			Command:   "mcp-weather-server",
		},
	})

	if err != nil {
		t.Errorf("MCP_Server register: cannot register server : %v \n", err)
	} else {
		t.Logf("MCP_Server register: success %v", got)
	}

	got, err = s.Register(context.Background(), &pbMcp.RegisterRequest{
		Config: &pbMcp.MCPServerConfig{
			Name:      "greet-server",
			Transport: pbMcp.MCPTransportType_TRANSPORT_TYPE_STDIO,
			Command:   "mcp-greet-server",
		},
	})

	if err != nil {
		t.Errorf("MCP_Server register: cannot register server : %v \n", err)
	} else {
		t.Logf("MCP_Server register: success %v", got)
	}

	got, err = s.Register(context.Background(), &pbMcp.RegisterRequest{
		Config: &pbMcp.MCPServerConfig{
			Name:      "get-weather",
			Transport: pbMcp.MCPTransportType_TRANSPORT_TYPE_STDIO,
			Command:   "mcp-weather-server",
		},
	})

	if err == nil {
		t.Errorf("error MCP_Server register: able to register duplicate: %v \n", got)
	} else {
		t.Logf("MCP_Server register: success duplicate identified %v", err)
	}

	got, err = s.Register(context.Background(), &pbMcp.RegisterRequest{
		Config: &pbMcp.MCPServerConfig{
			Name:      "mock-server",
			Transport: pbMcp.MCPTransportType_TRANSPORT_TYPE_SSE,
			Url:       "http://localhost:1080/mockserver/mcp",
		},
	})

	if err != nil {
		t.Errorf("MCP_Server register: cannot register remote server : %v \n", err)
	} else {
		t.Logf("MCP_Server register: success %v", got)
	}

}

func _list(t *testing.T) {

	s := mcpManagerServer.NewMCPManagerGrpcServer()

	got, err := s.List(context.Background(), &pbMcp.ListRequest{})
	if err != nil {
		t.Errorf("error MCP_Server list: not able to list: %v \n", err)
	} else {
		if TOTAL_MCP_ADDED == len(got.Servers) {
			t.Logf("MCP_Server register: success able to list - %d mcp servers found", len(got.Servers))
		} else {
			t.Errorf("error expected %d servers but got %d", TOTAL_MCP_ADDED, len(got.Servers))
		}

	}

}

func _check(t *testing.T) {}

func _get(t *testing.T) {}

func _delete(t *testing.T) {}
