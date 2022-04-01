package model

type Conf struct {
	Mysql struct {
		Driver   string
		Endpoint string
	}
	ConfAddr string
}
