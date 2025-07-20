package logoascii

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetLogo(t *testing.T) {
	logo := GetLogo()
	if logo == "" {
		t.Error("Expected non-empty logo")
	}

	assert.Equal(t, `
    __    __  __                       
   / /_  / /_/ /_____  ____  ___  ____ 
  / __ \/ __/ __/ __ \/_  / / _ \/ __ \
 / / / / /_/ /_/ /_/ / / /_/  __/ / / /
/_/ /_/\__/\__/ .___/ /___/\___/_/ .org 
             /_/                     
`, logo, "Logo does not match expected output")
}
