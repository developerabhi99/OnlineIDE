import { Terminal as XTerminal } from '@xterm/xterm'
import { useEffect, useRef } from 'react'
import '@xterm/xterm/css/xterm.css'
import { createWebSocket } from '../socket'

const Terminal = () => {
  const terminalRef = useRef() //
  const isRendered = useRef(false) // for checking its not rendered twice
  const termRef = useRef(null)
  const wsRef = useRef(null)  // for keeping websocket ref
  const inputBufferRef = useRef('') // ref to store input entered by user before enter

  const printPrompt = (cwd) => {
    termRef.current.write(`\r\n${cwd}> `)
  }

  useEffect(() => {
    if (isRendered.current) return
    isRendered.current = true

    const term = new XTerminal({
      rows: 20,
      cols: 150,
      cursorBlink: true,
      allowProposedApi: true, // needed for some clipboard functions
    })
    term.open(terminalRef.current)
    termRef.current = term

    

    const ws = createWebSocket(
        'ws://localhost:8080/ws',
        (data) => {
          if (data.startsWith('__CWD__:')) {
            const cwd = data.replace('__CWD__:', '').trim()
            termRef.current.write(`\r\n${cwd}> `) // show prompt
          } else {
           
            termRef.current.write(`\r\n${data}`)
            ws.send('__GET_CWD__')
          }
        },
        () => {
          termRef.current.writeln('Connected to backend!')
          ws.send('__GET_CWD__')
        },
        () => termRef.current.writeln('\r\nConnection closed.'),
        (err) => termRef.current.writeln('\r\nWebSocket error.')
      )
    wsRef.current = ws


    term.onData((data) => {

        //when user press enter 
      if (data === '\r') {
        if (ws.readyState === WebSocket.OPEN) {
          ws.send(inputBufferRef.current + '\n')
        }
        inputBufferRef.current = ''
      } else if (data === '\u007F') { //when user press back
        if (inputBufferRef.current.length > 0) {
          inputBufferRef.current = inputBufferRef.current.slice(0, -1)
          term.write('\b \b')
        }
      } else {
        inputBufferRef.current += data
        term.write(data)
      }
    })
  }, [])

  return <div ref={terminalRef} id="terminal" style={{ width: '100%', height: '400px' }}></div>
}

export default Terminal
