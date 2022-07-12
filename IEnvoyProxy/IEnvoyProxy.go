package IEnvoyProxy

import (
	"fmt"
	"encoding/json"
	"gitlab.com/yawning/obfs4.git/obfs4proxy"
	"net"
	"os"
	"runtime"
	"strconv"
	"time"
	dnsttclient "www.bamsoftware.com/git/dnstt.git/dnstt-client"
	hysteria "github.com/tobyxdd/hysteria/cmd"
	v2ray "github.com/v2fly/v2ray-core/envoy"
)

var meekPort = 47000

// MeekPort - Port where Obfs4proxy will provide its Meek service.
// Only use this after calling StartObfs4Proxy! It might have changed after that!
//
//goland:noinspection GoUnusedExportedFunction
func MeekPort() int {
	return meekPort
}

var obfs2Port = 47100

// Obfs2Port - Port where Obfs4proxy will provide its Obfs2 service.
// Only use this property after calling StartObfs4Proxy! It might have changed after that!
//
//goland:noinspection GoUnusedExportedFunction
func Obfs2Port() int {
	return obfs2Port
}

var obfs3Port = 47200

// Obfs3Port - Port where Obfs4proxy will provide its Obfs3 service.
// Only use this property after calling StartObfs4Proxy! It might have changed after that!
//
//goland:noinspection GoUnusedExportedFunction
func Obfs3Port() int {
	return obfs3Port
}

var obfs4Port = 47300

// Obfs4Port - Port where Obfs4proxy will provide its Obfs4 service.
// Only use this property after calling StartObfs4Proxy! It might have changed after that!
//
//goland:noinspection GoUnusedExportedFunction
func Obfs4Port() int {
	return obfs4Port
}

var scramblesuitPort = 47400

// ScramblesuitPort - Port where Obfs4proxy will provide its Scramblesuit service.
// Only use this property after calling StartObfs4Proxy! It might have changed after that!
//
//goland:noinspection GoUnusedExportedFunction
func ScramblesuitPort() int {
	return scramblesuitPort
}

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
var v2rayWechatPort = 47601
var v2rayWsPort = 47602

func V2raySrtpPort() int {
	return v2raySrtpPort
}

func V2rayWechatPort() int {
	return v2rayWechatPort
}

func V2rayWsPort() int {
	return v2rayWsPort
}

var obfs4ProxyRunning = false
var dnsttRunning = false
var hysteriaRunning = false
var v2rayRunning = false

// StateLocation - Override TOR_PT_STATE_LOCATION, which defaults to "$TMPDIR/pt_state".
var StateLocation string

func init() {
	if //goland:noinspection GoBoolExpressions
	runtime.GOOS == "android" {
		StateLocation = "/data/local/tmp"
	} else {
		StateLocation = os.Getenv("TMPDIR")
	}

	StateLocation += "/pt_state"
}

// Obfs4ProxyVersion - The version of Obfs4Proxy bundled with IPtProxy.
//
//goland:noinspection GoUnusedExportedFunction
func Obfs4ProxyVersion() string {
	return obfs4proxy.Obfs4proxyVersion
}

