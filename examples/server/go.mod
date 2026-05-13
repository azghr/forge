module github.com/azghr/forge/examples/server

go 1.26.1

require (
	github.com/azghr/forge/atomicfile v0.0.0
	github.com/azghr/forge/envconfig v0.0.0
	github.com/azghr/forge/flagsub v0.0.0
	github.com/azghr/forge/pathsafe v0.0.0
	github.com/azghr/forge/stopwatch v0.0.0
	github.com/azghr/forge/stringutil v0.0.0
	github.com/azghr/forge/validator v0.0.0
)

replace (
	github.com/azghr/forge/atomicfile => ../../atomicfile
	github.com/azghr/forge/envconfig => ../../envconfig
	github.com/azghr/forge/flagsub => ../../flagsub
	github.com/azghr/forge/pathsafe => ../../pathsafe
	github.com/azghr/forge/stopwatch => ../../stopwatch
	github.com/azghr/forge/stringutil => ../../stringutil
	github.com/azghr/forge/validator => ../../validator
)
