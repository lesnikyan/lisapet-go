package nodes

type GuardNo struct {
	msg string
}

type GenFinish struct {
	msg string
}

func (er *GenFinish) Error() string {
	return "Gen is finished"
}

func NewGenFinish() *GenFinish {
	return &GenFinish{}
}
