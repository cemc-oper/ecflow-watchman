package ecflow_watchman

import "fmt"

var (
	Version   = "Unknown version"
	BuildTime = "Unknown build time"
	GitCommit = "Unknown GitCommit"
)

func PrintVersionInformation() {
	fmt.Printf("Version %s (%s)\n", Version, GitCommit)
	fmt.Printf("Build at %s\n", BuildTime)
	fmt.Printf("Please visit https://github.com/perillaroc/ecflow-watchman for more information.\n")
}
