package main

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/viper"
	"gitlab.com/rosenpin/good-morning/config"
	"gitlab.com/rosenpin/good-morning/provider"
	"gitlab.com/rosenpin/good-morning/querier"
	"gitlab.com/rosenpin/good-morning/result"
	"gitlab.com/rosenpin/good-morning/url"
)

func main() {
	loadConfig()
	conf := config.Config{}
	viper.Unmarshal(&conf)

	querier := querier.JSONQuerier{}
	urlCreator := url.GoogleImagesCreator{}
	parser := result.GoogleImagesResultParser{}

	provider := provider.NewImageProvider(querier, urlCreator, parser, conf)

	link, err := provider.Provide()
	if err != nil {
		panic(err)
	}

	openBrowser(link)
}

func openBrowser(url string) bool {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}

func loadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/good-morning/")
	viper.AddConfigPath("$HOME/.good-morning")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}
