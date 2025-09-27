
const FileTreeNode = ({ onSelect,path,node }) => {
    return (
      <li onClick={(e)=>{
        e.stopPropagation()
        if (node.isDir) return ;
        onSelect(path);
      }}>
        {node.name!="node_modules" && node.name}
        {node.name!="node_modules" && (node.children && node.children.length) > 0 && (
          <ul>
            {node.children.map((child, i) => (
              <FileTreeNode onSelect={onSelect} path={path +'/'+child.name} key={i} node={child} />
            ))}
          </ul>
        )}
      </li>
    )
  }
  
  const FileTree = ({ tree,onSelect }) => {
    if (!tree || !tree.name) return null // handle empty case
  
    return (
      <ul>
        <FileTreeNode node={tree} path="" onSelect={onSelect} />
      </ul>
    )
  }
  
  export default FileTree
  