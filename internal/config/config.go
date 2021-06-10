package config

//init log
//init database
//kafka

type Cfg struct {
	DatabaseName string
	DatabasePwd  string
}

func InitConfig() Cfg {
	return Cfg{
		DatabaseName: "TEST",
		DatabasePwd:  "1234",
	}
}
