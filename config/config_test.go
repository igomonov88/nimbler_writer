package config

import (
"testing"
)

// Success and failure markers.
const (
	success = "\u2713"
	failed  = "\u2717"
)

func TestParse(t *testing.T) {
	t.Log("Given we need to test config reading:")
	{
		t.Logf("\tAnd we got file path with valid config file:")
		_, err := Parse("config.yaml")
		if err != nil {
			t.Fatalf("\t%sWhen should be able to parse config.yaml file: %s.", failed, err)
		}
		t.Logf("\t%s\tWhen should be able to parse config.yaml file.", success)
	}
	{
		t.Logf("\tAnd we got not existing file path:")
		_, err := Parse("conf.yaml")
		if err == nil {
			t.Fatalf("\t%s\tWhen should get an error.", failed)
		}
		t.Logf("\t%s\tWhen should get an error.", success)
	}
	{
		t.Logf("\tAnd we try to check provided values from config:")
		cfg, err := Parse("config.yaml")
		if err != nil {
			t.Fatalf("\t%s\tWhen should be able to use values from readed config file: %s.", failed, err)
		}
		{
			if len(cfg.Service.APIHost) == 0 || len(cfg.Zipkin.ReporterURI) == 0 {
				t.Fatalf("\t%s\tWhen nescessary config values should not be empty.", failed)
			}
			t.Logf("\t%s\tWhen nescessary config values should not being empty.", success)
		}
	}
}

