package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/upfluence/pkg/cfg"
	"github.com/upfluence/pkg/log"
	"gopkg.in/yaml.v2"
)

type BuildType string

const (
	Go   BuildType = "go"
	Ruby BuildType = "rb"

	defaultDist = "dist"
)

var defaultTags = map[string]string{
	"master":  "latest",
	"staging": "staging",
}

type Configuration struct {
	Verbose    bool      `yaml:"verbose"`
	Type       BuildType `yaml:"type"`
	Repository string    `yaml:"repository,omitempty"`

	Compiler *Compiler `yaml:"compiler,omitempty"`

	Docker *Docker `yaml:"docker,omitempty"`
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
	Binaries []Binary `yaml:"binaries"`
	Dist     string   `yaml:"dist"`
}

func (c Compiler) GetDist() string {
	if c.Dist != "" {
		return c.Dist
	}

	return defaultDist
}

type Docker struct {
	Dockerfile string            `yaml:"dockerfile"`
	Image      string            `yaml:"image"`
	Tags       map[string]string `yaml:"tags"`
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
