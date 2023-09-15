package kill_process

import "testing"

func Test_kill(t *testing.T) {

	err := Kill("notepad.exe")
	if err != nil {
		t.Error(err)
	}
}
