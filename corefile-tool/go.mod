module github.com/coredns/corefile-migration/corefile-tool

go 1.14

replace github.com/coredns/corefile-migration => ../

require (
	github.com/coredns/corefile-migration v0.0.0
	github.com/lithammer/dedent v1.1.0
	github.com/spf13/cobra v1.8.1
)
