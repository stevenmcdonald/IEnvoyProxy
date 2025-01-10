package v2ray

// copied and modified from main/commands/run.go

// feeding core.LoadConfig with a string reader containing the config
// JSON seems the simplest way to run v2ray as a library
//
// Golang's JSON support seems a little cumbersome for a couple
// string substitutions in a complex JSON snippet, so we just
// use fmt.Sprintf to assemble the config
//
// The JSON file we build should look similar to the client example config
// (that will be) documented here:
// https://gitlab.com/stevenmcdonald/envoy-proxy-examples/v2ray/

// We provide functions for starting and stopping several client services
// independently. Unfortunately you can't start them all and tell v2ray to
// use the one that works... but that's what Envoy is good at.

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	core "github.com/v2fly/v2ray-core/v5"
	_ "github.com/v2fly/v2ray-core/v5/main/distro/all"
)

var osWsSignals = make(chan os.Signal, 1)
var osWechatSignals = make(chan os.Signal, 1)
var osSrtpSignals = make(chan os.Signal, 1)

// getInbound
//
// @param port - port to listen for SOCKS5 connections
func getInbound(clientPort int) string {
	return fmt.Sprintf(`
      {
        "port": %d,
        "protocol": "socks",
        "sniffing": {
          "enabled": true,
          "destOverride": ["http", "tls"]
        },
        "settings": {
          "auth": "noauth"
        }
      }`, clientPort)
}

func getWsConfig(clientPort int, serverAddress, serverWsPort, wsPath, id string) string {
	return fmt.Sprintf(`
  {
    "log": {
      "loglevel": "error"
    },
    "inbounds": [%s
    ],
    "outbounds": [
      {
        "protocol": "vmess",
        "settings": {
          "vnext": [
            {
              "address": "%s",
              "port": %s,
              "users": [
                {
                  "id": "%s",
                  "alterId": 0
                }
              ]
            }
          ]
        },
        "streamSettings": {
          "network": "ws",
          "security": "tls",
          "wsSettings": {
            "path": "%s"
          }
        }
      }
    ]
  }`, getInbound(clientPort), serverAddress, serverWsPort, id, wsPath)
}

// getQUICConfig
//
// @param clientPort - port to listen on for SOCKS5 connections
//
// @param serverAddress - server address to connect to
//
// @param serverPort - server port to connect to
//
// @oaram type - type of QUIC obfuscation, should be "srtp" or "wechat-video"
func getQuicConfig(clientPort int, serverAddress, serverPort, quicType, id string) string {
	return fmt.Sprintf(`
  {
    "log": {
      "loglevel": "error"
    },
    "inbounds": [%s
    ],
    "outbounds": [
      {
        "protocol": "vmess",
        "settings": {
          "vnext": [
            {
              "address": "%s",
              "port": %s,
              "users": [
                {
                  "id": "%s",
                  "alterId": 0
                }
              ]
            }
          ]
        },
        "streamSettings": {
          "network": "quic",
          "quicSettings": {
            "security": "aes-128-gcm",
            "header": {
              "type": "%s"
            },
            "key": "0"
          }
        }
      }
    ]
  }`, getInbound(clientPort), serverAddress, serverPort, id, quicType)
}

func startServer(jsonConfig string) (*core.Instance, error) {
	reader := strings.NewReader(jsonConfig)

	config, err := core.LoadConfig(core.FormatJSON, reader)
	if err != nil {
		fmt.Printf("error reading config: %s\n", err)
		return nil, err
	}

	server, err := core.New(config)
	if err != nil {
		fmt.Printf("error creating server: %s\n", err)
		return nil, err
	}

	if err := server.Start(); err != nil {
		fmt.Printf("failed to start %s\n", err)

		_ = server.Close()

		return nil, err
	}

	return server, nil
}

// StartWs - start v2ray, websocket transport
//
// @param clientPort - client SOCKS port routed to the WS server
//
// @param serverAddress - IP or hostname of the server
//
// @param serverPort - port of the websocket server (probably 443)
//
// @param wsPath - path to the websocket on the server
//
// @param id - UUID used to authenticate with the server
//
// @returns error, if transport could not be started, or `nil` on success.
func StartWs(clientPort int, serverAddress, serverPort, wsPath, id string) error {
	server, err := startServer(getWsConfig(clientPort, serverAddress, serverPort, wsPath, id))
	if err != nil {
		return err
	}

	go func(server *core.Instance) {
		defer func(server *core.Instance) {
			_ = server.Close()
		}(server)

		{
			signal.Notify(osWsSignals, syscall.SIGTERM)
			<-osWsSignals
		}
	}(server)

	return nil
}

func StopWs() {
	osWsSignals <- syscall.SIGTERM
}

// StartSrtp - start v2ray, QUIC/SRTP transport
//
// @param clientPort - client SOCKS port routed to the WS server
//
// @param serverAddress - IP or hostname of the server
//
// @param serverPort - port of the websocket server (probably 443)
//
// @param id - UUID used to authenticate with the server
//
// @returns error, if transport could not be started, or `nil` on success.
func StartSrtp(clientPort int, serverAddress, serverPort, id string) error {
	server, err := startServer(getQuicConfig(clientPort, serverAddress, serverPort, "srtp", id))
	if err != nil {
		return err
	}

	go func(server *core.Instance) {
		defer func(server *core.Instance) {
			_ = server.Close()
		}(server)

		{
			signal.Notify(osSrtpSignals, syscall.SIGTERM)
			<-osSrtpSignals
		}
	}(server)

	return nil
}

func StopSrtp() {
	osSrtpSignals <- syscall.SIGTERM
}

// StartWechat - start v2ray, QUIC/Wechat-video transport
//
// @param clientPort - client SOCKS port routed to the WS server
//
// @param serverAddress - IP or hostname of the server
//
// @param serverPort - port of the websocket server (probably 443)
//
// @param id - UUID used to authenticate with the server
//
// @returns error, if transport could not be started, or `nil` on success.
func StartWechat(clientPort int, serverAddress, serverPort, id string) error {
	server, err := startServer(getQuicConfig(clientPort, serverAddress, serverPort, "wechat-video", id))
	if err != nil {
		return err
	}

	go func(server *core.Instance) {
		defer func(server *core.Instance) {
			_ = server.Close()
		}(server)

		{
			signal.Notify(osWechatSignals, syscall.SIGTERM)
			<-osWechatSignals
		}
	}(server)

	return nil
}

func StopWechat() {
	osWechatSignals <- syscall.SIGTERM
}
