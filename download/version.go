package download

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/xmdhs/gomclauncher/lang"
)

func Getversionlist(cxt context.Context, atype string, print func(string)) (*version, error) {
	var rep *http.Response
	var err error
	var b []byte
	r := newrandurls(atype)
	_, f := r.auto()
	for i := 0; i < 4; i++ {
		if err := func() error {
			if i == 3 {
				return fmt.Errorf("Getversionlist: %w", err)
			}
			rep, _, err = Aget(cxt, source(`https://piston-meta.mojang.com/mc/game/version_manifest.json`, f))
			if rep != nil {
				defer rep.Body.Close()
			}
			if err != nil {
				print(fmt.Sprint(lang.Lang("getversionlistfail"), fmt.Errorf("Getversionlist: %w", err), source(`https://piston-meta.mojang.com/mc/game/version_manifest.json`, f)))
				f = r.fail(f)
				return nil
			}
			b, err = ioutil.ReadAll(rep.Body)
			if err != nil {
				print(fmt.Sprint(lang.Lang("getversionlistfail"), fmt.Errorf("Getversionlist: %w", err), source(`https://piston-meta.mojang.com/mc/game/version_manifest.json`, f)))
				f = r.fail(f)
				return nil
			}
			return errors.New("")
		}(); err != nil {
			if err.Error() == "" {
				break
			} else {
				return nil, fmt.Errorf("Getversionlist: %w", err)
			}
		}
	}
	v := version{}
	err = json.Unmarshal(b, &v)
	v.atype = atype
	if err != nil {
		return nil, fmt.Errorf("Getversionlist: %w", err)
	}
	return &v, nil
}

type version struct {
	Latest   versionLatest    `json:"latest"`
	Versions []versionVersion `json:"versions"`
	atype    string
}

type versionLatest struct {
	Release  string `json:"release"`
	Snapshot string `json:"snapshot"`
}

type versionVersion struct {
	ID          string `json:"id"`
	ReleaseTime string `json:"releaseTime"`
	Time        string `json:"time"`
	Type        string `json:"type"`
	URL         string `json:"url"`
}

func (v version) Downjson(cxt context.Context, version, apath string, print func(string)) error {
	r := newrandurls(v.atype)
	_, f := r.auto()
	for _, vv := range v.Versions {
		if vv.ID == version {
			s := strings.Split(vv.URL, "/")
			path := apath + `/versions/` + vv.ID + `/` + vv.ID + `.json`
			if ver(path, s[len(s)-2]) {
				return nil
			}
			for i := 0; i < 4; i++ {
				if i == 3 {
					return FileDownLoadFail
				}
				err := get(cxt, source(vv.URL, f), path)
				if err != nil {
					print(fmt.Sprint(lang.Lang("weberr"), source(vv.URL, f), fmt.Errorf("Downjson: %w", err)))
					f = r.fail(f)
					continue
				}
				if !ver(path, s[len(s)-2]) {
					print(fmt.Sprint(lang.Lang("filecheckerr"), source(vv.URL, f)))
					f = r.fail(f)
					continue
				}
				break
			}
			return nil
		}
	}
	return NoSuch
}

var (
	NoSuch         = errors.New("no such")
	ErrFileChecker = errors.New("file checker")
)
