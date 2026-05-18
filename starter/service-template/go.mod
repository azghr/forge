module github.com/azghr/forge/starter/service-template

go 1.26.1

require (
	github.com/azghr/forge/envconfig v0.0.0
	github.com/azghr/forge/retry v0.0.0
	github.com/azghr/forge/stopwatch v0.0.0
)

replace (
	github.com/azghr/forge/envconfig => ../../envconfig
	github.com/azghr/forge/retry => ../../retry
	github.com/azghr/forge/stopwatch => ../../stopwatch
)
