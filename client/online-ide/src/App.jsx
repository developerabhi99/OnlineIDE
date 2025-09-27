import { useEffect, useState } from "react";
import Terminal from "./component/terminal";
import "./App.css";
import FileTree from "./component/Tree";
import AceEditor from "react-ace";
import "ace-builds/src-noconflict/mode-javascript";
import "ace-builds/src-noconflict/theme-monokai";
import "ace-builds/src-noconflict/ext-language_tools";

function App() {
  const [fileTree, setFileTree] = useState({});
  const [selectedFile,setSelectedFile]=useState('')
  const [selectedCode,setSelectedCode]=useState('')

  const getFileTree = async () => {
    const response = await fetch("http://localhost:8080/files");
    const result = await response.json();

    //console.log(result)

    setFileTree(result);
  };

  useEffect(() => {
    getFileTree();
  }, []);

  useEffect(() => {
    const fileWs = new WebSocket("ws://localhost:8080/fileWatcher");

    fileWs.onmessage = (tree) => {
      const updatedTree = JSON.parse(tree.data);
      setFileTree(updatedTree);
    };

    return () => fileWs.close();
  }, []);

  useEffect(() => {
    const handleKeyDown = (e) => {
      if ((e.ctrlKey || e.metaKey) && e.key === "s") {
        e.preventDefault(); // stop browser save
        saveFile();
      }
    };

    const saveFile = async () => {
      if (!selectedFile) return;
      try {
        await fetch(`http://localhost:8080/saveFile`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            path: selectedFile,
            content: selectedCode,
          }),
        });
        console.log("File saved:", selectedFile);
      } catch (err) {
        console.error("Save failed", err);
      }
    };

    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [selectedFile, selectedCode]);


  return (
    <div className="playground-container">
      {/* Editor + FileTree */}
      <div className="editor-container">
        <div className="file-container">
          <h4>Explorer</h4>
          <FileTree tree={fileTree} onSelect={async (path)=> {

            const res = await fetch(`http://localhost:8080/fileCode/${path}`)
            //const fileData=await res.json()
          //  console.log("fileData ",res)
            const fileData = await res.json()
              setSelectedFile(path)
              console.log("constructed path ",path)
              console.log("fileData ",fileData.content)
              setSelectedCode(fileData.content)
          }

     
           
            } 
            
            
            />
        </div>

        <div className="file-editor">
        {selectedFile && <p>{selectedFile.replaceAll('/',' >')}</p>}
          <pre>
            <AceEditor
              height="60vh"
              width="100%"
              mode="javascript"
              theme="monokai"
              fontSize="16px"
              highlightActiveLine={true}
              setOptions={{
                enableLiveAutocompletion: true,
                showLineNumbers: true,
                tabSize: 2,
              }}
              value={selectedCode}
              onChange={(e)=> setSelectedCode(e)}
            />
          </pre>
        </div>
      </div>

      {/* Terminal */}
      <div className="terminal-container">
        <Terminal />
      </div>
    </div>
  );
}

export default App;
