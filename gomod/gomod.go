package gomod

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"go.x2ox.com/sorbifolia/random"
	"golang.org/x/mod/modfile"
)

func New(repoURL, branch string, cdn bool) *Package {
	return &Package{
		RepoURL: repoURL,
		Branch:  branch,
		CDN:     cdn,
	}
}

func Parse(filename string) ([]*Package, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var arr []*Package
	if err = json.Unmarshal(data, &arr); err != nil {
		return nil, err
	}
	return arr, nil
}

type Package struct {
	RepoURL    string `json:"repo_url"`
	Branch     string `json:"branch"`
	CDN        bool   `json:"cdn"`
	mainModule string
	path       string
}

func (p *Package) FindModule() ([]string, error) {
	if err := p.gitClone(); err != nil {
		return nil, err
	}
	return p.findModule(""), nil
}

func (p Package) Output(packageName string) error {
	return (packageData{
		Main:    p.mainModule,
		PkgName: packageName,
		Repo:    p.RepoURL,
		Branch:  p.Branch,
		ReadMe:  p.readMe(),
	}).WriteToFile()
}

func (p Package) Clean() error {
	return os.RemoveAll(p.path)
}

func (p *Package) gitClone() error {
	if p.path == "" {
		p.path = ".go.mod.data-" + random.NewMathRand().RandString(10)
	}
	if _, err := git.PlainClone(p.path, false, &git.CloneOptions{
		URL:           p.RepoURL,
		ReferenceName: plumbing.NewBranchReferenceName(p.Branch),
		SingleBranch:  true,
	}); err != nil {
		return err
	}

	return nil
}

func (p *Package) findModule(dir string) []string {
	var (
		files, err = os.ReadDir(filepath.Join(p.path, dir))
		arr        []string
	)
	if err != nil {
		return nil
	}

	for _, v := range files {
		if v.IsDir() {
			arr = append(arr, p.findModule(filepath.Join(dir, v.Name()))...)
			continue
		}
		if v.Name() == "go.mod" {
			if pkgName, _ := parseModFile(filepath.Join(p.path, dir, "go.mod")); pkgName != "" {
				arr = append(arr, pkgName)

				if p.mainModule == "" || p.mainModule > pkgName {
					p.mainModule = pkgName
				}
			}
		}
	}

	return arr
}

func (p Package) readMe() string {
	repo := strings.TrimPrefix(p.RepoURL, "https://github.com/")

	if p.CDN {
		return fmt.Sprintf("https://cdn.jsdelivr.net/gh/%s@%s/README.md", repo, p.Branch)
	}
	return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/README.md", repo, p.Branch)
}

func parseModFile(filename string) (string, error) {
	bts, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	var file *modfile.File
	if file, err = modfile.Parse(filename, bts, func(_, ver string) (string, error) {
		return ver, nil
	}); err != nil {
		return "", err
	}
	return file.Module.Mod.Path, nil
}
