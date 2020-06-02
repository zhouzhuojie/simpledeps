package simpledeps

import (
	"fmt"
)

// Package defines a struct that holds a node in the dependency graph
// It leverages DependsOn to define their relationships
type Package struct {
	Name      string
	DependsOn []*Package

	// dependedByInstalled is a variable to store what other
	// packages depends on this package at runtime
	dependedByInstalled map[string]struct{}
}

// NewPackage creates a new instance of Package
func NewPackage(name string) *Package {
	return &Package{
		Name:                name,
		DependsOn:           nil,
		dependedByInstalled: make(map[string]struct{}),
	}
}

// PackageManager helps to manage the dependency graph of packages
type PackageManager struct {
	allPackages map[string]*Package
	installed   map[string]*Package
}

func NewPackageManager() *PackageManager {
	return &PackageManager{
		allPackages: make(map[string]*Package),
		installed:   make(map[string]*Package),
	}
}

// Install is a recursive function that installs that package
// It works as a dfs travesal into the dependency graph
func (pm *PackageManager) Install(name string) ([]string, error) {
	p, ok := pm.allPackages[name]
	if !ok {
		return nil, fmt.Errorf("package %s not defined", name)
	}

	// if already installed, break the recursion early
	if _, ok := pm.installed[name]; ok {
		return nil, nil
	}

	pm.installed[name] = pm.allPackages[name]

	ret := []string{}
	for i := range p.DependsOn {
		dep := p.DependsOn[i]

		// update dependedByInstalled
		dep.dependedByInstalled[name] = struct{}{}

		// do Install() recursion
		path, err := pm.Install(dep.Name)
		if err != nil {
			return nil, err
		}
		ret = append(path, ret...)
	}
	ret = append(ret, name)
	return ret, nil
}

// Remove is a recursive function that removes the package
// It works as a dfs traversal into the dependency graph
func (pm *PackageManager) Remove(name string) ([]string, error) {
	// if not installed, directly return
	p, ok := pm.installed[name]
	if !ok {
		return nil, nil
	}

	// if there are other packages that depends on this package, return err
	if len(p.dependedByInstalled) != 0 {
		return nil, fmt.Errorf("there are other packages that depends on %s", p.Name)
	}

	delete(pm.installed, name)

	ret := []string{}
	for i := range p.DependsOn {
		dep := p.DependsOn[i]
		delete(dep.dependedByInstalled, name)
		path, _ := pm.Remove(dep.Name)
		ret = append(ret, path...)
	}
	ret = append([]string{name}, ret...)
	return ret, nil
}

// List lists the installed packages
func (pm *PackageManager) List() []string {
	ret := []string{}
	for _, p := range pm.installed {
		if p != nil {
			ret = append(ret, p.Name)
		}
	}
	return ret
}

// Define defines the package and its dependencies
// For example, a -> b, c
//	    pm := NewPackageManager()
//		pm.Define(a, b, c)
func (pm *PackageManager) Define(p *Package, deps ...*Package) {
	pm.allPackages[p.Name] = p
	for i := range deps {
		pm.allPackages[deps[i].Name] = deps[i]
		p.DependsOn = append(p.DependsOn, deps[i])
	}
}
