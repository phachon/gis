package app

import (
	"github.com/spf13/viper"
	"os"
	"flag"
	"github.com/phachon/go-logger"
	"strings"
	"github.com/fatih/color"
)

var (
	flagConf = flag.String("conf", "config.toml", "please input conf path")
)

var (
	Version = "v0.2"

	author = "phachon"

	address = "https://github.com/phachon/gis"

	RootPath = ""

	AppPath = ""

	Conf = viper.New()

	Log = go_logger.NewLogger()
)

// 启动初始化
func init()  {
	initFlag()
	initPath()
	initConfig()
	initPoster()
	initLog()
}

// print
func initPoster() {
	fg := color.New(color.FgBlue)
	fg.Println(` 
             _        
   ____ _   (_)  _____
  / __  /  / /  / ___/
 / /_/ /  / /  (__  )
 \__, /  /_/  /____/
/____/
		`)
	fg.Printf(" Version: %s \r\n", Version)
	fg.Printf(" Address: %s \r\n", address)
	fg.Printf(" Author : %s \r\n", author)
	fg.Println("-------------------------------------")
}

// init flag
func initFlag() {
	flag.Parse()
}

// init dir and path
func initPath() {
	AppPath, _ = os.Getwd()
	RootPath = strings.Replace(AppPath, "app", "", 1)
}

// init config
func initConfig()  {

	if *flagConf == "" {
		Log.Error("config file not found!")
		os.Exit(1)
	}

	Conf.SetConfigType("toml")
	Conf.SetConfigFile(*flagConf)
	err := Conf.ReadInConfig()
	if err != nil {
		Log.Error("Fatal error config file: "+err.Error())
		os.Exit(1)
	}

	file := Conf.ConfigFileUsed()
	if file != "" {
		Log.Info("Use config file: " + file)
	}
}

// init log
func initLog() {

	Log.Detach("console")

	// console adapter config
	consoleLevelStr := Conf.GetString("log.console.level")
	consoleConfig := &go_logger.ConsoleConfig{
		Color: Conf.GetBool("log.console.color"), // show text color
		JsonFormat: Conf.GetBool("log.console.jsonFormat"), // text json format
	}
	Log.Attach("console", Log.LoggerLevel(consoleLevelStr), consoleConfig)

	// file adapter config
	fileLevelStr := Conf.GetString("log.file.level")
	levelFilenameConf := Conf.GetStringMapString("log.file.levelFilename")
	levelFilename := map[int]string{}
	if len(levelFilenameConf) > 0 {
		for levelStr, levelFile := range levelFilenameConf {
			level := Log.LoggerLevel(levelStr)
			levelFilename[level] = levelFile
		}
	}
	fileConfig := &go_logger.FileConfig{
		Filename:  Conf.GetString("log.file.filename"),
		LevelFileName : levelFilename,
		MaxSize: Conf.GetInt64("log.file.maxSize"),
		MaxLine: Conf.GetInt64("log.file.maxLine"),
		DateSlice: Conf.GetString("log.file.dateSlice"),
		JsonFormat: Conf.GetBool("log.file.jsonFormat"),
	}
	Log.Attach("file", Log.LoggerLevel(fileLevelStr), fileConfig)
}