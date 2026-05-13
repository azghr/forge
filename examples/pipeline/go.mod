module github.com/azghr/forge/examples/pipeline

go 1.26.1

require (
	github.com/azghr/forge/cache v0.0.0
	github.com/azghr/forge/lockutil v0.0.0
	github.com/azghr/forge/mathutil v0.0.0
	github.com/azghr/forge/multityperror v0.0.0
	github.com/azghr/forge/orderedset v0.0.0
	github.com/azghr/forge/priorityqueue v0.0.0
	github.com/azghr/forge/queue v0.0.0
	github.com/azghr/forge/retry v0.0.0
	github.com/azghr/forge/workerpool v0.0.0
)

replace (
	github.com/azghr/forge/cache => ../../cache
	github.com/azghr/forge/lockutil => ../../lockutil
	github.com/azghr/forge/mathutil => ../../mathutil
	github.com/azghr/forge/multityperror => ../../multityperror
	github.com/azghr/forge/orderedset => ../../orderedset
	github.com/azghr/forge/priorityqueue => ../../priorityqueue
	github.com/azghr/forge/queue => ../../queue
	github.com/azghr/forge/retry => ../../retry
	github.com/azghr/forge/workerpool => ../../workerpool
)
