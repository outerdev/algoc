module github.com/outerdev/algoc

go 1.14

require (
	github.com/algorand/go-algorand v0.0.0-20200320145517-cea8009ba7ee
	github.com/algorand/go-algorand-sdk v1.2.1
	github.com/go-yaml/yaml v2.1.0+incompatible
	github.com/manifoldco/promptui v0.7.0
	github.com/mitchellh/go-ps v1.0.0
	github.com/spf13/cobra v0.0.6
	golang.org/x/crypto v0.0.0-20200302210943-78000ba7a073 // indirect
	golang.org/x/sys v0.0.0-20200202164722-d101bd2416d5
	gopkg.in/yaml.v2 v2.2.4 // indirect
)

replace github.com/algorand/go-algorand => ./go-algorand
