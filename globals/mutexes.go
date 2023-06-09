package globals

import "sync"

var Get_commits_rw_mutex = &sync.RWMutex{}
