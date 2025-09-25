
const FileTreeNode = ({ node }) => {
    return (
      <li>
        {node.name}
        {node.children && node.children.length > 0 && (
          <ul>
            {node.children.map((child, i) => (
              <FileTreeNode key={i} node={child} />
            ))}
          </ul>
        )}
      </li>
    )
  }
  
  const FileTree = ({ tree }) => {
    if (!tree || !tree.name) return null // handle empty case
  
    return (
      <ul>
        <FileTreeNode node={tree} />
      </ul>
    )
  }
  
  export default FileTree
  