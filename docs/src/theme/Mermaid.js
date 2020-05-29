// according to https://github.com/facebook/docusaurus/issues/1258#issuecomment-594393744

// use in *.mdx like:

// import Mermaid from '@theme/Mermaid'
//
// <Mermaid chart={`
// flowchart TD
//     cr([Create Request]) --> backoffice[Backoffice Server REST]
// `}/>

import React, {useEffect} from "react"
import mermaid from "mermaid"
import styles from './mermaid.module.css';
import cn from 'classnames'

mermaid.initialize({
  startOnLoad: true,
  logLevel: 'fatal',
  securityLevel: 'strict',
  arrowMarkerAbsolute: false,
  theme: "neutral",
  flowchart: {
    useMaxWidth: true,
    htmlLabels: true,
    rankSpacing: 65,
    nodeSpacing: 30,
    curve: "basis"
  },
  sequence:{
    useMaxWidth: true,
  },
  gantt:{
    useMaxWidth: true,
  }
})

const Mermaid = ({chart}) => {
  useEffect(() => {
    mermaid.contentLoaded()
  }, [])
  return <div className={cn(styles.graph, "mermaid")}>{chart}</div>
}

export default Mermaid
