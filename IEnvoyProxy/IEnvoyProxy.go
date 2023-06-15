package IEnvoyProxy

import (
	"fmt"
	"encoding/json"
	"log"
	"net"
	"strconv"
	"time"
	"os"

	hysteria "github.com/apernet/hysteria/app/cmd"
	v2ray "github.com/v2fly/v2ray-core/envoy"
	snowflakeclient "git.torproject.org/pluggable-transports/snowflake.git/v2/client"
)

var hysteriaPort = 47500

// HysteriaPort - Port where Hysteria will provide its service.
// Only use this property after calling StartHysteria! It might have changed after that!
//
//goland:noinspection GoUnusedExportedFunction
func HysteriaPort() int {
	return hysteriaPort
}


var v2raySrtpPort = 47600
var v2rayWechatPort = 47700
var v2rayWsPort = 47800
var snowflakePort = 47900

func V2raySrtpPort() int {
	return v2raySrtpPort
}

func V2rayWechatPort() int {
	return v2rayWechatPort
}

func V2rayWsPort() int {
	return v2rayWsPort
}

// SnowflakePort - Port where Snowflake will provide its service.
// Only use this property after calling StartSnowflake! It might have changed after that!
//
//goland:noinspection GoUnusedExportedFunction
func SnowflakePort() int {
	return snowflakePort
}

var hysteriaRunning = false
var v2rayWsRunning = false
var v2raySrtpRunning = false
var v2rayWechatRunning = false
var snowflakeRunning = false

/// Hysteria

type HysteriaListen struct {
	Listen string `json:"listen"`
}

type HysteriaConfig struct {
	Server		string			`json:"server"`
	Protocol	string			`json:"protocol"`
	Obfs		string			`json:"obfs"`
	Socks5		HysteriaListen	`json:"socks5"`
	Up_mbps		int				`json:"up_mbps"`
	Down_mbps	int				`json:"down_mbps"`
	Ca			string			`json:"ca"`
	Alpn		string			`json:"alpn"`
}

// StartHysteria -- Start the Hysteria client
//
// @param server Hysteria server hostname or IP and port, e.g. "192.168.64.2:32323"
//
// @param obfs Essentially a password, used to obfuscate the connection,
// MUST use the same value on client and server
//
// @param ca Path to Root CA used by server (for self signed certs)
func StartHysteria(server, obfs, ca string) int {
	log.Println("Starting Hysteria")
	if hysteriaRunning {
		log.Printf("Hysteria already running on %d", hysteriaPort)
		return hysteriaPort
	}

	hysteriaRunning = true

	hysteriaPort = findPort(hysteriaPort)

	// Hysteria uses a JSON file for config, creating JSON
	// to pass in seems like the path of least resistance
	listenAddr := fmt.Sprintf("127.0.0.1:%d", hysteriaPort)

	listenConf := HysteriaListen{listenAddr}
	conf := HysteriaConfig{
		server,
		"wechat-video",
		obfs,
		listenConf,
		1000, // up_mbps
		1000, // down_mbps
		ca,
		"Envoy",
	}

	confJson, err := json.Marshal(conf)

	if err != nil {
		fmt.Println(err)
		return 0
	}

	// fmt.Printf("config: %s", string(confJson))

	go hysteria.Start(&confJson)
	log.Printf("Hysteria started on port %d", hysteriaPort)

	return hysteriaPort
}

func StopHysteria() {
	if !hysteriaRunning {
		return
	}

	go hysteria.Stop()

	hysteriaRunning = false
}

/// V2Ray

// StartV2RayWs - Start V2Ray client for websocket transport
//
// @param serverAddress - Hostname of WS web server proxy
//
// @oaram serverPort - Port of the WS listener (probably 443)
//
// @param wsPath - path the websocket
//
// @param id - v2ray UUID for auth
func StartV2RayWs(serverAddress, serverPort, wsPath, id string) int {
	if v2rayWsRunning {
		return v2rayWsPort
	}

	v2rayWsPort = findPort(v2rayWsPort)
	clientPort := strconv.Itoa(v2rayWsPort)

	v2rayWsRunning = true

	go v2ray.StartWs(&clientPort, &serverAddress, &serverPort, &wsPath, &id)

	return v2rayWsPort
}

func StopV2RayWs() {
	if !v2rayWsRunning {
		return
	}

	go v2ray.StopWs()

	v2rayWsRunning = false
}

