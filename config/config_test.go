package config

import (
	"testing"

)

func makeConf(tableSchema string) *Config {
	conf := &Config{
		Redshift: RedshiftConfig {
			TableSchema: tableSchema,
		},
	}
	return conf
}

func TestValidateTableSchemaConfig(t *testing.T) {

	testCases := []struct {
		conf       *Config
		hasError   bool
		errMessage string
	}{
		{
			conf:       makeConf(""),
			hasError:   true,
			errMessage: "TableSchema definition missing from Redshift configuration. More information: https://www.hauserdocs.io",
		},
		{
			conf:       makeConf("test"),
			hasError:   false,
			errMessage: "",
		},
		{
			conf:       makeConf("search_path"),
			hasError:   false,
			errMessage: "",
		},
	}

	for _, tc := range testCases {
		err := tc.conf.Redshift.ValidateTableSchema()
		if tc.hasError && err == nil {
			t.Errorf("expected Redshift.ValidateTableSchema() to return an error when config.Config.Redshift.TableSchema is empty")
		}
		if tc.hasError && err.Error() != tc.errMessage {
			t.Errorf("expected Redshift.ValidateTableSchema() to return \n%s \nwhen config.Config.Redshift.TableSchema is empty, returned \n%s \ninstead", tc.errMessage, err)
		}
		if !tc.hasError && err != nil {
			t.Errorf("unexpected error thrown for TableSchema %s: %s", tc.conf.Redshift.TableSchema, err)
		}
	}
}