package config

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/upfluence/pkg/cfg"
	"github.com/upfluence/pkg/log"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

type BuildType string

const (
	Go       BuildType = "go"
	Ruby     BuildType = "rb"
	Frontend BuildType = "frontend"
	Python   BuildType = "py"
	Node     BuildType = "node"

	defaultDist = "dist"
)

var (
	defaultTags = map[string]string{
		"master":  "latest",
		"staging": "staging",
	}

	defaultEnvs = map[string]string{
		"master":  "production",
		"staging": "staging",
	}
)

type Configuration struct {
	Verbose    bool      `yaml:"verbose"`
	Type       BuildType `yaml:"type"`
	Repository string    `yaml:"repository,omitempty"`

	Compiler *Compiler `yaml:"compiler,omitempty"`

	Docker   *Docker   `yaml:"docker,omitempty"`
	Deployer *Deployer `yaml:"deployer,omitempty"`
}

func (c Configuration) GetRepo() string {
	if r := c.Repository; r != "" {
		return r
	}

	dir, err := os.Getwd()

	if err != nil {
		log.Errorf("Getwd: %s", err.Error())

		return "."
	}

	for _, path := range append(
		strings.Split(os.Getenv("GOPATH"), ":"),
		strings.Split(os.Getenv("GOROOT"), ":")...,
	) {
		dir = strings.TrimPrefix(dir, path+"/src/")
	}

	return dir
}

func (c Configuration) GetPath() string {
	dir, err := os.Getwd()

	if err != nil {
		log.Errorf("Getwd: %s", err.Error())

		return "."
	}

	return dir
}
func (c Configuration) GetVerbose() bool {
	if c.Verbose {
		return true
	}

	return cfg.FetchBool("VERBOSE", false)
}

func (c Configuration) GetCompiler() Compiler {
	if c.Compiler != nil {
		return *c.Compiler
	}

	return Compiler{}
}

func (c Configuration) GetBuilder() Docker {
	if c.Docker != nil {
		return *c.Docker
	}

	return Docker{}
}

func (c Configuration) GetDeployer() Deployer {
	if c.Deployer != nil {
		return *c.Deployer
	}

	return Deployer{}
}

type Binary struct {
	Path string `yaml:"path"`
	Name string `yaml:"name"`
}

func (b Binary) GetPath() string {
	if b.Path == "" {
		log.Fatalf("Path empty in the binary %+v", b)
	}

	return b.Path
}

func (b Binary) GetName() string {
	if b.Name != "" {
		return b.Name
	}

	splittedPath := strings.Split(b.Path, "/")

	return splittedPath[len(splittedPath)-1]
}

type Compiler struct {
	Binaries []Binary          `yaml:"binaries"`
	Dist     string            `yaml:"dist"`
	CGO      string            `yaml:"cgo"`
	Args     map[string]string `yaml:"args"`
}

func (c Compiler) GetCGO() string {
	if c.CGO == "" {
		return "0"
	}

	return c.CGO
}

func (c Compiler) GetDist() string {
	if c.Dist != "" {
		return c.Dist
	}

	return defaultDist
}

type BuildArg struct {
	Provider string `yaml:"provider"`
	Value    string `yaml:"value"`
}

type Docker struct {
	Dockerfile string              `yaml:"dockerfile"`
	Image      string              `yaml:"image"`
	Tags       map[string]string   `yaml:"tags"`
	BuildArgs  map[string]BuildArg `yaml:"build_args"`
}

func (d Docker) GetBuildArgs() map[string]string {
	res := make(map[string]string)

	for k, ba := range d.BuildArgs {
		switch ba.Provider {
		case "static":
			res[k] = ba.Value
		case "env":
			v := ba.Value

			if v == "" {
				v = k
			}

			res[k] = os.Getenv(v)
		}
	}

	if _, ok := res["GITHUB_TOKEN"]; !ok && os.Getenv("GITHUB_TOKEN") != "" {
		res["GITHUB_TOKEN"] = os.Getenv("GITHUB_TOKEN")
	}

	return res
}

func (d Docker) GetTag(branch string) string {
	if t, ok := d.Tags[branch]; ok {
		return t
	}

	return defaultTags[branch]
}

func (d Docker) GetDockerfile() string {
	if d.Dockerfile != "" {
		return d.Dockerfile
	}

	return "Dockerfile"
}

func (d Docker) GetImage() string {
	if v := strings.Split(d.Image, "/"); len(v) == 2 {
		return d.Image
	}

	if d.Image == "" {
		log.Fatalf("docker: Image name empty")
	}

	return fmt.Sprintf("upfluence/%s", d.Image)
}

type Deployer struct {
	URL  string            `yaml:"url"`
	Envs map[string]string `yaml:"envs"`

	AccessToken string `yaml:"access_token"`
}

func (d *Deployer) Client() *http.Client {
	if d.AccessToken == "" {
		return http.DefaultClient
	}

	return oauth2.NewClient(
		context.Background(),
		oauth2.StaticTokenSource(&oauth2.Token{AccessToken: d.AccessToken}),
	)
}

func (d *Deployer) GetEnv(branch string) string {
	if t, ok := d.Envs[branch]; ok {
		return t
	}

	return defaultEnvs[branch]
}

func ParseFile(path string) (*Configuration, error) {
	var buf, err = ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	cfg := Configuration{}

	if err := yaml.Unmarshal(buf, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
