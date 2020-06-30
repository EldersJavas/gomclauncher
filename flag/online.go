package flag

import (
	"fmt"
	"os"
	"time"

	"github.com/xmdhs/gomclauncher/auth"
)

func (f *Flag) Aonline() {
	if f.Email == "" {
		fmt.Println("请设置邮箱")
		fmt.Println("比如 -email xxx@xxx.xx")
		os.Exit(0)
	}
	err := gmlconfig[f.Email].setonline(f.Email, f.Passworld)
	if err != nil {
		if err.Error() == "have" {
			a := auth.Auth{
				AccessToken: gmlconfig[f.Email].AccessToken,
				ClientToken: gmlconfig[f.Email].ClientToken,
			}
			atime := time.Now().Unix()
			if atime-gmlconfig[f.Email].Time > 1200 {
				if err := auth.Validate(a); err != nil {
					err = auth.Refresh(&a)
					if err != nil {
						if err.Error() == "not ok" {
							fmt.Println("请尝试重新登录")
							os.Exit(0)
						} else {
							fmt.Println("可能是网络问题，可再次尝试")
							os.Exit(0)
						}
					}
					aconfig := gmlconfig[f.Email]
					aconfig.Name = a.Username
					aconfig.UUID = a.ID
					aconfig.AccessToken = a.AccessToken
					aconfig.Time = time.Now().Unix()
					aconfig.ClientToken = a.ClientToken
					gmlconfig[f.Email] = aconfig
					saveconfig()
				}
			}
		} else if err.Error() == "not ok" {
			fmt.Println("账户名或密码错误")
			os.Exit(0)
		} else {
			fmt.Println(err)
			os.Exit(0)
		}
	}
	f.Userproperties = gmlconfig[f.Email].Userproperties
	f.AccessToken = gmlconfig[f.Email].AccessToken
	f.Name = gmlconfig[f.Email].Name
	f.UUID = gmlconfig[f.Email].UUID
}
