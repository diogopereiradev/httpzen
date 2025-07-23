package logoascii

import "fmt"

func GetLogo(context string) string {
	return fmt.Sprintf(`
    __    __  __                       
   / /_  / /_/ /_____  ____  ___  ____ 
  / __ \/ __/ __/ __ \/_  / / _ \/ __ \
 / / / / /_/ /_/ /_/ / / /_/  __/ / / /
/_/ /_/\__/\__/ .___/ /___/\___/_/%s
             /_/                     
`, context)
}
