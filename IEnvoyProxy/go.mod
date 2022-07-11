module IEnvoyProxy

go 1.16

replace (
	git.torproject.org/pluggable-transports/snowflake.git/v2 => ../snowflake
	github.com/lucas-clemente/quic-go => github.com/tobyxdd/quic-go v0.27.1-0.20220516000630-9265b64059b0
	github.com/pion/dtls/v2 => github.com/pion/dtls/v2 v2.0.12
	github.com/tobyxdd/hysteria v1.0.5 => ../hysteria
	github.com/v2fly/v2ray-core v4.15.0+incompatible => ../v2ray-core
	gitlab.com/yawning/obfs4.git => ../obfs4
	www.bamsoftware.com/git/dnstt.git => ../dnstt
)

require (
	git.torproject.org/pluggable-transports/snowflake.git/v2 v2.2.0
	github.com/gopherjs/gopherjs v0.0.0-20210420193930-a4630ec28c79 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/lucas-clemente/quic-go v0.27.2 // indirect
	github.com/miekg/dns v1.1.49 // indirect
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	github.com/songgao/water v0.0.0-20200317203138-2b4b6d7c09d8 // indirect
	github.com/tobyxdd/hysteria v1.0.5
	github.com/v2fly/v2ray-core v4.15.0+incompatible // indirect
	github.com/v2fly/v2ray-core/v5 v5.0.7 // indirect
	gitlab.com/yawning/obfs4.git v0.0.0-20220204003609-77af0cba934d
	golang.org/x/tools v0.1.8-0.20211022200916-316ba0b74098 // indirect
	www.bamsoftware.com/git/dnstt.git v0.0.0-00010101000000-000000000000
)
