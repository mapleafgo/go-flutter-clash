package go_flutter_clash

import (
	"errors"
	"os"
	"path/filepath"

	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/hub/executor"
	"github.com/Dreamacro/clash/hub/route"
	flutter "github.com/go-flutter-desktop/go-flutter"
	"github.com/go-flutter-desktop/go-flutter/plugin"
	"github.com/mapleafgo/go-flutter-clash/go/config"
)

const channelName = "go_flutter_clash"

// GoFlutterClashPlugin implements flutter.Plugin and handles method.
type GoFlutterClashPlugin struct{}

var _ flutter.Plugin = &GoFlutterClashPlugin{} // compile-time type check

// InitPlugin initializes the plugin.
func (p *GoFlutterClashPlugin) InitPlugin(messenger plugin.BinaryMessenger) error {
	channel := plugin.NewMethodChannel(messenger, channelName, plugin.StandardMethodCodec{})
	channel.HandleFunc("init", p.initClash)
	channel.HandleFunc("start", p.start)
	return nil
}

func (p *GoFlutterClashPlugin) initClash(arguments interface{}) (reply interface{}, err error) {
	if params, ok := arguments.([]interface{}); ok {
		var homeDir string
		if params[0] != nil {
			homeDir = params[0].(string)
			if !filepath.IsAbs(homeDir) {
				currentDir, _ := os.Getwd()
				homeDir = filepath.Join(currentDir, homeDir)
			}
			C.SetHomeDir(homeDir)
			return nil, nil
		}
	}
	return nil, errors.New("props error")
}

func (p *GoFlutterClashPlugin) start(arguments interface{}) (reply interface{}, err error) {
	if params, ok := arguments.([]interface{}); ok {
		var profile, fcc string
		if params[0] != nil {
			profile = params[0].(string)
		}
		if params[1] != nil {
			fcc = params[1].(string)
		}
		cfg, err := config.Parse(profile, fcc)
		if err != nil {
			return nil, err
		}
		go route.Start("127.0.0.1:9090", cfg.General.Secret)
		executor.ApplyConfig(cfg, true)
		return nil, nil
	}
	return nil, errors.New("props error")
}
