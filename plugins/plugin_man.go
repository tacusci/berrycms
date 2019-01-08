// Copyright (c) 2019, tacusci ltd
//
// Licensed under the GNU GENERAL PUBLIC LICENSE Version 3 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.gnu.org/licenses/gpl-3.0.html
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package plugins

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/gofrs/uuid"

	"github.com/robertkrimen/otto"
)

var manager = &Manager{
	dir:     "./plugins",
	plugins: []Plugin{},
}

type Plugin struct {
	uuid     string
	filePath string
	src      string
	VM       *otto.Otto
	Document *goquery.Document
}

func (p *Plugin) UUID() string { return p.uuid }

func (p *Plugin) ParseFile() error {
	if p.filePath != "" && p.filePath != "-" {
		data, err := ioutil.ReadFile(p.filePath)
		if err != nil {
			return err
		}
		p.src = string(data)
	}
	return nil
}

func (p *Plugin) Call(funcName string, this interface{}, argumentList ...interface{}) (otto.Value, error) {
	if err := p.ParseFile(); err != nil {
		return otto.Value{}, err
	}

	if _, err := p.VM.Run(p.src); err != nil {
		return otto.Value{}, err
	}

	return p.VM.Call(funcName, this, argumentList)
}

// Manager contains plugin collection and add utility and concurrent protection for executing
type Manager struct {
	sync.Mutex
	dir     string
	plugins []Plugin
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
	m.Lock()
	defer m.Unlock()
	m.plugins = []Plugin{}
}

func (m *Manager) Plugins() *[]Plugin {
	return &m.plugins
}

func (m *Manager) loadFromDir(dir string) error {
	files, err := ioutil.ReadDir(m.dir)
	if err != nil {
		return err
	}
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
	return nil
}

func (m *Manager) loadPlugin(fileFullPath string) error {
	m.Lock()
	defer m.Unlock()

	if uuidV4, err := uuid.NewV4(); err == nil {
		plugin := Plugin{
			uuid:     uuidV4.String(),
			VM:       otto.New(),
			filePath: fileFullPath,
			Document: &goquery.Document{},
		}

		if err := plugin.ParseFile(); err != nil {
			return err
		}

		plugin.VM.Set("UUID", plugin.uuid)
		plugin.VM.Set("InfoLog", PluginInfoLog)
		plugin.VM.Set("DebugLog", PluginDebugLog)
		plugin.VM.Set("ErrorLog", PluginErrorLog)
		plugin.VM.Run(plugin.src)

		m.plugins = append(m.plugins, plugin)
	} else {
		return err
	}

	return nil
}
