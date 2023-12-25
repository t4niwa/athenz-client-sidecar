package usecase

import (
	"crypto/rand"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/kpango/glg"
)

// TODO: remove this testing function
// func testingCache(c *clientd, report func()) {
func testing(c *clientd) {
	// go func() {
	// 	glg.Warn("🌟pprof server start~")
	// 	http.ListenAndServe("localhost:8080", nil)
	// }()
	go func() {

		time.Sleep(5 * time.Second)
		addDummy := func(startIndex int) {
			glg.Warn("🌟Generating fake tokens for testing ...")
			for i := startIndex; i < startIndex+20000; i++ {
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
			// f, _ := os.Create(fmt.Sprintf("/tmp/pprof/profile.pb.%07d.gz", startIndex))
			// defer f.Close()
			// runtime.GC()
			// pprof.WriteHeapProfile(f) // heap dump to file
			// report()                  // log cache size
		}

		// 20k each call for each type, to 100k token (total 200k at+rt)
		addDummy(0)
		addDummy(20000)
		addDummy(40000)
		addDummy(60000)
		addDummy(80000)
		addDummy(100000)
		// addDummy(120000)
		// addDummy(140000)
		// addDummy(160000)
		// addDummy(180000)

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
