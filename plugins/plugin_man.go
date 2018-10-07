package plugins

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/tacusci/logging"

	"github.com/robertkrimen/otto"
)

func NewManager() *Manager {
	man := &Manager{
		pluginsDirPath: "./plugins",
		Plugins:        &[]Plugin{},
	}
	man.load()
	return man
}

type Manager struct {
	pluginsDirPath string
	Plugins        *[]Plugin
}

func (m *Manager) load() {
	pluginFiles, err := ioutil.ReadDir(m.pluginsDirPath)
	if err != nil {
		logging.Error(err.Error())
		return
	}

	for _, file := range pluginFiles {
		plugin := m.loadPlugin(file)
		if plugin != nil {
			*m.Plugins = append(*m.Plugins, *plugin)
		}
	}
}

func (m *Manager) loadPlugin(file os.FileInfo) *Plugin {
	if m.validatePlugin(file) {
		plugin := &Plugin{filePath: fmt.Sprintf("%s%s%s", m.pluginsDirPath, string(filepath.Separator), file.Name())}
		if plugin.loadRuntime() {
			return plugin
		}
	}
	return nil
}

func (m *Manager) NewExtPlugin() *Plugin {
	return &Plugin{}
}

func (m *Manager) ExecAll() {
	for _, plugin := range *m.Plugins {
		plugin.Run()
	}
}

func (m *Manager) validatePlugin(fi os.FileInfo) bool {
	return strings.Contains(fi.Name(), ".js")
}

type Plugin struct {
	runtime  *otto.Otto
	filePath string
}

func (p *Plugin) loadRuntime() bool {
	p.runtime = otto.New()
	return p.runtime != nil
}

func (p *Plugin) setApiFuncs() {
	if p.runtime != nil {
		p.runtime.Set("InfoLog", PluginInfoLog)
	}
}

func (p *Plugin) Run() bool {

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

	p.setApiFuncs()

	if _, err := p.runtime.Run(buff.String()); err != nil {
		logging.Error(err.Error())
		return false
	}

	return true
}

func (p *Plugin) Call(funcName string) {
	p.runtime.Call(funcName, nil, nil)
}
