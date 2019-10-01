package config

type LockfileUpdates struct {
	Enabled  *bool    `mapstructure:"enabled,omitempty" yaml:"enabled,omitempty" json:"enabled,omitempty"`
	Settings Settings `mapstructure:"settings,omitempty" yaml:"settings,omitempty" json:"settings"`
}
