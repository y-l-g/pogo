package pogo

import (
	"testing"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

func parsePogoConfig(t *testing.T, input string) (*Pogo, error) {
	t.Helper()

	d := caddyfile.NewTestDispenser(input)
	p := &Pogo{}
	err := p.UnmarshalCaddyfile(d)
	return p, err
}

func TestUnmarshalDefaultWorkerAndNamedPool(t *testing.T) {
	p, err := parsePogoConfig(t, `pogo {
		worker worker.php
		num_threads 2
		max_wait 5s

		pool external_api {
			worker public/api-worker.php
			num_threads 7
			max_wait 10s
		}
	}`)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if p.Worker != "worker.php" || p.NumThreads != 2 {
		t.Fatalf("unexpected default config: %#v", p)
	}
	if len(p.Pools) != 1 || p.Pools[0].Name != "external_api" {
		t.Fatalf("unexpected pools: %#v", p.Pools)
	}
	if p.Pools[0].Worker != "public/api-worker.php" || p.Pools[0].NumThreads != 7 {
		t.Fatalf("unexpected external_api config: %#v", p.Pools[0])
	}
}

func TestPoolConfigsBuildDefaultPoolFromTopLevelWorker(t *testing.T) {
	p := &Pogo{
		Worker:     "worker.php",
		NumThreads: 2,
		Pools: []PoolConfig{{
			Name:   "external_api",
			Worker: "public/api-worker.php",
		}},
	}

	configs, err := p.poolConfigs()
	if err != nil {
		t.Fatalf("pool config failed: %v", err)
	}

	if len(configs) != 2 {
		t.Fatalf("expected 2 pool configs, got %d", len(configs))
	}
	if configs[0].Name != defaultPoolName || configs[0].Worker != "worker.php" {
		t.Fatalf("unexpected default pool config: %#v", configs[0])
	}
}

func TestValidatePoolConfigsRequiresDefault(t *testing.T) {
	err := validatePoolConfigs([]PoolConfig{{
		Name:   "external_api",
		Worker: "worker.php",
	}})
	if err == nil {
		t.Fatal("expected missing default pool error")
	}
}

func TestValidatePoolConfigsRejectsDuplicatePool(t *testing.T) {
	err := validatePoolConfigs([]PoolConfig{
		{Name: defaultPoolName, Worker: "a.php"},
		{Name: defaultPoolName, Worker: "b.php"},
	})
	if err == nil {
		t.Fatal("expected duplicate pool error")
	}
}

func TestPoolConfigsRejectsDuplicateDefaultPool(t *testing.T) {
	err := (&Pogo{
		Worker: "public/default-worker.php",
		Pools: []PoolConfig{{
			Name:   defaultPoolName,
			Worker: "public/other-worker.php",
		}},
	}).Provision(caddy.Context{})
	if err == nil {
		t.Fatal("expected duplicate default pool error")
	}
}

func TestValidatePoolConfigsRequiresWorker(t *testing.T) {
	err := validatePoolConfigs([]PoolConfig{{
		Name: defaultPoolName,
	}})
	if err == nil {
		t.Fatal("expected missing worker error")
	}
}
