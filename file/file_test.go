package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsFile(t *testing.T) {
	assert.Equal(t, true, IsFile("./testdata/Readable.txt"))
	assert.Equal(t, true, IsFile("./testdata/symbolic.txt"))
	assert.Equal(t, false, IsFile("./testdata"))
	assert.Equal(t, false, IsFile("./testdata/NoReadableDir"))
	assert.Equal(t, true, IsFile("./testdata/.hidden.txt"))
	assert.Equal(t, false, IsFile("abcdef"))
}

func TestExists(t *testing.T) {
	assert.Equal(t, true, Exists("./testdata/Readable.txt"))
	assert.Equal(t, true, Exists("./testdata/symbolic.txt"))
	assert.Equal(t, true, Exists("./testdata"))
	assert.Equal(t, true, Exists("/"))
	assert.Equal(t, false, Exists("abcdef"))
}

func TestIsDir(t *testing.T) {
	assert.Equal(t, false, IsDir("./testdata/Readable.txt"))
	assert.Equal(t, false, IsDir("./testdata/symbolic.txt"))
	assert.Equal(t, true, IsDir("./testdata"))
	assert.Equal(t, true, IsDir("/"))
	assert.Equal(t, false, IsDir("abcdef"))
}

func TestIsSymlink(t *testing.T) {
	assert.Equal(t, false, IsSymlink("./testdata/Readable.txt"))
	assert.Equal(t, true, IsSymlink("./testdata/symbolic.txt"))
	assert.Equal(t, false, IsSymlink("./testdata/"))
	assert.Equal(t, false, IsSymlink("/"))
	assert.Equal(t, false, IsSymlink("abcdef"))
}

func TestIsZero(t *testing.T) {
	assert.Equal(t, true, IsZero("./testdata/Readable.txt"))
	assert.Equal(t, true, IsZero("./testdata/symbolic.txt"))
	assert.Equal(t, false, IsZero("./testdata/"))
	assert.Equal(t, false, IsZero("abcdef"))
}

func TestIsReadable(t *testing.T) {
	assert.Equal(t, true, IsReadable("./testdata/Readable.txt"))
	assert.Equal(t, true, IsReadable("./testdata/symbolic.txt"))
	assert.Equal(t, true, IsReadable("./testdata/"))
	assert.Equal(t, false, IsReadable("abcdef"))
}

func TestIsWritable(t *testing.T) {
	assert.Equal(t, true, IsWritable("./testdata/Writable.txt"))
	assert.Equal(t, true, IsWritable("./testdata/symbolic.txt"))
	assert.Equal(t, true, IsWritable("./testdata"))
	assert.Equal(t, false, IsWritable("abcdef"))
}
func TestIsExecutable(t *testing.T) {
	assert.Equal(t, true, IsExecutable("./testdata/Executable.txt"))
	assert.Equal(t, true, IsExecutable("./testdata/symbolic.txt"))
	assert.Equal(t, true, IsExecutable("./testdata"))
	assert.Equal(t, false, IsExecutable("./testdata/NonExecutable.txt"))
	assert.Equal(t, false, IsExecutable("abcdef"))
}

func TestIsHiddenFile(t *testing.T) {
	assert.Equal(t, false, IsHiddenFile("./testdata/Executable.txt"))
	assert.Equal(t, true, IsHiddenFile("./testdata/.hidden.txt"))
	assert.Equal(t, false, IsHiddenFile("/tmp/mimixbox"))
	assert.Equal(t, false, IsHiddenFile("./testdata"))
	assert.Equal(t, false, IsHiddenFile("abcdef"))
	assert.Equal(t, false, IsHiddenFile(".abcdef"))
	assert.Equal(t, false, IsHiddenFile(".HiddenDir"))
}
