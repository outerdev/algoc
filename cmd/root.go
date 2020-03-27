package cmd

import (
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"

	. "github.com/outerdev/algoc/config"
	. "github.com/outerdev/algoc/errors"
	. "github.com/outerdev/algoc/lock"
	. "github.com/outerdev/algoc/prompt"

	"github.com/algorand/go-algorand/cmd/kmd/codes"
	"github.com/algorand/go-algorand/daemon/kmd"
	"github.com/algorand/go-algorand/daemon/kmd/server"
	"github.com/algorand/go-algorand/logging"
)

const (
	kmdLogFileName = "kmd.log"
	kmdLogFilePerm = 0640
)

func runKmd(dataDir string, timeoutSecs uint64) {
	// Use logging package instead of stdin/stdout
	log := logging.NewLogger()
	log.SetLevel(logging.Info)

	// Parse timeout duration. 0 timeout -> nil timeout
	var timeout *time.Duration
	if timeoutSecs != 0 {
		t := time.Duration(timeoutSecs) * time.Second
		timeout = &t
	}

	// We have a dataDir now, so use log files
	kmdLogFilePath := filepath.Join(dataDir, kmdLogFileName)
	kmdLogFileMode := os.O_CREATE | os.O_WRONLY | os.O_APPEND
	logFile, err := os.OpenFile(kmdLogFilePath, kmdLogFileMode, kmdLogFilePerm)
	if err != nil {
		log.Errorf("failed to open log file: %s", err)
		os.Exit(codes.ExitCodeKMDLogError)
	}
	log.SetOutput(logFile)

	// Prevent swapping with mlockall if supported by the platform
	TryMlockall(log)

	// Create a "kill" channel to allow the server to shut down gracefully
	kill := make(chan os.Signal)

	// Timeouts can also send on the kill channel; because signal.Notify
	// will not block, this shouldn't cause an issue. From docs: "Package
	// signal will not block sending to c"
	signal.Notify(kill, os.Interrupt, unix.SIGTERM, unix.SIGINT)
	signal.Ignore(unix.SIGHUP)

	// Build a kmd StartConfig
	startConfig := kmd.StartConfig{
		DataDir: dataDir,
		Kill:    kill,
		Log:     log,
		Timeout: timeout,
	}

	// Start the kmd server
	died, sock, err := kmd.Start(startConfig)
	if err == server.ErrAlreadyRunning {
		log.Errorf("couldn't start kmd: %s", err)
		os.Exit(codes.ExitCodeKMDAlreadyRunning)
	}
	if err != nil {
		log.Errorf("couldn't start kmd: %s", err)
		os.Exit(codes.ExitCodeKMDError)
	}

	log.Infof("started kmd on sock: %s", sock)

	// Wait until the kmd server exits
	<-died
	log.Infof("kmd server died. exiting...")
}

type Config struct {
	Token string `yaml:"token" action:"prompt"`
	Host  string `yaml:"host" action:"prompt,url"`
}

var config Config

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	SetConfigFileName("algoc.yaml")
	if err := LoadConfig(&config); IsConfigNotPresent(err) {
		if err := PromptForValues(&config); err != nil {
			Fatal(err)
		}
		if err := WriteConfig(config); err != nil {
			Fatal(err)
		}
	} else if err != nil {
		Fatal(err)
	}
}

func init() {
	initConfig()

	configDir, err := LocateConfigDir()
	if err != nil {
		panic(err)
	}

	dataDir := configDir + "/data"
	if !DirExists(dataDir) {
		err := os.Mkdir(dataDir, 0755)
		if err != nil {
			panic(err)
		}
	}

	timeoutSecs := uint64(0)
	runKmd(dataDir, timeoutSecs)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		Fatal(err)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "algoc",
	Short: "Small command line utility to interact with the Algorand network",
	Long: `Command line utility to interact with the Algorand network.

Manage and issue assets with a few commands.`,
}
