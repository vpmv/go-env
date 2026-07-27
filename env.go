package env

import (
	"os"

	"github.com/vpmv/go-env/internal/format"
)

var GlobalEnv = `ENV`

// SetDelimiter sets the delimiter used in slice-types.
// Side effect: Sets new string formatter instance
func SetDelimiter(d string) {
	format.Stringer = format.NewStringFormatter(format.WithDelimiter(d))
}

type Environment string

func (e Environment) String() string {
	return string(e)
}

const (
	Development Environment = "development"
	Production  Environment = "production"
	Testing     Environment = "testing"
	Staging     Environment = "staging"
)

var commonEnvs = map[string]Environment{
	`dev`:     Development,
	`devel`:   Development,
	`develop`: Development,
	`prod`:    Production,
	`test`:    Testing,
	`stage`:   Staging,
}

func GetEnv() Environment {
	return Environment(os.Getenv(GlobalEnv))
}

func ParseEnv(env Environment) Environment {
	if environment, ok := commonEnvs[env.String()]; ok {
		return environment
	}
	return env
}

func SetEnv(override bool, env ...Environment) {
	if !override && GetEnv() != `` {
		return
	}

	if len(env) > 0 {
		_ = os.Setenv(GlobalEnv, ParseEnv(env[0]).String())
	} else {
		_ = os.Setenv(GlobalEnv, Development.String())
	}
}

// IsEnv - test if the environment is <environment>
func IsEnv(environment Environment) bool {
	return GetString(GlobalEnv, Development.String()) == environment.String()
}

// IsDevelopment - test if the environment is development
func IsDevelopment() bool {
	return IsEnv(Development)
}

// IsTesting - test if the environment is testing
func IsTesting() bool {
	return IsEnv(Testing)
}

// IsStaging - test if the environment is staging
func IsStaging() bool {
	return IsEnv(Staging)
}

// IsProduction - test if the environment is production
func IsProduction() bool {
	return IsEnv(Production)
}
