module IEnvoyProxy

go 1.16

replace (
	git.torproject.org/pluggable-transports/snowflake.git/v2 => ../snowflake
	github.com/lucas-clemente/quic-go => github.com/tobyxdd/quic-go v0.27.1-0.20220516000630-9265b64059b0
	github.com/pion/dtls/v2 => github.com/pion/dtls/v2 v2.0.12
	github.com/tobyxdd/hysteria v1.0.5 => ../hysteria
	github.com/v2fly/v2ray-core v4.15.0+incompatible => ../v2ray-core
)

require (
	github.com/lucas-clemente/quic-go v0.27.2 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.17.0 // indirect
	github.com/pion/dtls/v2 v2.1.5 // indirect
	github.com/pion/sctp v1.8.2 // indirect
	github.com/pion/transport v0.13.0 // indirect
	github.com/tobyxdd/hysteria v1.0.5
	github.com/v2fly/v2ray-core v4.15.0+incompatible
	github.com/v2fly/v2ray-core/v5 v5.0.7 // indirect
	golang.org/x/crypto v0.0.0-20220516162934-403b01795ae8 // indirect
	golang.org/x/mobile v0.0.0-20221110043201-43a038452099 // indirect
)
