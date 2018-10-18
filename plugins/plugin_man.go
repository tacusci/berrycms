package plugins

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Manager struct {
	dir string
}

// Load finds all plugins in provided directory and loads then into manager
func (m *Manager) Load(dir string) error {
	m.dir = dir

	if err := m.loadFromDir(m.dir); err != nil {
		return err
	}

	return nil
}

func (m *Manager) loadFromDir(dir string) error {
	if files, err := ioutil.ReadDir(m.dir); err == nil {
		for i := range files {
			file := files[i]
			// if found directory, call this function to process that directory too
			if file.IsDir() {
				m.loadFromDir(fmt.Sprintf("%s%s%s", dir, string(os.PathSeparator), file.Name()))
			}
			fileNameParts := strings.Split(file.Name(), ".")
			if len(fileNameParts) > 1 {
				if fileNameParts[len(fileNameParts)-1] == ".js" {

				}
			}
		}
	} else {
		return err
	}
}
