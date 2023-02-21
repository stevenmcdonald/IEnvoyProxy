package IEnvoyProxy

import (
	"fmt"
	"encoding/json"
	"log"
	"net"
	"os"
	"strconv"
	"time"
	dnsttclient "www.bamsoftware.com/git/dnstt.git/dnstt-client"
	hysteria "github.com/tobyxdd/hysteria/cmd"
	v2ray "github.com/v2fly/v2ray-core/envoy"
)

var dnsttPort = 57000

// DnsttPort - Port where Dnstt will provide its service.
// Only use this property after calling StartDnstt! It might have changed after that!
//
//goland:noinspection GoUnusedExportedFunction
func DnsttPort() int {
	return dnsttPort
}

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

func V2raySrtpPort() int {
	return v2raySrtpPort
}

func V2rayWechatPort() int {
	return v2rayWechatPort
}

func V2rayWsPort() int {
	return v2rayWsPort
}

var dnsttRunning = false
var hysteriaRunning = false
var v2rayWsRunning = false
var v2raySrtpRunning = false
var v2rayWechatRunning = false


// StartDnstt - Start the Dnstt client.
//
// @param ttDomain	subdomain name for DNSTT
//
// @param dohURL OPTIONAL. URL of a DoH resolver. Use either this or `dotAddr`.
//
// @param dotAddr OPTIONAL. Address of a DoT resolver. Use either this or `dohURL`.
//
// @param pubkey The DNSTT's server public key (as hex digits).
//
// @return Port number where Dnstt will listen on, if no error happens during start up.
//
//goland:noinspection GoUnusedExportedFunction
func StartDnstt(ttDomain, dohURL, dotAddr, pubkey string) int {
	if dnsttRunning {
		return dnsttPort
	}

	dnsttRunning = true

	dnsttPort = findPort(dnsttPort)

	// From the dnstt docs:
	//
	// In -doh and -dot modes, the program's TLS fingerprint is camouflaged with
	// uTLS by default. The specific TLS fingerprint is selected randomly from a
	// weighted distribution. You can set your own distribution (or specific single
	// fingerprint) using the -utls option. The special value "none" disables uTLS.
	//     -utls '3*Firefox,2*Chrome,1*iOS'
	//     -utls Firefox
	//     -utls none
	var utlsDistribution = "3*Firefox,1*iOS"
	var listenAddr = fmt.Sprintf("localhost:%d", dnsttPort)

	go dnsttclient.Start(&ttDomain, &listenAddr, &dohURL, &dotAddr, &pubkey, &utlsDistribution)

	return dnsttPort
}

// StopDnstt - Stop the Dnstt client.
//
//goland:noinspection GoUnusedExportedFunction
func StopDnstt() {
	if !dnsttRunning {
		return
	}

	go dnsttclient.Stop()

	dnsttRunning = false
}

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
	if hysteriaRunning {
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

	return hysteriaPort
}

func StopHysteria() {
	if !hysteriaRunning {
		return
	}

	go hysteria.Stop()

	hysteriaRunning = false
}

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
	if v2raySrtpRunning {
		return v2raySrtpPort
	}

	v2raySrtpPort = findPort(v2raySrtpPort)
	clientPort := strconv.Itoa(v2raySrtpPort)

	v2raySrtpRunning = true

	go v2ray.StartSrtp(&clientPort, &serverAddress, &serverPort, &id)

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
	if v2rayWechatRunning {
		return v2rayWechatPort
	}

	v2rayWechatPort = findPort(v2rayWechatPort)
	clientPort := strconv.Itoa(v2rayWechatPort)

	v2rayWechatRunning = true

	go v2ray.StartWechat(&clientPort, &serverAddress, &serverPort, &id)

	return v2rayWechatPort
}

func StopV2RayWechat() {
	if !v2rayWechatRunning {
		return
	}

	go v2ray.StopWechat()

	v2rayWechatRunning = false
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
