package context

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

var (
	routerTitle  = &sync.Map{}
	routerRegexp = regexp.MustCompile(`(.*):[^/]+(.*)`)
)

// SetRouterTitle 设定路由标题
func SetRouterTitle(method, router, title string) {
	routerTitle.Store(fmt.Sprintf("%s-%s", method, router), title)
}

// GetRouterTitleAndKey 获取路由标题和键
func GetRouterTitleAndKey(method, router string) (string, string) {
	key := fmt.Sprintf("%s-%s", method, router)
	vv, ok := routerTitle.Load(key)
	if ok {
		return vv.(string), key
	}

	var title string
	routerTitle.Range(func(vk, vv interface{}) bool {
		vkey := vk.(string)
		if !strings.Contains(vkey, "/:") {
			return true
		}

		rkey := "^" + routerRegexp.ReplaceAllString(vkey, "$1[^/]+$2") + "$"
		b, _ := regexp.MatchString(rkey, key)
		if b {
			title = vv.(string)
			key = vkey
		}
		return !b
	})

	return title, key
}
