// Package version is a convenience utility that provides SDK
// consumers with a ready-to-use version command that
// produces apps versioning information based on flags
// passed at compile time.
//
// Configure the version command
//
// The version command can be just added to your cobra root command.
// At build time, the variables Name, Version, Commit, and BuildTags
// can be passed as build flags as shown in the following example:
//
//  go build -X github.com/cosmos/cosmos-sdk/version.Name=gaia \
//   -X github.com/cosmos/cosmos-sdk/version.ServerName=gaiad \
//   -X github.com/cosmos/cosmos-sdk/version.ClientName=gaiacli \
//   -X github.com/cosmos/cosmos-sdk/version.Version=1.0 \
//   -X github.com/cosmos/cosmos-sdk/version.Commit=f0f7b7dab7e36c20b757cebce0e8f4fc5b95de60 \
//   -X "github.com/okex/exchain/libs/cosmos-sdk/version.BuildTags=linux darwin amd64"
package version

import (
	"encoding/json"
	"fmt"
	"runtime"
	"runtime/debug"
)

var (
	// application's name
	Name = ""
	// server binary name
	ServerName = "<appd>"
	// client binary name
	ClientName = "<appcli>"
	// application's version string
	Version = ""
	// commit
	Commit = ""
	// build tags
	BuildTags = ""
	CosmosSDK  = ""
	Tendermint = ""
)

// Info defines the application version information.
type Info struct {
	Name       string     `json:"name" yaml:"name"`
	ServerName string     `json:"server_name" yaml:"server_name"`
	ClientName string     `json:"client_name" yaml:"client_name"`
	Version    string     `json:"version" yaml:"version"`
	GitCommit  string     `json:"commit" yaml:"commit"`
	BuildTags  string     `json:"build_tags" yaml:"build_tags"`
	GoVersion  string     `json:"go" yaml:"go"`
	BuildDeps  []buildDep `json:"build_deps" yaml:"build_deps"`
	CosmosSDK  string `json:"cosmos_sdk" yaml:"cosmos_sdk"`
	Tendermint string `json:"tendermint" yaml:"tendermint"`
}

func NewInfo() Info {
	return Info{
		Name:       Name,
		ServerName: ServerName,
		ClientName: ClientName,
		Version:    Version,
		CosmosSDK:  CosmosSDK,
		Tendermint: Tendermint,
		GitCommit:  Commit,
		BuildTags:  BuildTags,
		GoVersion:  fmt.Sprintf("go version %s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH),
		BuildDeps:  depsFromBuildInfo(),
	}
}

func (v Info) String() string {
	return fmt.Sprintf(`%s: %s
git commit: %s
build tags: %s
cosmos-sdk: %s
tendermint: %s
%s`,v.Name, v.Version, v.GitCommit, v.BuildTags, v.CosmosSDK, v.Tendermint, v.GoVersion)
}

func depsFromBuildInfo() (deps []buildDep) {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return nil
	}

	for _, dep := range buildInfo.Deps {
		deps = append(deps, buildDep{dep})
	}

	return
}

type buildDep struct {
	*debug.Module
}

func (d buildDep) String() string {
	if d.Replace != nil {
		return fmt.Sprintf("%s@%s => %s@%s", d.Path, d.Version, d.Replace.Path, d.Replace.Version)
	}

	return fmt.Sprintf("%s@%s", d.Path, d.Version)
}

func (d buildDep) MarshalJSON() ([]byte, error)      { return json.Marshal(d.String()) }
func (d buildDep) MarshalYAML() (interface{}, error) { return d.String(), nil }
