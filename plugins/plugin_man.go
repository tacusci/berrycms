package plugins

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gobuffalo/uuid"
	"github.com/tacusci/logging"

	"github.com/robertkrimen/otto"
)

var pluginsList *[]Plugin = &[]Plugin{}

func NewManager() *Manager {
	man := &Manager{
		pluginsDirPath: "./plugins",
		plugins:        pluginsList,
	}
	man.CompileAll()
	return man
}

type Manager struct {
	sync.Mutex
	pluginsDirPath string
	plugins        *[]Plugin
}

func (m *Manager) GetPlugins() *[]Plugin {
	m.Lock()
	defer m.Unlock()
	return m.plugins
}

func (m *Manager) SetPlugins(plugins []Plugin) {
	m.Lock()
	defer m.Unlock()
	m.plugins = &plugins
}

func (m *Manager) LoadPlugins() error {
	err := m.load()
	if err != nil {
		logging.Error(fmt.Sprintf("Unable to load plugins, -> %s", err.Error()))
		return err
	}
	return nil
}

func (m *Manager) UnloadPlugins() {
	m.SetPlugins([]Plugin{})
}

func (m *Manager) load() error {
	pluginFiles, err := ioutil.ReadDir(m.pluginsDirPath)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir("./plugins", os.ModeDir)
			if os.IsPermission(err) {
				logging.Error(fmt.Sprintf("Unable to create plugins dir, permission denied -> %s", err.Error()))
				return err
			}
			pluginFiles, err = ioutil.ReadDir(m.pluginsDirPath)
			if err != nil {
				return err
			}
		}
		return err
	}

	for _, file := range pluginFiles {
		plugin := m.loadPlugin(file)
		if plugin != nil {
			plugins := m.GetPlugins()
			m.SetPlugins(append(*plugins, *plugin))
		}
	}
	return nil
}

func (m *Manager) loadPlugin(file os.FileInfo) *Plugin {
	if m.validatePlugin(file) {
		if uuidV4, err := uuid.NewV4(); err == nil {
			plugin := &Plugin{
				UUID:     uuidV4.String(),
				filePath: fmt.Sprintf("%s%s%s", m.pluginsDirPath, string(filepath.Separator), file.Name()),
			}
			if plugin.loadRuntime() {
				return plugin
			}
		} else {
			return nil
		}
	}
	return nil
}

func (m *Manager) NewExtPlugin() *Plugin {
	return &Plugin{}
}

func (m *Manager) CompileAll() {
	for _, plugin := range *m.plugins {
		plugin.Compile()
		plugin.setGlobalConsts()
	}
}

func (m *Manager) validatePlugin(fi os.FileInfo) bool {
	return strings.Contains(fi.Name(), ".js")
}

type Plugin struct {
	runtime   *otto.Otto
	UUID      string
	filePath  string
	compiled  bool
	WaitGroup sync.WaitGroup
}

func (p *Plugin) loadRuntime() bool {
	p.runtime = otto.New()
	return p.runtime != nil
}

func (p *Plugin) setApiFuncs() {
	if p.runtime != nil {
		p.runtime.Set("InfoLog", PluginInfoLog)
		p.runtime.Set("DebugLog", PluginDebugLog)
		p.runtime.Set("ErrorLog", PluginErrorLog)
	}
}

func (p *Plugin) setGlobalConsts() {
	p.runtime.Set("UUID", p.UUID)
}

func (p *Plugin) Compile() bool {

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

	p.compiled = true

	return p.compiled
}

func (p *Plugin) Call(funcName string, this interface{}, argumentList ...interface{}) (otto.Value, error) {
	if argumentList != nil {
		return p.runtime.Call(funcName, this, argumentList)
	}
	return otto.Value{}, errors.New("Argument list is nil")
}
