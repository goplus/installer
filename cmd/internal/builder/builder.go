package builder

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Build struct {
	Source bool // if true, OS and Arch must be empty
	OS     string
	Arch   string
	Root   string

	Verbose bool
	Repo    string
	Branch  string
	Out     string
}

func (b *Build) Do() error {
	work, err := ioutil.TempDir("", "bindist")
	if err != nil {
		return err
	}
	defer func() {
		_ = os.RemoveAll(work)
	}()
	b.Root = work

	_, err = run(work, "git", "clone", b.Repo)
	if err != nil {
		return err
	}
	src := filepath.Join(work, "gop")
	if b.Branch != "" {
		_, err = run(src, "git", "checkout", b.Branch)
		if err != nil {
			return err
		}
	}

	if b.OS == "windows" {
		_, err = run(src, "cmd", "/c", "make.bat")
	} else {
		_, err = run(src, "bash", "make.bash")
	}
	if err != nil {
		return err
	}

	var (
		version string
		_       = version
	)
	// TODO: windows ?
	fullVersion, err := run("", filepath.Join(src, "bin/gop"), "version")
	if err != nil {
		return err
	}
	fullVersion = bytes.TrimSpace(fullVersion)
	v := bytes.SplitN(fullVersion, []byte(" "), 2)
	version = string(v[0])

	//base := fmt.Sprintf("gop.%s.%s-%s", version, b.OS, b.Arch)
	switch b.OS {
	case "darwin":
		err := b.buildDarwinPkg()
		if err != nil {
			return err
		}
	case "windows":
	default:
		panic("not support")
	}
	return nil
}

func run(dir, name string, args ...string) ([]byte, error) {
	buf := new(bytes.Buffer)
	absName, err := lookPath(name)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(absName, args...)
	var output io.Writer = buf
	//if b.verbose {
	if true {
		log.Printf("running %q %q", absName, args)
		output = io.MultiWriter(buf, os.Stdout)
	}
	cmd.Stdout = output
	cmd.Stderr = output
	cmd.Dir = dir
	//cmd.Env = b.env()
	if err := cmd.Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s", err.Error())
		return nil, fmt.Errorf("%s %s: %v", name, strings.Join(args, " "), err)
	}
	return buf.Bytes(), nil
}

var cleanEnv = []string{
	"GOARCH",
	"GOBIN",
	"GOHOSTARCH",
	"GOHOSTOS",
	"GOOS",
	"GOROOT",
	"GOROOT_FINAL",
}

func (b *Build) env() []string {
	env := os.Environ()
	for i := 0; i < len(env); i++ {
		for _, c := range cleanEnv {
			if strings.HasPrefix(env[i], c+"=") {
				env = append(env[:i], env[i+1:]...)
			}
		}
	}
	final := "/usr/local/go"
	if b.OS == "windows" {
		final = `c:\go`
	}
	env = append(env,
		"GOARCH="+b.Arch,
		"GOHOSTARCH="+b.Arch,
		"GOHOSTOS="+b.OS,
		"GOOS="+b.OS,
		"GOROOT="+b.Root,
		"GOROOT_FINAL="+final,
	)
	return env
}

func lookPath(prog string) (absPath string, err error) {
	absPath, err = exec.LookPath(prog)
	if err == nil {
		return
	}
	//t, ok := windowsDeps[prog]
	//if !ok {
	//	return
	//}
	//for _, dir := range t.commonDirs {
	//	for _, ext := range []string{"exe", "bat"} {
	//		absPath = filepath.Join(dir, prog+"."+ext)
	//		if _, err1 := os.Stat(absPath); err1 == nil {
	//			err = nil
	//			os.Setenv("PATH", os.Getenv("PATH")+";"+dir)
	//			return
	//		}
	//	}
	//}
	return
}
