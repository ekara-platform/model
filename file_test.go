package model

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrlExists(t *testing.T) {
	p := "http://www.google.com/my_path"
	b, _ := FileExist(p)
	assert.False(t, b)

}

func TestFileExists(t *testing.T) {
	p := "./testFile.txt"
	c := []byte("filecontent\n")

	f, e := os.Create(p)
	assert.Nil(t, e)

	_, e = f.Write(c)
	assert.Nil(t, e)
	f.Close()
	b, _ := FileExist(p)
	assert.True(t, b)

	e = os.Remove(p)
	assert.Nil(t, e)
	b, _ = FileExist(p)
	assert.False(t, b)
}

func TestDirExists(t *testing.T) {
	p := "./testDir"

	e := os.MkdirAll(p, 0755)
	assert.Nil(t, e)
	assert.True(t, DirExist(p))

	e = os.Remove(p)
	assert.Nil(t, e)
	assert.False(t, DirExist(p))
}
