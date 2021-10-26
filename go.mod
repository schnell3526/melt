module github.com/schnell3526/melt

go 1.17

require (
	github.com/dsnet/compress v0.0.1
	github.com/klauspost/compress v1.13.6
	github.com/spf13/cobra v1.2.1
	github.com/ulikunitz/xz v0.5.10
)

require (
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)

replace github.com/schnell3526/melt/cmd => ./cmd

replace github.com/schnell3526/melt/util/prefix/zip => ./util/prefix/zip
