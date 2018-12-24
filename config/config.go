package config

import (
	"errors"
	"io/ioutil"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Warehouse       string
	FsApiToken      string
	Backoff         duration
	BackoffStepsMax int
	CheckInterval   duration
	TmpDir          string
	ListExportLimit int
	GroupFilesByDay bool

	// for debug only; can point to localhost
	ExportURL string

	// aws: s3 + redshift
	S3       S3Config
	Redshift RedshiftConfig

	// gcloud: GCS + BigQuery
	GCS      GCSConfig
	BigQuery BigQueryConfig
}

type S3Config struct {
	Bucket  string
	Region  string
	Timeout duration
	S3Only  bool
}

type RedshiftValidator interface {
	ValidateTableSchema() error
}

type RedshiftConfig struct {
	Host        string
	Port        string
	DB          string
	User        string
	Password    string
	ExportTable string
	SyncTable   string
	TableSchema string
	Credentials string
	VarCharMax  int
	RedshiftValidator
}

type Validator struct {
	TableSchema string
}

func (v Validator) ValidateTableSchema() error {
	if v.TableSchema == "" {
		return errors.New("TableSchema definition missing from Redshift configuration. More information: https://www.hauserdocs.io")
	}
	return nil
}

type GCSConfig struct {
	Bucket  string
	GCSOnly bool
}

type BigQueryConfig struct {
	Project     string
	Dataset     string
	ExportTable string
	SyncTable   string
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

func Load(filename string) (*Config, error) {
	var conf Config

	tomlData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if _, err := toml.Decode(string(tomlData), &conf); err != nil {
		return nil, err
	}

	conf.Redshift.RedshiftValidator = &Validator{
		TableSchema: conf.Redshift.TableSchema,
	}

	return &conf, nil
}
