package trigger

import (
	"fmt"
	"testing"

	"github.com/asaskevich/EventBus"
)

func TestName(t *testing.T) {
	bus := EventBus.New()
	bus.Subscribe("main:calculator", calculator)
	bus.Publish("main:calculator", 20, 40)
	bus.Unsubscribe("main:calculator", calculator)
}

func calculator(a int, b int) {
	fmt.Printf("%d\n", a+b)
}
