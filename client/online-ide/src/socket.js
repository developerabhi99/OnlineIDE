export const createWebSocket=(url,onMessage,onOpen,onClose,onError)=>{

    const ws= new WebSocket(url);

    ws.onopen=()=>{
        if(onOpen) onOpen()
    }
    
    ws.onclose=()=>{
        if(onClose) onClose()
    }

    ws.onmessage=(evt)=>{
        if(onMessage) onMessage(evt.data)
    }
    
    ws.onerror=()=>{
        if(onError) onError
    }

    return ws

}