package globals

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Create a function getCurrentTime that returns current clock time in this format "17-03-54" where 17 is hour, 03 is minute and 54 is second.
func getCurrentTime() string {
    return time.Now().Format("15-04-05")
}

func GenericRecover () {
    if r := recover(); r != nil {
        log.Println("Recovered in f", r)
        os.WriteFile("panic" + getCurrentTime() + ".txt", []byte(fmt.Sprintf("%v", r)), 0644)
    }
}