// StartObfs4Proxy - Start the Obfs4Proxy.
//
// This will test, if the default ports are available. If not, it will increment them until there is.
// Only use the port properties after calling this, they might have been changed!
//
// @param logLevel Log level (ERROR/WARN/INFO/DEBUG). Defaults to ERROR if empty string.
//
// @param enableLogging Log to TOR_PT_STATE_LOCATION/obfs4proxy.log.
//
// @param unsafeLogging Disable the address scrubber.
//
// @param proxy HTTP, SOCKS4 or SOCKS5 proxy to be used behind Obfs4proxy. E.g. "socks5://127.0.0.1:12345"
//
// @return Port number where Obfs4Proxy will listen on for Obfs4(!), if no error happens during start up.
//	If you need the other ports, check MeekPort, Obfs2Port, Obfs3Port and ScramblesuitPort properties!
//
//goland:noinspection GoUnusedExportedFunction
func StartObfs4Proxy(logLevel string, enableLogging, unsafeLogging bool, proxy string) int {
	if obfs4ProxyRunning {
		return obfs4Port
	}

	obfs4ProxyRunning = true

	for !IsPortAvailable(meekPort) {
		meekPort++
	}

	if meekPort >= obfs2Port {
		obfs2Port = meekPort + 1
	}

	for !IsPortAvailable(obfs2Port) {
		obfs2Port++
	}

	if obfs2Port >= obfs3Port {
		obfs3Port = obfs2Port + 1
	}

	for !IsPortAvailable(obfs3Port) {
		obfs3Port++
	}

	if obfs3Port >= obfs4Port {
		obfs4Port = obfs3Port + 1
	}

	for !IsPortAvailable(obfs4Port) {
		obfs4Port++
	}

	if obfs4Port >= scramblesuitPort {
		scramblesuitPort = obfs4Port + 1
	}

	for !IsPortAvailable(scramblesuitPort) {
		scramblesuitPort++
	}

	fixEnv()

	if len(proxy) > 0 {
		_ = os.Setenv("TOR_PT_PROXY", proxy)
	} else {
		_ = os.Unsetenv("TOR_PT_PROXY")
	}

	go obfs4proxy.Start(&meekPort, &obfs2Port, &obfs3Port, &obfs4Port, &scramblesuitPort, &logLevel, &enableLogging, &unsafeLogging)

	return obfs4Port
}

// StopObfs4Proxy - Stop the Obfs4Proxy.
//
//goland:noinspection GoUnusedExportedFunction
func StopObfs4Proxy() {
	if !obfs4ProxyRunning {
		return
	}

	go obfs4proxy.Stop()

	obfs4ProxyRunning = false
}

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

	fixEnv()

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
	Alpn		string			`json:alpn`
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

	fmt.Printf("config: %s", string(confJson))

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

// StartV2ray -- Start v2ray client
//
// @param serverAddress - IP or hostname of the server for SRTP and wechat-video
//
// @param serverWsAddress - Hostname of WS web server proxy
//
// @param serverWsPort - port of the WS server (TLS is assumed, so probably 443)
//
// @param serverWsPath - websocket path (should be the same in the v2ray config and http proxy host)
//
// @param serverSrtpPort - port for (fake) SRTP connections
//
// @param serverWechatPort - port for (fake) Wechat video connections
//
// @param id - UUID for authentication with the server
//
// returns the client Websocket port, call the helper functions for the other ports
//
func StartV2Ray(serverAddress, serverWsAddress, serverWsPort, serverWsPath, serverSrtpPort, serverWechatPort, id string) int {
	if v2rayRunning {
		return v2rayWsPort
	}

	v2rayWsPort = findPort(v2rayWsPort)
	v2raySrtpPort = findPort(v2rayWsPort + 1)
	v2rayWechatPort = findPort(v2raySrtpPort + 1)

	// convert to strings
	wsport := strconv.Itoa(v2rayWsPort)
	srtpport := strconv.Itoa(v2raySrtpPort)
	wechatport := strconv.Itoa(v2rayWechatPort)

	v2rayRunning = true

	go v2ray.Start(&wsport, &srtpport, &wechatport, &serverAddress, &serverWsPort, &serverSrtpPort, &serverWechatPort, &serverWsPath, &id)

	return v2rayWsPort
}

func StopV2ray() {
	if !v2rayRunning {
		return
	}

	go v2ray.Stop()

	v2rayRunning = false
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

// Hack: Set some environment variables that are either
// required, or values that we want. Have to do this here, since we can only
// launch this in a thread and the manipulation of environment variables
// from within an iOS app won't end up in goptlib properly.
//
// Note: This might be called multiple times when using different functions here,
// but that doesn't necessarily mean, that the values set are independent each
// time this is called. It's still the ENVIRONMENT, we're changing here, so there might
// be race conditions.
func fixEnv() {
	_ = os.Setenv("TOR_PT_CLIENT_TRANSPORTS", "meek_lite,obfs2,obfs3,obfs4,scramblesuit")
	_ = os.Setenv("TOR_PT_MANAGED_TRANSPORT_VER", "1")

	_ = os.Setenv("TOR_PT_STATE_LOCATION", StateLocation)
}
