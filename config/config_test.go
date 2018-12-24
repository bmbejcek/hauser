package config

import (
	"testing"

)

func makeValidator(tableSchema string) *Validator {
	return &Validator{
		TableSchema: tableSchema,
	}
}

func TestValidateTableSchemaConfig(t *testing.T) {

	testCases := []struct {
		validator       *Validator
		hasError   bool
		errMessage string
	}{
		{
			validator:       makeValidator(""),
			hasError:   true,
			errMessage: "TableSchema definition missing from Redshift configuration. More information: https://www.hauserdocs.io",
		},
		{
			validator:       makeValidator("test"),
			hasError:   false,
			errMessage: "",
		},
		{
			validator:       makeValidator("search_path"),
			hasError:   false,
			errMessage: "",
		},
	}

	for _, tc := range testCases {
		err := tc.validator.ValidateTableSchema()
		if tc.hasError && err == nil {
			t.Errorf("expected Redshift.ValidateTableSchema() to return an error when config.Config.Redshift.TableSchema is empty")
		}
		if tc.hasError && err.Error() != tc.errMessage {
			t.Errorf("expected Redshift.ValidateTableSchema() to return \n%s \nwhen config.Config.Redshift.TableSchema is empty, returned \n%s \ninstead", tc.errMessage, err)
		}
		if !tc.hasError && err != nil {
			t.Errorf("unexpected error thrown for TableSchema %s: %s", tc.validator.TableSchema, err)
		}
	}
}