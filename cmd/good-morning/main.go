package main

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/viper"
	"gitlab.com/rosenpin/good-morning/config"
	"gitlab.com/rosenpin/good-morning/querier"
)

func main() {
	loadConfig()
	conf := config.Config{}
	viper.Unmarshal(&conf)

	querier := querier.NewQuerier(conf)
	result, err := querier.Query()
	if err != nil {
		panic(err)
	}

	r, ok := result.(map[string]interface{})
	if !ok {
		fmt.Println("invalid result", ok)
		fmt.Println(result)
		return
	}

	items, ok := r["items"].([]interface{})
	if !ok {
		fmt.Println("invalid result", ok)
		fmt.Println(result)
		return
	}

	if len(items) != 1 {
		fmt.Println("invalid items", items)
		fmt.Println(result)
		return
	}

	item := items[0].(map[string]interface{})

	openBrowser(item["link"].(string))
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
