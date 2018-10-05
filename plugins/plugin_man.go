package plugins

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/tacusci/logging"

	"github.com/robertkrimen/otto"
)

func NewManager() *Manager {
	man := &Manager{}
	man.load("./plugins")
	return man
}

type Manager struct {
	pluginsDirPath string
	plugins        *[]Plugin
}

func (m *Manager) load(dirPath string) {
	pluginFiles, err := ioutil.ReadDir(dirPath)
	if err != nil {
		logging.Error(err.Error())
		return
	}

	m.pluginsDirPath = dirPath

	for _, file := range pluginFiles {
		plugin := Plugin{filePath: fmt.Sprintf("%s%s%s", dirPath, string(filepath.Separator), file.Name())}
		if plugin.loadRuntime() {
			*m.plugins = append(*m.plugins, plugin)
		}
	}
}

type Plugin struct {
	runtime  *otto.Otto
	filePath string
}

func (p *Plugin) loadRuntime() bool {
	f, err := os.Open(p.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		logging.Error(err.Error())
		return false
	}
	defer f.Close()

	buff := bytes.NewBuffer(nil)

	if _, err := buff.ReadFrom(f); err != nil {
		logging.Error(err.Error())
		return false
	}

	runtime := otto.New()
	if _, err := runtime.Run(buff.String()); err != nil {
		logging.Error(err.Error())
		return false
	}

	p.runtime = runtime

	return true
}
