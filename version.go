package ecflow_watchman

import "fmt"

var (
	Version    = "Unknown version"
	BuildTime  = "Unknown build time"
	GitCommit  = "Unknown GitCommit"
	ProjectUrl = "https://github.com/perillaroc/ecflow-watchman"
)

func PrintVersionInformation() {
	fmt.Printf("Version %s (%s)\n", Version, GitCommit)
	fmt.Printf("Build at %s\n", BuildTime)
	fmt.Printf("Please visit %s for more information.\n", ProjectUrl)
}
