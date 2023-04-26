package builddeploy

import (
	"time"
)

var config = map[string]struct{
    Timer *time.Timer
    Seconds int
    Branch string
}{}
