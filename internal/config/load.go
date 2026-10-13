package config

import "github.com/hashicorp/hcl/v2/hclsimple"

func LoadAgentConfig(path string) (AgentConfig, error) {
	var cfg AgentConfig
	if err := hclsimple.DecodeFile(path, nil, &cfg); err != nil {
		return AgentConfig{}, err
	}
	if err := cfg.Validate(); err != nil {
		return AgentConfig{}, err
	}
	return cfg, nil
}

func LoadAggregatorConfig(path string) (AggregatorConfig, error) {
	var cfg AggregatorConfig
	if err := hclsimple.DecodeFile(path, nil, &cfg); err != nil {
		return AggregatorConfig{}, err
	}
	if err := cfg.Validate(); err != nil {
		return AggregatorConfig{}, err
	}
	return cfg, nil
}
