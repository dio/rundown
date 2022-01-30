package archives

import (
	"path/filepath"

	"github.com/codeclysm/extract"
)

type Archive interface {
	Version() string
	BinaryName() string
	URLPattern() string
	Renamer() extract.Renamer
}

type Proxy struct {
	VersionUsed string
}

func (p *Proxy) Version() string {
	if p.VersionUsed != "" {
		return p.VersionUsed
	}
	return "1.21.0"
}

func (p *Proxy) BinaryName() string {
	return "envoy"
}

func (p *Proxy) URLPattern() string {
	return "https://archive.tetratelabs.io/envoy/download/v%s/envoy-v%s-%s-amd64.tar.xz"
}

func (p *Proxy) Renamer() extract.Renamer {
	return func(name string) string {
		baseName := filepath.Base(name)
		if baseName == p.BinaryName() {
			return baseName
		}
		return name
	}
}

type ExtAuthz struct {
	VersionUsed string
}

func (e *ExtAuthz) Version() string {
	if e.VersionUsed != "" {
		return e.VersionUsed
	}
	return "0.6.0-rc0"
}

func (e *ExtAuthz) BinaryName() string {
	return "auth_server"
}

func (e *ExtAuthz) URLPattern() string {
	return "https://github.com/dio/authservice/releases/download/v%s/auth_server_%s_%s_amd64.tar.gz"
}

func (e *ExtAuthz) Renamer() extract.Renamer {
	return func(name string) string {
		if name == e.BinaryName()+".stripped" {
			return e.BinaryName()
		}
		return name
	}
}
