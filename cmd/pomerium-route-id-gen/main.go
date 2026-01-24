package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/pomerium/pomerium/config"
	"gopkg.in/yaml.v3"
)

type OutputFormat string

const (
	OutputFormatText OutputFormat = "text"
	OutputFormatJSON OutputFormat = "json"
)

type RouteInfo struct {
	RouteID string `json:"route_id"`
	From    string `json:"from"`
}

func main() {
	configPath := flag.String("config", "", "Path to Pomerium configuration file")
	outputFormat := flag.String("output", "text", "Output format (text|json)")
	flag.StringVar(outputFormat, "o", "text", "Output format (text|json) (shorthand)")
	flag.Parse()

	if *configPath == "" {
		fmt.Fprintln(os.Stderr, "Error: --config flag is required")
		flag.Usage()
		os.Exit(1)
	}

	data, err := os.ReadFile(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
		os.Exit(1)
	}

	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing YAML: %v\n", err)
		os.Exit(1)
	}

	var options struct {
		Policies []config.Policy `mapstructure:"policy"`
		Routes   []config.Policy `mapstructure:"routes"`
	}

	decodeHook := mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToSliceHookFunc(","),
		config.DecodePolicyHookFunc(),
	)

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:     &options,
		DecodeHook: decodeHook,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating decoder: %v\n", err)
		os.Exit(1)
	}

	if err := decoder.Decode(raw); err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding configuration: %v\n", err)
		os.Exit(1)
	}

	allPolicies := append(options.Policies, options.Routes...)
	var results []RouteInfo

	for i := range allPolicies {
		p := &allPolicies[i]
		id, err := p.RouteID()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to calculate RouteID for %s: %v\n", p.From, err)
			continue
		}
		results = append(results, RouteInfo{
			RouteID: id,
			From:    p.From,
		})
	}

	switch OutputFormat(*outputFormat) {
	case OutputFormatJSON:
		printJSON(results)
	case OutputFormatText:
		printText(results)
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown output format %q\n", *outputFormat)
		os.Exit(1)
	}
}

func printText(results []RouteInfo) {
	for i, res := range results {
		if i > 0 {
			fmt.Println("---")
		}
		fmt.Printf("RouteID: %s\n", res.RouteID)
		fmt.Printf("From: %s\n", res.From)
	}
}

func printJSON(results []RouteInfo) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(results); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}
}
