module github.com/azghr/forge/examples/config

go 1.26.1

require (
	github.com/azghr/forge/atomicfile v0.0.0
	github.com/azghr/forge/jsonmerge v0.0.0
	github.com/azghr/forge/pathsafe v0.0.0
)

replace (
	github.com/azghr/forge/atomicfile => ../../atomicfile
	github.com/azghr/forge/jsonmerge => ../../jsonmerge
	github.com/azghr/forge/pathsafe => ../../pathsafe
)
