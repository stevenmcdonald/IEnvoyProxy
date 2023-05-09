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
	git.torproject.org/pluggable-transports/snowflake.git/v2 v2.0.0-00010101000000-000000000000
	github.com/lucas-clemente/quic-go v0.27.2 // indirect
	github.com/tobyxdd/hysteria v1.0.5
	github.com/v2fly/v2ray-core v4.15.0+incompatible
	github.com/v2fly/v2ray-core/v5 v5.0.7 // indirect
	golang.org/x/mobile v0.0.0-20230427221453-e8d11dd0ba41 // indirect
	golang.org/x/sync v0.0.0-20220819030929-7fc1605a5dde // indirect
	golang.org/x/tools v0.1.12 // indirect
)
