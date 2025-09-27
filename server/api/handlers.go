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

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type FileNode struct {
	Name     string     `json:"name"`
	IsDir    bool       `json:"isDir"`
	Children []FileNode `json:"children"`
}

var cwd string
var startCwd string

var Logger *log.Logger

func init() {
	cwdDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}
	cwd = filepath.Join(cwdDir, "user")
	startCwd = filepath.Join(cwdDir, "user")

	// Open log file (create if not exists, append if exists)
	file, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Create a logger that writes to both file and stdout
	Logger = log.New(file, "", log.LstdFlags)
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

		// if strings.HasPrefix(command, "npm run dev") {
		// 	cmd := exec.Command("npm", "run", "dev")

		// 	// Capture stdout and stderr
		// 	stdout, _ := cmd.StdoutPipe()
		// 	stderr, _ := cmd.StderrPipe()

		// 	if err := cmd.Start(); err != nil {
		// 		conn.WriteMessage(websocket.TextMessage, []byte("Error starting dev server: "+err.Error()))
		// 		return
		// 	}

		// 	// Stream stdout
		// 	go func() {
		// 		scanner := bufio.NewScanner(stdout)
		// 		for scanner.Scan() {
		// 			line := scanner.Text()
		// 			conn.WriteMessage(websocket.TextMessage, []byte(line))
		// 		}
		// 	}()

		// 	// Stream stderr
		// 	go func() {
		// 		scanner := bufio.NewScanner(stderr)
		// 		for scanner.Scan() {
		// 			line := scanner.Text()
		// 			conn.WriteMessage(websocket.TextMessage, []byte(line))
		// 		}
		// 	}()

		// 	// Don’t call cmd.Wait() here if you want it to keep running
		// 	continue
		// }

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

func FileTreeWatcher(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Websocket failed to upgrade ", err.Error())
		return
	}

	defer conn.Close()

	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Println("Failed to create watcher ", err)
		return
	}

	defer watcher.Close()

	err = watcher.Add("./user")
	if err != nil {
		log.Println("Failed to watch ./user:", err)
		return
	}

	filepath.Walk("./user", func(path string, info os.FileInfo, err error) error {
		if err == nil && info.IsDir() {
			_ = watcher.Add(path)
		}
		return nil
	})

	log.Println("Started watching ./user")

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			log.Println("File change detected:", event)

			// Regenerate file tree
			tree, err := generateFileTree("./user")
			if err != nil {
				continue
			}

			// Send updated tree to frontend
			treeJson, _ := json.Marshal(tree)
			conn.WriteMessage(websocket.TextMessage, treeJson)

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Watcher error:", err)
		}
	}
}

func LogInfo(msg string) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		Logger.Printf("[INFO] %s:%d %s", file, line, msg)
	} else {
		Logger.Printf("[INFO] %s", msg)
	}
}

func LogError(msg string) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		Logger.Printf("[ERROR] %s:%d %s", file, line, msg)
	} else {
		Logger.Printf("[ERROR] %s", msg)
	}
}

func GetFileCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	relPath := vars["path"]
	//fmt.Fprintf(w, "path is %s", relPath)

	absPath := filepath.Join(startCwd, relPath)

	data, err := os.ReadFile(absPath)

	if err != nil {
		http.Error(w, "Unable to read file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"path":    relPath,
		"content": string(data),
	})

}

func SaveFileCode(w http.ResponseWriter, r *http.Request) {

	type SaveRequest struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}
	var req SaveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	// vars := mux.Vars(r)
	// relPath := vars["path"]
	//fmt.Fprintf(w, "path is %s", relPath)

	absPath := filepath.Join(startCwd, req.Path)

	//absPath := filepath.Join(startCwd, req.Path)
	if !strings.HasPrefix(absPath, startCwd) {
		http.Error(w, "Invalid path", http.StatusForbidden)
		return
	}

	err := os.WriteFile(absPath, []byte(req.Content), 0644)
	if err != nil {
		http.Error(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File saved successfully"))

}
