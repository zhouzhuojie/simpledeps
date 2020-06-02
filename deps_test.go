package simpledeps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func generateSamplePackageManager() *PackageManager {
	// a -> b,c
	// b -> c,d
	// c -> f
	// e -> c
	// g -> h

	// install a -> f,c,d,b,a // order
	// remove c -> error
	// # remove a -> a,b,c,d,f //
	// install e -> e
	// remove e -> e  // print all packages that are removed
	// list -> a,b,c,d,f // no order

	a := NewPackage("a")
	b := NewPackage("b")
	c := NewPackage("c")
	d := NewPackage("d")
	e := NewPackage("e")
	f := NewPackage("f")
	g := NewPackage("g")
	h := NewPackage("h")

	pm := NewPackageManager()
	pm.Define(a, b, c)
	pm.Define(b, c, d)
	pm.Define(c, f)
	pm.Define(e, c)
	pm.Define(g, h)
	return pm
}

func TestSimpleDeps(t *testing.T) {
	t.Run("happy code path for e", func(t *testing.T) {
		// install a -> f,c,d,b,a // order
		// remove c -> error
		// install e -> e
		// remove e -> e  // print all packages that are removed
		// list -> a,b,c,d,f // no order
		pm := generateSamplePackageManager()

		// should be able to install a
		installPath, err := pm.Install("a")
		assert.NoError(t, err)
		assert.Equal(t, []string{"d", "f", "c", "b", "a"}, installPath)

		// should not be able to remove c
		_, err = pm.Remove("c")
		assert.Error(t, err)

		// should be able to install e
		installPath, err = pm.Install("e")
		assert.NoError(t, err)
		assert.Equal(t, []string{"e"}, installPath)

		// should be able to remove e
		removePath, err := pm.Remove("e")
		assert.NoError(t, err)
		assert.Equal(t, []string{"e"}, removePath)

		// should be able to list the current install packages
		list := pm.List()
		assert.ElementsMatch(t, []string{"a", "b", "c", "d", "f"}, list)
	})

	t.Run("happy code path for a", func(t *testing.T) {
		// install a -> f,c,d,b,a
		// remove c -> error
		// remove a -> a,b,c,d,f
		// install e -> f, c, e
		// remove e -> e, c, f
		// list -> empty
		pm := generateSamplePackageManager()

		// should be able to install a
		installPath, err := pm.Install("a")
		assert.NoError(t, err)
		assert.Equal(t, []string{"d", "f", "c", "b", "a"}, installPath)

		// should not be able to remove c
		_, err = pm.Remove("c")
		assert.Error(t, err)

		// should be able to remove a
		removePath, err := pm.Remove("a")
		assert.NoError(t, err)
		assert.Equal(t, []string{"a", "b", "d", "c", "f"}, removePath)

		// should be able to install e
		installPath, err = pm.Install("e")
		assert.NoError(t, err)
		assert.Equal(t, []string{"f", "c", "e"}, installPath)

		// should be able to remove e
		removePath, err = pm.Remove("e")
		assert.NoError(t, err)
		assert.Equal(t, []string{"e", "c", "f"}, removePath)
	})
}
