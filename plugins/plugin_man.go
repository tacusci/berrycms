package plugins

import (
	"fmt"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"os"
	"strings"
)

var manager = &Manager{
	dir:     "./plugins",
	scripts: []otto.Script{},
}

// Manager contains plugin collection and add utility and concurrent protection for executing
type Manager struct {
	dir     string
	scripts []otto.Script
}

// NewManager retrieves pointer to only single instance plugin manager
func NewManager() *Manager {
	return manager
}

// Load finds all plugins in provided directory and loads then into manager
func (m *Manager) Load() error {
	m.Unload()

	if err := m.loadFromDir(m.dir); err != nil {
		return err
	}

	return nil
}

func (m *Manager) Unload() {
	m.scripts = []otto.Script{}
}

func (m *Manager) loadFromDir(dir string) error {
	if files, err := ioutil.ReadDir(m.dir); err == nil {
		for i := range files {
			file := files[i]
			fileFullPath := fmt.Sprintf("%s%s%s", dir, string(os.PathSeparator), file.Name())
			// if found directory, call this function to process that directory too
			if file.IsDir() {
				m.loadFromDir(fileFullPath)
			}
			fileNameParts := strings.Split(file.Name(), ".")
			if len(fileNameParts) > 1 {
				if fileNameParts[len(fileNameParts)-1] == "js" {
					m.loadPlugin(fileFullPath)
				}
			}
		}
	} else {
		return err
	}
	return nil
}

func (m *Manager) loadPlugin(pluginFileInfo string) error {
	return nil
}
