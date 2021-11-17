package config

type Options struct {
	CpuProfile          bool
	TestData            bool
	Wipe                bool
	YesToAll            bool
	Port                uint
	Addr                string
	Sql                 string
	SqlUsername         string
	SqlPassword         string
	SqlAddress          string
	ActivityLogLoc      string
	AdminHiddenPassword string
	AdminPagesDisabled  bool
	NoRobots            bool
	NoSitemap           bool
	LogFileName         string
	AutoCertDomain      string
}

type option struct {
	label       string
	defaultVal  interface{}
	description string
}

func register(ptr interface{}) {
}

func Resolve() Options {

	return Options{}
}
