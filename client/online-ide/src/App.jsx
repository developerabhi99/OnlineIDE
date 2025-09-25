import { useEffect, useState } from 'react'
import Terminal from './component/terminal'
import './App.css'
import FileTree from './component/Tree'

function App() {
  const [fileTree, setFileTree] = useState({})


  const getFileTree= async ()=>{
    const response= await fetch('http://localhost:8080/files')
    const result=await response.json();

    console.log(result)

    setFileTree(result)

  }

  useEffect(()=>{
    getFileTree()
  },[])

  return (
    <div className="playground-container">
    {/* Editor + FileTree */}
    <div className="editor-container">
      <div className="file-container">
        <h4>Explorer</h4>
        <FileTree tree={fileTree} />
      </div>

      <div className="file-editor">
        <pre>
{`function hello() {
console.log("Hello World");
}`}
        </pre>
      </div>
    </div>

    {/* Terminal */}
    <div className="terminal-container">
      <Terminal />
    </div>
  </div>
  )
}

export default App
