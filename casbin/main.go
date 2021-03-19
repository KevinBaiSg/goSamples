package main

import (
	"fmt"

	"github.com/casbin/casbin"
)

func main() {
	enforcer := casbin.NewEnforcer("/Users/kevin/working/golang/goSamples/casbin/model.conf", "/Users/kevin/working/golang/goSamples/casbin/policy.csv")
	ret := enforcer.AddRoleForUser("kevin", "alice")
	fmt.Println(ret)
	role, err := enforcer.GetUsersForRole("kevin")
	if err != nil {
		return
	}
	fmt.Println(role)
}
