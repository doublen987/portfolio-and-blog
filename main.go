package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/doublen987/Projects/MySite/server/functionality"
	"github.com/doublen987/Projects/MySite/server/webportal"
)

type configuration struct {
	ServerAddress      string `json:"webserver"`
	DatabaseType       uint8  `json:"databasetype"`
	DatabaseConnection string `json:"dbconnection"`
	FrontEnd           string `json:"frontend"`
}

func startServer() {
	config, err := functionality.ExtractConfiguration("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	functionality.SetJWTSecret(config.JWTSecretKey)

	address := config.Host + ":" + config.Port
	tlsaddress := config.Host + ":" + config.PortTLS
	log.Println("Starting web server on addres: ", address)
	log.Println("Connecting to database: ", config.DBConnection)

	httpErrChan, httpIsErrChan := webportal.RunAPI(config.Databasetype, address, config.CertPEM, config.KeyPEM, tlsaddress, config.DBConnection, config.DBName, config.FileStorageType)
	for true {
		select {
		case err := <-httpErrChan:
			log.Print("HTTP Error: ", err)
		case err := <-httpIsErrChan:
			log.Print("HTTPS Error: ", err)
		}
	}
}

var PIDFile = "/tmp/daemonize"

func savePID(pid int, fileName string) {

	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("Unable to create pid file : %v\n", err)
		os.Exit(1)
	}

	defer file.Close()

	_, err = file.WriteString(strconv.Itoa(pid))

	if err != nil {
		log.Printf("Unable to create pid file : %v\n", err)
		os.Exit(1)
	}

	file.Sync() // flush to disk

}

func SayHelloWorld(w http.ResponseWriter, r *http.Request) {
	html := "Hello World"

	w.Write([]byte(html))
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage : %s [start|stop] \n ", os.Args[0]) // return the program name back to %s
		os.Exit(0)                                            // graceful exit
	}

	if strings.ToLower(os.Args[1]) == "main" {
		instanceName := os.Args[2]
		// Make arrangement to remove PID file upon receiving the SIGTERM from kill command
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)

		go func() {
			signalType := <-ch
			signal.Stop(ch)
			fmt.Println("Exit command received. Exiting...")

			// this is a good place to flush everything to disk
			// before terminating.
			fmt.Println("Received signal type : ", signalType)

			// remove PID file
			os.Remove(PIDFile + instanceName + ".pid")

			os.Exit(0)

		}()

		startServer()
	}

	if strings.ToLower(os.Args[1]) == "start" {

		instanceName := os.Args[2]
		fmt.Println(PIDFile + instanceName + ".pid")
		// check if daemon already running.
		if _, err := os.Stat(PIDFile + instanceName + ".pid"); err == nil {
			fmt.Println("Already running or /tmp/daemonize" + instanceName + ".pid file exist.")
			os.Exit(1)
		}

		cmd := exec.Command(os.Args[0], "main", instanceName)
		cmd.Start()
		fmt.Println("Daemon process ID is : ", cmd.Process.Pid)
		savePID(cmd.Process.Pid, PIDFile+instanceName+".pid")
		os.Exit(0)
	}

	// upon receiving the stop command
	// read the Process ID stored in PIDfile
	// kill the process using the Process ID
	// and exit. If Process ID does not exist, prompt error and quit

	if strings.ToLower(os.Args[1]) == "stop" {
		instanceName := os.Args[2]

		fmt.Println(PIDFile + instanceName + ".pid")
		if _, err := os.Stat(PIDFile + instanceName + ".pid"); err == nil {

			data, err := ioutil.ReadFile(PIDFile + instanceName + ".pid")
			if err != nil {
				fmt.Println("Not running")
				os.Exit(1)
			}
			ProcessID, err := strconv.Atoi(string(data))

			if err != nil {
				fmt.Println("Unable to read and parse process id found in ", PIDFile+instanceName+".pid")
				os.Exit(1)
			}

			process, err := os.FindProcess(ProcessID)

			if err != nil {
				fmt.Printf("Unable to find process ID [%v] with error %v \n", ProcessID, err)
				os.Exit(1)
			}
			// remove PID file
			os.Remove(PIDFile + instanceName + ".pid")

			fmt.Printf("Killing process ID [%v] now.\n", ProcessID)
			// kill process and exit immediately
			err = process.Kill()

			if err != nil {
				fmt.Printf("Unable to kill process ID [%v] with error %v \n", ProcessID, err)
				os.Exit(1)
			} else {
				fmt.Printf("Killed process ID [%v]\n", ProcessID)
				os.Exit(0)
			}

		} else {

			fmt.Println("Not running.")
			os.Exit(1)
		}
	} else {
		fmt.Printf("Unknown command : %v\n", os.Args[1])
		fmt.Printf("Usage : %s [start|stop]\n", os.Args[0]) // return the program name back to %s
		os.Exit(1)
	}

}
