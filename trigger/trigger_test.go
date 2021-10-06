package trigger

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/sadlil/go-trigger"
)

func TestSimpleTrigger(t *testing.T) {
	trigger.On("first-event", func() {
		// Do Some Task Here.
		fmt.Println("Done")
	})
	trigger.Fire("first-event")
}

//using skill
//necessary factor
//target effect  type (damage,defense,cure...),CD,pp-value
//condition
//function
//buff
//

func TestTriggerType(t *testing.T) {
	fmt.Println(reflect.TypeOf(TriggerType_Skill))
}
