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
	Pools []PoolConfig `json:"pools,omitempty"`

	manager *manager
}

type PoolConfig struct {
	Name       string         `json:"name,omitempty"`
	Worker     string         `json:"worker,omitempty"`
	NumThreads int            `json:"num_threads,omitempty"`
	MaxWait    caddy.Duration `json:"max_wait,omitempty"`
}

func (Pogo) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "pogo",
		New: func() caddy.Module { return new(Pogo) },
	}
}

func (p *Pogo) Provision(_ caddy.Context) error {
	if err := validatePoolConfigs(p.Pools); err != nil {
		return err
	}

	pools := make(map[string]*pool, len(p.Pools))
	for _, cfg := range p.Pools {
		maxWait := time.Duration(cfg.MaxWait)
		if maxWait <= 0 {
			maxWait = 30 * time.Second
		}

		workers := frankenphpCaddy.RegisterWorkers("m#Pogo/"+cfg.Name, cfg.Worker, cfg.NumThreads)
		pools[cfg.Name] = newPool(cfg.Name, workers, maxWait)
	}

	p.manager = newManager(pools)

	globalManagerMu.Lock()
	globalManager = p.manager
	globalManagerMu.Unlock()

	return nil
}

func validatePoolConfigs(configs []PoolConfig) error {
	if len(configs) == 0 {
		return fmt.Errorf("pogo requires at least one pool")
	}

	seen := make(map[string]struct{}, len(configs))
	for _, cfg := range configs {
		if cfg.Name == "" {
			return fmt.Errorf("pogo pool name is required")
		}
		if cfg.Worker == "" {
			return fmt.Errorf("pogo pool %q worker is required", cfg.Name)
		}
		if _, exists := seen[cfg.Name]; exists {
			return fmt.Errorf("duplicate pogo pool %q", cfg.Name)
		}
		seen[cfg.Name] = struct{}{}
	}

	if _, ok := seen[defaultPoolName]; !ok {
		return fmt.Errorf("pogo requires a %q pool", defaultPoolName)
	}

	return nil
}

func (p *Pogo) Cleanup() error {
	globalManagerMu.Lock()
	if globalManager == p.manager {
		globalManager = nil
	}
	globalManagerMu.Unlock()

	if p.manager != nil {
		p.manager.close()
	}

	return nil
}

func (p *Pogo) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		for d.NextBlock(0) {
			switch d.Val() {
			case "pool":
				cfg, err := unmarshalPool(d)
				if err != nil {
					return err
				}
				p.Pools = append(p.Pools, cfg)
			default:
				return d.Errf(`unrecognized subdirective "%s"`, d.Val())
			}
		}
	}

	return nil
}

func unmarshalPool(d *caddyfile.Dispenser) (PoolConfig, error) {
	cfg := PoolConfig{}

	if !d.NextArg() {
		return cfg, d.ArgErr()
	}
	cfg.Name = d.Val()

	if d.NextArg() {
		return cfg, d.Errf(`too many arguments for "pool": %s`, d.Val())
	}

	for d.NextBlock(1) {
		switch d.Val() {
		case "worker":
			if !d.NextArg() {
				return cfg, d.ArgErr()
			}
			cfg.Worker = d.Val()
		case "num_threads":
			if !d.NextArg() {
				return cfg, d.ArgErr()
			}
			n, err := strconv.Atoi(d.Val())
			if err != nil {
				return cfg, d.WrapErr(err)
			}
			cfg.NumThreads = n
		case "max_wait":
			if !d.NextArg() {
				return cfg, d.ArgErr()
			}
			duration, err := caddy.ParseDuration(d.Val())
			if err != nil {
				return cfg, d.WrapErr(err)
			}
			cfg.MaxWait = caddy.Duration(duration)
		default:
			return cfg, d.Errf(`unrecognized pool subdirective "%s"`, d.Val())
		}
	}

	return cfg, nil
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
