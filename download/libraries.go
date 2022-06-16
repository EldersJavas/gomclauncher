package download

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xmdhs/gomclauncher/auth"
	"github.com/xmdhs/gomclauncher/internal"
	"github.com/xmdhs/gomclauncher/lang"
	"github.com/xmdhs/gomclauncher/launcher"
)

type Libraries struct {
	librarie   launcher.LauncherjsonX115
	assetIndex assets
	typee      string
	cxt        context.Context
	print      func(string)
	path       string
	*randurls
}

func Newlibraries(cxt context.Context, b []byte, typee string, print func(string), apath string) (Libraries, error) {
	mod := launcher.Modsjson{}
	var url, id string
	l := launcher.LauncherjsonX115{}
	err := json.Unmarshal(b, &mod)
	r := newrandurls(typee)
	if err != nil {
		return Libraries{}, fmt.Errorf("Newlibraries: %w", err)
	}
	if mod.InheritsFrom != "" {
		b, err := ioutil.ReadFile(apath + `/versions/` + mod.InheritsFrom + "/" + mod.InheritsFrom + ".json")
		if err != nil {
			return Libraries{}, fmt.Errorf("Newlibraries: %w", err)
		}
		err = json.Unmarshal(b, &l)
		if err != nil {
			return Libraries{}, fmt.Errorf("Newlibraries: %w", err)
		}
		modlibraries2(mod.Libraries, &l)
		l.ID = mod.ID
	} else {
		err = json.Unmarshal(b, &l)
		if err != nil {
			return Libraries{}, fmt.Errorf("Newlibraries: %w", err)
		}
	}
	for i := range l.Libraries {
		launcher.FullLibraryX115(&l.Libraries[i], "")
	}
	url = l.AssetIndex.URL
	id = l.AssetIndex.ID
	path := apath + "/assets/indexes/" + id + ".json"
	if !ver(path, l.AssetIndex.Sha1) {
		err := assetsjson(cxt, r, url, path, typee, l.AssetIndex.Sha1, print)
		if err != nil {
			return Libraries{}, fmt.Errorf("Newlibraries: %w", err)
		}
	}
	bb, err := ioutil.ReadFile(path)
	if err != nil {
		return Libraries{}, fmt.Errorf("Newlibraries: %w", err)
	}
	a := assets{}
	err = json.Unmarshal(bb, &a)
	if err != nil {
		return Libraries{}, fmt.Errorf("Newlibraries: %w", err)
	}
	return Libraries{
		print:      print,
		librarie:   l,
		assetIndex: a,
		typee:      typee,
		cxt:        cxt,
		randurls:   r,
		path:       apath,
	}, nil
}

type assets struct {
	Objects map[string]asset `json:"objects"`
}

type asset struct {
	Hash string `json:"hash"`
}

func get(cxt context.Context, u, path string) error {
	reps, timer, err := Aget(cxt, u)
	if reps != nil {
		defer reps.Body.Close()
	}
	if err != nil {
		return fmt.Errorf("get: %w", err)
	}
	if reps.StatusCode != 200 {
		return fmt.Errorf("get: %w", &ErrHTTPCode{code: reps.StatusCode})
	}
	_, err = os.Stat(path)
	if err != nil {
		dir, _ := filepath.Split(path)
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			return fmt.Errorf("get: %w", err)
		}
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("get: %w", err)
	}
	defer f.Close()
	bw := bufio.NewWriter(f)
	for {
		timer.Reset(5 * time.Second)
		i, err := io.CopyN(bw, reps.Body, 100000)
		if err != nil && err != io.EOF {
			return fmt.Errorf("get: %w", err)
		}
		if i == 0 {
			break
		}
	}
	err = bw.Flush()
	if err != nil {
		return fmt.Errorf("get: %w", err)
	}
	return nil
}

type ErrHTTPCode struct {
	code int
}

func (e *ErrHTTPCode) Error() string {
	return fmt.Sprintf("http code: %d", e.code)
}

func modlibraries2(l []launcher.Librarie, Launcherjson *launcher.LauncherjsonX115) {
	for _, v := range l {
		l := librarie2LibraryX115(&v)
		Launcherjson.Libraries = append(Launcherjson.Libraries, *l)
	}
}

var mirror = map[string]map[string]string{
	"bmclapi": {
		`launchermeta.mojang.com`:          `bmclapi2.bangbang93.com`,
		`piston-meta.mojang.com`:           `bmclapi2.bangbang93.com`,
		`launcher.mojang.com`:              `bmclapi2.bangbang93.com`,
		`resources.download.minecraft.net`: `bmclapi2.bangbang93.com/assets`,
		`libraries.minecraft.net`:          `bmclapi2.bangbang93.com/maven`,
		`files.minecraftforge.net/maven`:   `bmclapi2.bangbang93.com/maven`,
		`maven.minecraftforge.net`:         `bmclapi2.bangbang93.com/maven`,
	},
	"mcbbs": {
		`launchermeta.mojang.com`:          `download.mcbbs.net`,
		`piston-meta.mojang.com`:           `download.mcbbs.net`,
		`launcher.mojang.com`:              `download.mcbbs.net`,
		`resources.download.minecraft.net`: `download.mcbbs.net/assets`,
		`libraries.minecraft.net`:          `download.mcbbs.net/maven`,
		`files.minecraftforge.net/maven`:   `download.mcbbs.net/maven`,
		`maven.minecraftforge.net`:         `download.mcbbs.net/maven`,
	},
}

func source(url, types string) string {
	m, ok := mirror[types]
	if ok {
		for k, v := range m {
			if strings.Contains(url, k) {
				return strings.Replace(url, k, v, 1)
			}
		}
	}
	return url
}

func Aget(cxt context.Context, aurl string) (*http.Response, *time.Timer, error) {
	return internal.HttpGet(cxt, aurl, auth.Transport, nil)
}

func assetsjson(cxt context.Context, r *randurls, url, path, typee, sha1 string, print func(string)) error {
	var err error
	_, f := r.auto()
	for i := 0; i < 4; i++ {
		if i == 3 {
			return err
		}
		err = get(cxt, source(url, f), path)
		if err != nil {
			f = r.fail(f)
			print(lang.Lang("weberr") + " " + fmt.Errorf("assetsjson: %w", err).Error() + " " + url)
			continue
		}
		if !ver(path, sha1) {
			f = r.fail(f)
			err = ErrFileChecker
			print(lang.Lang("filecheckerr") + " " + url)
			continue
		}
		break
	}
	return nil
}
