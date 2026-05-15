package pogo

import (
	"fmt"
	"strconv"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	frankenphpCaddy "github.com/dunglas/frankenphp/caddy"
)

func init() {
	caddy.RegisterModule(Pogo{})
	httpcaddyfile.RegisterGlobalOption("pogo", parseGlobalOption)
}

type Pogo struct {
	Worker     string         `json:"worker,omitempty"`
	NumThreads int            `json:"num_threads,omitempty"`
	MaxWait    caddy.Duration `json:"max_wait,omitempty"`

	pool *pool
}

func (Pogo) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "pogo",
		New: func() caddy.Module { return new(Pogo) },
	}
}

func (p *Pogo) Provision(_ caddy.Context) error {
	if p.Worker == "" {
		return fmt.Errorf("pogo worker is required")
	}

	maxWait := time.Duration(p.MaxWait)
	if maxWait <= 0 {
		maxWait = 30 * time.Second
	}

	workers := frankenphpCaddy.RegisterWorkers("m#Pogo", p.Worker, p.NumThreads)
	p.pool = newPool(workers, maxWait)

	globalPoolMu.Lock()
	globalPool = p.pool
	globalPoolMu.Unlock()

	return nil
}

func (p *Pogo) Cleanup() error {
	globalPoolMu.Lock()
	if globalPool == p.pool {
		globalPool = nil
	}
	globalPoolMu.Unlock()

	if p.pool != nil {
		p.pool.close()
	}

	return nil
}

func (p *Pogo) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		for d.NextBlock(0) {
			switch d.Val() {
			case "worker":
				if !d.NextArg() {
					return d.ArgErr()
				}
				p.Worker = d.Val()
			case "num_threads":
				if !d.NextArg() {
					return d.ArgErr()
				}
				n, err := strconv.Atoi(d.Val())
				if err != nil {
					return d.WrapErr(err)
				}
				p.NumThreads = n
			case "max_wait":
				if !d.NextArg() {
					return d.ArgErr()
				}
				duration, err := caddy.ParseDuration(d.Val())
				if err != nil {
					return d.WrapErr(err)
				}
				p.MaxWait = caddy.Duration(duration)
			default:
				return d.Errf(`unrecognized subdirective "%s"`, d.Val())
			}
		}
	}

	return nil
}

func parseGlobalOption(d *caddyfile.Dispenser, _ any) (any, error) {
	app := &Pogo{}
	if err := app.UnmarshalCaddyfile(d); err != nil {
		return nil, err
	}

	return httpcaddyfile.App{
		Name:  "pogo",
		Value: caddyconfig.JSON(app, nil),
	}, nil
}

var (
	_ caddy.Module       = (*Pogo)(nil)
	_ caddy.Provisioner  = (*Pogo)(nil)
	_ caddy.CleanerUpper = (*Pogo)(nil)
)
