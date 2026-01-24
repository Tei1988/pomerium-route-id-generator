package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCLI runs integration tests for the pomerium-route-id-gen CLI tool.
func TestCLI(t *testing.T) {
	// Build binary
	cmd := exec.Command("go", "build", "-o", "pomerium-route-id-gen", "./cmd/pomerium-route-id-gen")
	err := cmd.Run()
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	defer os.Remove("pomerium-route-id-gen")

	t.Run("basic text output", func(t *testing.T) {
		configFile := "test-basic.yaml"
		content := `
policy:
  - from: https://a.example.com
    to: http://httpbin.org
`
		err := os.WriteFile(configFile, []byte(content), 0644)
		assert.NoError(t, err)
		defer os.Remove(configFile)

		out, err := exec.Command("./pomerium-route-id-gen", "--config", configFile).Output()
		assert.NoError(t, err)
		assert.Contains(t, string(out), "RouteID: fd1bc82bb90632f6")
		assert.Contains(t, string(out), "From: https://a.example.com")
	})

	t.Run("basic json output", func(t *testing.T) {
		configFile := "test-json.yaml"
		content := `
policy:
  - from: https://a.example.com
    to: http://httpbin.org
`
		err := os.WriteFile(configFile, []byte(content), 0644)
		assert.NoError(t, err)
		defer os.Remove(configFile)

		out, err := exec.Command("./pomerium-route-id-gen", "--config", configFile, "--output", "json").Output()
		assert.NoError(t, err)

		var results []map[string]string
		err = json.Unmarshal(out, &results)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "fd1bc82bb90632f6", results[0]["route_id"])
		assert.Equal(t, "https://a.example.com", results[0]["from"])
	})

	t.Run("both policy and routes", func(t *testing.T) {
		configFile := "test-both.yaml"
		content := `
policy:
  - from: https://policy.example.com
    to: http://httpbin.org
routes:
  - from: https://route.example.com
    to: http://httpbin.org
`
		err := os.WriteFile(configFile, []byte(content), 0644)
		assert.NoError(t, err)
		defer os.Remove(configFile)

		out, err := exec.Command("./pomerium-route-id-gen", "--config", configFile).Output()
		assert.NoError(t, err)
		assert.Contains(t, string(out), "From: https://policy.example.com")
		assert.Contains(t, string(out), "From: https://route.example.com")
		assert.Contains(t, string(out), "---")
	})
}
