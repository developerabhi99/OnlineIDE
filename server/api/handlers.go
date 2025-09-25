package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gorilla/websocket"
)

type FileNode struct {
	Name     string     `json:"name"`
	IsDir    bool       `json:"isDir"`
	Children []FileNode `json:"children"`
}

var cwd string
var startCwd string

func init() {
	cwdDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}
	cwd = filepath.Join(cwdDir, "user")
	startCwd = filepath.Join(cwdDir, "user")
}

var cmd *exec.Cmd
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

		if command == "__GET_CWD__" {
			conn.WriteMessage(websocket.TextMessage, []byte("__CWD__:"+cwd))
			continue
		}

		if strings.HasPrefix(command, "cd ") {
			dir := strings.TrimSpace(command[3:])
			newPath := filepath.Join(cwd, dir)

			log.Println("New path ", newPath)
			// Resolve the absolute path
			absPath, err := filepath.Abs(newPath)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Error: invalid path"))
				continue
			}
			log.Println("absPath path ", absPath)
			log.Println("startCwd path ", startCwd)
			// Ensure user cannot go above the startCwd
			if !strings.HasPrefix(absPath, startCwd) {
				conn.WriteMessage(websocket.TextMessage, []byte("Error: cannot go below user folder"))
				continue
			}

			// Check if it exists and is a directory
			fi, err := os.Stat(absPath)
			if err != nil || !fi.IsDir() {
				conn.WriteMessage(websocket.TextMessage, []byte("Error: directory does not exist"))
				continue
			}

			// Update current working directory
			cwd = absPath
			//cmd.Dir = cwd
			conn.WriteMessage(websocket.TextMessage, []byte("__CWD__:"+cwd))
			continue
		}

		// if strings.HasPrefix(command, "cd ") {
		// 	//cwd = currentUserFile + "/user"
		// 	// Extract the path after "cd "
		// 	dir := strings.TrimSpace(command[3:])

		// 	// Build the new path relative to current cwd
		// 	newPath := filepath.Join(cwd, dir)

		// 	// Check if it exists and is a directory
		// 	fi, err := os.Stat(newPath)
		// 	if err != nil || !fi.IsDir() {
		// 		conn.WriteMessage(websocket.TextMessage, []byte("Error: directory does not exist"))
		// 		continue
		// 	}

		// 	// Update current working directory

		// 	cwd = newPath
		// 	cmd.Dir = cwd

		// 	// Send updated cwd to frontend
		// 	conn.WriteMessage(websocket.TextMessage, []byte("__CWD__:"+cwd))
		// 	continue
		// }

		// Selecting shell based on OS

		// currentUserFile, err := os.Getwd()
		// if err != nil {
		// 	log.Fatalf("Failed to get cwd: %v", err)
		// }

		// if runtime.GOOS == "windows" {
		// 	cmd = exec.Command("cmd", "/C", command)
		// 	cmd.Dir = currentUserFile + "/user"
		// } else {
		// 	cmd = exec.Command("sh", "-c", command)
		// 	cmd.Dir = currentUserFile + "/user"
		// }

		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/C", command)
		} else {
			cmd = exec.Command("sh", "-c", command)
		}
		//log.Println("before cwd", cmd.Dir)
		// if cwd == startCwd {
		// 	continue
		// } else {

		// }
		cmd.Dir = cwd
		//log.Println("after cwd", cmd.Dir)

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

func GetFileTree(w http.ResponseWriter, r *http.Request) {

	fileTree, err := generateFileTree("./user")
	if err != nil {
		http.Error(w, "Failed to generate file tree: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fileTree)

}

func generateFileTree(path string) (FileNode, error) {
	initialDir, err := os.Stat(path)

	if err != nil {
		return FileNode{}, err
	}

	node := FileNode{
		Name:  initialDir.Name(),
		IsDir: initialDir.IsDir(),
	}
	// if initail file is directory
	if initialDir.IsDir() {
		//read the directory

		entries, err := os.ReadDir(path)

		if err != nil {
			return FileNode{}, err
		}

		// looping through entries if given path is directory
		for _, entry := range entries {
			childPath := filepath.Join(path, entry.Name())
			childNode, err := generateFileTree(childPath)

			if err != nil {
				return FileNode{}, err
			}

			node.Children = append(node.Children, childNode)
		}
	}
	return node, nil
}
