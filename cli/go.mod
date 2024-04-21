module cli

go 1.22.2

require github.com/apple/pkl-go v0.6.0

require agnostos.com/config v1.0.0

require (
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
)

replace agnostos.com/config => ./config
