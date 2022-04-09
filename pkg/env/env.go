package env

import (
	"flag"
	"fmt"
	"strings"
)

var (
	active Environment
	dev    Environment = &environment{value: "dev"}
	test   Environment = &environment{value: "test"}
	uat    Environment = &environment{value: "uat"}
	live   Environment = &environment{value: "live"}
)

// Environment used for env variable setting
type Environment interface {
	Value() string
	IsDev() bool
	IsTest() bool
	IsUat() bool
	IsLive() bool
	t()
}

type environment struct {
	value string
}

func (e *environment) Value() string {
	return e.value
}

func (e *environment) IsDev() bool {
	return e.value == "dev"
}

func (e *environment) IsTest() bool {
	return e.value == "test"
}

func (e *environment) IsUat() bool {
	return e.value == "uat"
}

func (e *environment) IsLive() bool {
	return e.value == "live"
}

func (e *environment) t() {}

func init() {
	env := flag.String("env", "", "please input operation env:\n dev: dev in local \n test: test env\n uat:user acceptance testing env\n live: product live env\n")
	flag.Parsed()

	switch strings.ToUpper(strings.TrimSpace(*env)) {
	case "dev":
		active = dev
	case "test":
		active = test
	case "uat":
		active = uat
	case "live":
		active = live
	default:
		active = test
		fmt.Println("Warning: '-env' not set, will use default env: 'test'")
	}
}

// Active return current env info
func Active() Environment {
	return active
}
