module github.com/azghr/forge/examples/cli

go 1.26.1

require (
	github.com/azghr/forge/envconfig v0.0.0
	github.com/azghr/forge/flagsub v0.0.0
	github.com/azghr/forge/option v0.0.0
	github.com/azghr/forge/regexcache v0.0.0
	github.com/azghr/forge/shellquote v0.0.0
	github.com/azghr/forge/sliceutil v0.0.0
	github.com/azghr/forge/tablewriter v0.0.0
)

replace (
	github.com/azghr/forge/envconfig => ../../envconfig
	github.com/azghr/forge/flagsub => ../../flagsub
	github.com/azghr/forge/option => ../../option
	github.com/azghr/forge/regexcache => ../../regexcache
	github.com/azghr/forge/shellquote => ../../shellquote
	github.com/azghr/forge/sliceutil => ../../sliceutil
	github.com/azghr/forge/tablewriter => ../../tablewriter
)
