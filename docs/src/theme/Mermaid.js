// according to https://github.com/facebook/docusaurus/issues/1258#issuecomment-594393744

// use in *.mdx like:

// import Mermaid from '@theme/Mermaid'
//
// <Mermaid chart={`
// flowchart TD
//     cr([Create Request]) --> backoffice[Backoffice Server REST]
// `}/>

import React, { useEffect } from "react"
import mermaid from "mermaid"

mermaid.initialize({
  startOnLoad: true
})

const Mermaid = ({ chart }) => {
  useEffect(() => {
    mermaid.contentLoaded()
  }, [])
  return <div className="mermaid">{chart}</div>
}

export default Mermaid
