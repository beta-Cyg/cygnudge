module cygnus.beta/cygnudge/server

go 1.21.5

require (
	cygnus.beta/cygnudge v0.0.1
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/google/uuid v1.5.0
)

require (
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.30.0 // indirect
)

replace cygnus.beta/cygnudge v0.0.1 => ../
