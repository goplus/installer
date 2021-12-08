package builder

import (
	"os"
	"path/filepath"

	"github.com/goplus/build/cmd/internal/assets"
)

type darwinEntry struct {
	distributionPath string
	scriptsPath      string
	darwinPkg        string
}

func (b *Build) buildDarwinPkg() error {
	entry, err := initDarwinEntry(b)
	if err != nil {
		return err
	}

	pkgDest := filepath.Join(b.Root, "target/pkgdest")
	err = os.MkdirAll(pkgDest, 0755)
	if err != nil {
		return err
	}
	_, err = run("", "pkgbuild",
		"--identifier", "org.goplus.gop",
		"--version", "1.0",
		"--scripts", entry.scriptsPath,
		"--root", entry.darwinPkg,
		filepath.Join(pkgDest, "org.goplus.gop.pkg"))
	if err != nil {
		return err
	}

	// TODO: package name
	targ := filepath.Join(b.Root, "target/gop.pkg")
	_, err = run("", "productbuild",
		"--distribution", entry.distributionPath,
		"--package-path", pkgDest,
		targ)
	_, err = run("", "cp", targ, b.Out)
	if err != nil {
		return err
	}
	return nil
}

func initDarwinEntry(b *Build) (*darwinEntry, error) {
	e := new(darwinEntry)

	// darwinPkg
	darwinPkg := filepath.Join(b.Root, "target/darwinpkg")
	err := os.MkdirAll(darwinPkg, 0755)
	if err != nil {
		return nil, err
	}
	e.darwinPkg = darwinPkg

	// etc/paths.d/gop
	gopData, err := assets.F.ReadFile("res/darwin/etc/paths.d/gop")
	if err != nil {
		return nil, err
	}
	err = writeByteToFile(gopData, filepath.Join(darwinPkg, "etc/paths.d/gop"))
	if err != nil {
		return nil, err
	}

	// gop src
	err = copySrcToDarwinPkg(filepath.Join(b.Root, "gop"), darwinPkg)
	if err != nil {
		return nil, err
	}

	// darwin distribution script/postinstall
	darwin := filepath.Join(b.Root, "target/darwin")
	err = os.MkdirAll(darwin, 0755)
	if err != nil {
		return nil, err
	}

	// darwin distribution
	distributionData, err := assets.F.ReadFile("res/darwin/Distribution")
	if err != nil {
		return nil, err
	}
	distributionPath := filepath.Join(darwin, "Distribution")
	err = writeByteToFile(distributionData, distributionPath)
	if err != nil {
		return nil, err
	}
	e.distributionPath = distributionPath

	// darwin script/postinstall
	postinstallData, err := assets.F.ReadFile("res/darwin/scripts/postinstall")
	if err != nil {
		return nil, err
	}
	scriptPath := filepath.Join(darwin, "scripts/postinstall")
	err = writeByteToFile(postinstallData, scriptPath)
	if err != nil {
		return nil, err
	}
	e.scriptsPath = scriptPath

	return e, nil
}

func writeByteToFile(data []byte, path string) error {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func copySrcToDarwinPkg(form, darwinPkg string) error {
	to := filepath.Join(darwinPkg, "usr/local")
	err := os.MkdirAll(to, 0755)
	if err != nil {
		return err
	}
	_, err = run("", "cp", "-a", form, to)
	if err != nil {
		return err
	}

	// rm .git .github
	var cleanFiles = []string{
		".git",
		".gitattributes",
		".github",
		".gitignore",
	}
	for _, p := range cleanFiles {
		_ = os.RemoveAll(filepath.Join(to, "gop", p))
	}

	return nil
}
