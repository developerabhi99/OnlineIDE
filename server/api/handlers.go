package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func MainHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Welcome to Online Code Editor by developerabhi99")

}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Websocket upgrade failed ", err)
	}

	defer conn.Close()

	fmt.Println("WebSocket Upgraded successfully")

	for {
		msgType, msg, err := conn.ReadMessage()

		if err != nil {
			log.Println("Unable to read message ", err)
			break
		}

		command := string(msg)

		//fmt.Printf("Recieved msg:%s\n", msg)

		log.Printf("Executing command: %s\n", command)

		// Choose shell based on OS
		var cmd *exec.Cmd
		currentUserFile, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get cwd: %v", err)
		}

		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/C", command)
			cmd.Dir = currentUserFile + "/user"
		} else {
			cmd = exec.Command("sh", "-c", command)
			cmd.Dir = currentUserFile + "/user"
		}

		if command == "__GET_CWD__" {
			cwd := currentUserFile + "/user"
			conn.WriteMessage(websocket.TextMessage, []byte("__CWD__:"+cwd))
			continue
		}

		// Capture combined stdout+stderr
		output, err := cmd.CombinedOutput()
		if err != nil {
			output = append(output, []byte("\nCommand error: "+err.Error())...)
		}

		if writeErr := conn.WriteMessage(msgType, output); writeErr != nil {
			log.Println("Error while writing message:", writeErr)
			break
		}
	}

}
