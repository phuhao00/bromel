package bromel

type CallBackFn func(component Component)

type Operation struct {
	IsAsynchronous bool
	CB             CallBackFn
	Ret            chan interface{}
}
