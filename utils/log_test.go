package utils

import (
	"fmt"
	"os"
	"stringinator-go/constants"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_openLogFiles(t *testing.T) {

	_, err := openLogFiles()
	if err != nil {
		fmt.Errorf("openLogFiles() error = %v", err)
		return
	}

	if _, err := os.Stat(constants.LogFile); err == nil {
		fmt.Printf("File exists\n")
		assert.True(t, true)
	} else {
		fmt.Printf("File does not exist\n")
		assert.True(t, false)
	}
}

func TestConfigureLogs(t *testing.T) {

	if got := ConfigureLogs("level"); got == nil {
		t.Errorf("Failed to cofnigure logs")
	}

	if got := ConfigureLogs("info"); got == nil {
		t.Errorf("Failed to cofnigure logs")
	}
}
