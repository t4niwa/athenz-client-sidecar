package usecase

import (
	"crypto/rand"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"sort"
	"strings"
	"time"

	"github.com/kpango/glg"
)

func testing(c *clientd) {
	go func() {

		time.Sleep(30 * time.Second)
		addDummy := func(startIndex, count int) {
			glg.Warn("🌟Generating fake tokens for testing ...")
			for i := startIndex; i < startIndex+count; i++ {
				// generate random token
				key := encode(fmt.Sprintf("domain%d", i), fmt.Sprintf("role%d", i), "proxyForPrincipal")

				atBuf := make([]byte, 750) // average access token size
				rand.Read(atBuf)
				rtBuf := make([]byte, 550) // average role token size
				rand.Read(rtBuf)

				// store token in cache
				c.access.StoreTokenCache(key, string(atBuf), fmt.Sprintf("domain%d", i), fmt.Sprintf("role%d", i))
				c.role.StoreTokenCache(key, string(rtBuf), fmt.Sprintf("domain%d", i), fmt.Sprintf("role%d", i))
			}

			glg.Warn("🌟Generating fake tokens for testing ...END")
			time.Sleep(500 * time.Millisecond)
			f, _ := os.Create(fmt.Sprintf("/tmp/pprof/profile.pb.%07d.gz", startIndex))
			defer f.Close()
			runtime.GC()
			pprof.WriteHeapProfile(f) // heap dump to file
			report(c)                 // log cache size
		}

		// 20k each call for each type, to 100k token (total 200k at+rt)
		addDummy(0, 2000)
		addDummy(2000, 2000)
		addDummy(4000, 2000)
		addDummy(6000, 2000)
		addDummy(8000, 12000)
		addDummy(20000, 20000)
		addDummy(40000, 20000)
		addDummy(60000, 20000)
		addDummy(80000, 120000)
		addDummy(200000, 20000)
		addDummy(220000, 20000)
		addDummy(240000, 20000)
		addDummy(260000, 20000)
		addDummy(280000, 120000)

		// ls | xargs -r -L 1 go tool pprof -inuse_space -top
	}()
}

func encode(domain, role, principal string) string {
	cacheKeySeparator := ";"
	roleSeparator := ","
	roles := strings.Split(role, roleSeparator)
	sort.Strings(roles)

	s := []string{domain, strings.Join(roles, roleSeparator), principal}
	if principal == "" {
		return strings.Join(s[:2], cacheKeySeparator)
	}
	return strings.Join(s, cacheKeySeparator)
}
