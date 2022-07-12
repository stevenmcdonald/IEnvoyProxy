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
	git.torproject.org/pluggable-transports/goptlib.git v1.1.0 // indirect
	github.com/lucas-clemente/quic-go v0.27.2 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.17.0 // indirect
	github.com/pion/dtls/v2 v2.1.5 // indirect
	github.com/pion/sctp v1.8.2 // indirect
	github.com/pion/transport v0.13.0 // indirect
	github.com/tobyxdd/hysteria v1.0.5
	github.com/v2fly/v2ray-core v4.15.0+incompatible
	github.com/v2fly/v2ray-core/v5 v5.0.7 // indirect
	gitlab.com/yawning/obfs4.git v0.0.0-20220204003609-77af0cba934d
	golang.org/x/crypto v0.0.0-20220516162934-403b01795ae8 // indirect
	golang.org/x/mobile v0.0.0-20220518205345-8578da9835fd // indirect
	golang.org/x/tools v0.1.8-0.20211022200916-316ba0b74098 // indirect
	www.bamsoftware.com/git/dnstt.git v0.0.0-00010101000000-000000000000
)
