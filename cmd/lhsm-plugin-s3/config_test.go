package main

import (
	"os"
	"reflect"
	"testing"

	"github.intel.com/hpdd/policy/pdm/lhsmd/config"
)

func TestLoadConfig(t *testing.T) {
	loaded, err := loadConfig("./test-fixtures/lhsm-plugin-s3.test")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &s3Config{
		Archives: archiveSet{
			&archiveConfig{
				Name:   "2",
				ID:     2,
				Region: "us-west-1",
				Bucket: "hpdd-test-bucket",
				Prefix: "archive-test",
			},
		},
	}

	if !reflect.DeepEqual(loaded, expected) {
		t.Fatalf("\nexpected: \n\n%#v\ngot: \n\n%#v\n\n", expected, loaded)
	}
}

func TestMergedConfig(t *testing.T) {
	os.Setenv(config.AgentConnEnvVar, "foo://bar:1234")
	os.Setenv(config.PluginMountpointEnvVar, "/foo/bar/baz")
	os.Setenv(config.ConfigDirEnvVar, "./test-fixtures")

	merged, err := getMergedConfig()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &s3Config{
		AgentAddress: "foo://bar:1234",
		ClientRoot:   "/foo/bar/baz",
		Region:       "us-east-1",
		Archives: archiveSet{
			&archiveConfig{
				Name:   "2",
				ID:     2,
				Region: "us-west-1",
				Bucket: "hpdd-test-bucket",
				Prefix: "archive-test",
			},
		},
	}

	if !reflect.DeepEqual(merged, expected) {
		t.Fatalf("\nexpected: \n\n%#v\ngot: \n\n%#v\n\n", expected, merged)
	}
}

func TestArchiveValidation(t *testing.T) {
	loaded, err := loadConfig("./test-fixtures/lhsm-plugin-s3.test")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	for _, archive := range loaded.Archives {
		if err := archive.checkValid(); err != nil {
			t.Fatalf("err: %s", err)
		}
	}

	loaded, err = loadConfig("./test-fixtures/lhsm-plugin-s3-badarchive")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	for _, archive := range loaded.Archives {
		if err := archive.checkValid(); err == nil {
			t.Fatalf("expected %s to fail validation", archive)
		}
	}
}