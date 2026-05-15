package pogo

import (
	"testing"

	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

func parsePogoConfig(t *testing.T, input string) (*Pogo, error) {
	t.Helper()

	d := caddyfile.NewTestDispenser(input)
	p := &Pogo{}
	err := p.UnmarshalCaddyfile(d)
	return p, err
}

func TestUnmarshalMultiplePools(t *testing.T) {
	p, err := parsePogoConfig(t, `pogo {
		pool default {
			worker public/pogo-worker.php
			num_threads 2
			max_wait 5s
		}
		pool external_api {
			worker public/api-worker.php
			num_threads 7
			max_wait 10s
		}
	}`)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if len(p.Pools) != 2 {
		t.Fatalf("expected 2 pools, got %d", len(p.Pools))
	}
	if p.Pools[0].Name != defaultPoolName || p.Pools[1].Name != "external_api" {
		t.Fatalf("unexpected pool names: %#v", p.Pools)
	}
	if p.Pools[1].Worker != "public/api-worker.php" || p.Pools[1].NumThreads != 7 {
		t.Fatalf("unexpected external_api config: %#v", p.Pools[1])
	}
}

func TestValidatePoolConfigsRequiresDefault(t *testing.T) {
	err := validatePoolConfigs([]PoolConfig{{
		Name:   "external_api",
		Worker: "public/pogo-worker.php",
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

func TestValidatePoolConfigsRequiresWorker(t *testing.T) {
	err := validatePoolConfigs([]PoolConfig{{
		Name: defaultPoolName,
	}})
	if err == nil {
		t.Fatal("expected missing worker error")
	}
}