func StartV2raySrtp(serverAddress, serverPort, id string) int {
	log.Println("Starting V2Ray SRTP")
	if v2raySrtpRunning {
		log.Printf("V2Ray SRTP already running on %d", v2raySrtpPort)
		return v2raySrtpPort
	}

	v2raySrtpPort = findPort(v2raySrtpPort)
	clientPort := strconv.Itoa(v2raySrtpPort)

	v2raySrtpRunning = true

	go v2ray.StartSrtp(&clientPort, &serverAddress, &serverPort, &id)
	log.Printf("V2Ray SRTP started on %d", v2raySrtpPort)

	return v2raySrtpPort
}

func StopV2RaySrtp() {
	if !v2raySrtpRunning {
		return
	}

	go v2ray.StopSrtp()

	v2raySrtpRunning = false
}

func StartV2RayWechat(serverAddress, serverPort, id string) int {
	log.Println("Starting V2Ray WeChat")
	if v2rayWechatRunning {
		log.Printf("V2Ray WeChat already running on %d", v2rayWechatPort)
		return v2rayWechatPort
	}

	v2rayWechatPort = findPort(v2rayWechatPort)
	clientPort := strconv.Itoa(v2rayWechatPort)

	v2rayWechatRunning = true

	go v2ray.StartWechat(&clientPort, &serverAddress, &serverPort, &id)
	log.Printf("V2Ray WeChat started on %d", v2rayWechatPort)

	return v2rayWechatPort
}

func StopV2RayWechat() {
	if !v2rayWechatRunning {
		return
	}

	go v2ray.StopWechat()

	v2rayWechatRunning = false
}

/// Snowflake

// StartSnowflake - Start the Snowflake client.
//
// @param ice Comma-separated list of ICE servers.
//
// @param url URL of signaling broker.
//
// @param front Front domain.
//
// @param ampCache OPTIONAL. URL of AMP cache to use as a proxy for signaling.
//
//	Only needed when you want to do the rendezvous over AMP instead of a domain fronted server.
//
// @param logFile Name of log file. OPTIONAL. Defaults to no log.
//
// @param logToStateDir Resolve the log file relative to Tor's PT state dir.
//
// @param keepLocalAddresses Keep local LAN address ICE candidates.
//
// @param unsafeLogging Prevent logs from being scrubbed.
//
// @param maxPeers Capacity for number of multiplexed WebRTC peers. DEFAULTs to 1 if less than that.
//
// @return Port number where Snowflake will listen on, if no error happens during start up.
//
//goland:noinspection GoUnusedExportedFunction
func StartSnowflake(ice, url, front, ampCache, logFile string, logToStateDir, keepLocalAddresses, unsafeLogging bool, maxPeers int) int {
	if snowflakeRunning {
		return snowflakePort
	}

	snowflakeRunning = true

	for !IsPortAvailable(snowflakePort) {
		snowflakePort++
	}

	fixEnv()

	go snowflakeclient.Start(&snowflakePort, &ice, &url, &front, &ampCache, &logFile, &logToStateDir, &keepLocalAddresses, &unsafeLogging, &maxPeers)

	return snowflakePort
}

// StopSnowflake - Stop the Snowflake client.
//
//goland:noinspection GoUnusedExportedFunction
func StopSnowflake() {
	if !snowflakeRunning {
		return
	}

	go snowflakeclient.Stop()

	snowflakeRunning = false
}

// SnowflakeClientConnected - Interface to use when clients connect
// to the snowflake proxy. For use with StartSnowflakeProxy
type SnowflakeClientConnected interface {
	// Connected - callback method to handle snowflake proxy client connections.
	Connected()
}


///////////////////
// Helper functions

// in IPtProxy, this handles the PT state directoy stuff...
// we only have snowflake for now, and that only needs a couple env
// vars set.
func fixEnv() {
	_ = os.Setenv("TOR_PT_CLIENT_TRANSPORTS", "snowflake")
	_ = os.Setenv("TOR_PT_MANAGED_TRANSPORT_VER", "1")
}


func findPort(port int) int {
	temp := port
	for !IsPortAvailable(temp) {
		temp++
	}
	return temp
}

// IsPortAvailable - Checks to see if a given port is not in use.
//
// @param port The port to check.
func IsPortAvailable(port int) bool {
	address := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))

	conn, err := net.DialTimeout("tcp", address, 500*time.Millisecond)

	if err != nil {
		return true
	}

	_ = conn.Close()

	return false
}
