package flag

import (
	"context"
	"errors"
	"fmt"

	"github.com/xmdhs/gomclauncher/auth"
	"github.com/xmdhs/gomclauncher/download"
	"github.com/xmdhs/gomclauncher/lang"
)

func (f *Flag) Authlib() {
	if f.Outmsg {
		err := download.Downauthlib(context.TODO(), func(s string) { fmt.Println(s) })
		if err != nil {
			fmt.Println(lang.Lang("authlibdownloadfailed"))
			panic(err)
		}
	}
	api, err := auth.Getauthlibapi(f.ApiAddress)
	if err != nil {
		if errors.Is(err, auth.JsonNotTrue) {
			panic(lang.Lang("auth.JsonNotTrue"))
		} else {
			fmt.Println(lang.Lang("webfail"))
			panic(err)
		}
	}
	f.ApiAddress = api
}
