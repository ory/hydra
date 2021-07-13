// according to https://github.com/facebook/docusaurus/issues/1258#issuecomment-594393744

// use in *.mdx like:

// import Mermaid from '@theme/Mermaid'
//
// <Mermaid chart={`
// flowchart TD
//     cr([Create Request]) --> backoffice[Backoffice Server REST]
// `}/>

import React, { useEffect, useState } from 'react'
import mermaid from 'mermaid'
import styles from './mermaid.module.css'
import cn from 'classnames'

mermaid.initialize({
  startOnLoad: true,
  logLevel: 'fatal',
  securityLevel: 'strict',
  arrowMarkerAbsolute: false,
  theme: 'neutral',
  flowchart: {
    useMaxWidth: true,
    htmlLabels: true,
    rankSpacing: 65,
    nodeSpacing: 30,
    curve: 'basis'
  },
  sequence: {
    useMaxWidth: true
  },
  gantt: {
    useMaxWidth: true
  }
})

const Mermaid = ({ chart }) => {
  const [zoomed, setZoomed] = useState(false)
  const [svg, setSvg] = useState(undefined)
  const [id] = useState(`mermaid-${Math.random().toString(36).substr(2, -1)}`)
  const toggle = () => setZoomed(!zoomed)

  useEffect(() => {
    mermaid.render(id, chart, (svg) => {
      setSvg(svg)
    })
  }, [])

  return (
    <>
      <div
        onClick={toggle}
        className={cn(styles.graph, styles.pointer)}
        dangerouslySetInnerHTML={{ __html: svg }}
      />
      <div
        onClick={toggle}
        className={cn(styles.overlay, styles.pointer, styles.graph, {
          [styles.visible]: zoomed
        })}
      >
        <div
          onClick={(e) => e.stopPropagation()}
          className={cn(styles.backdrop, styles.graph)}
          dangerouslySetInnerHTML={{ __html: svg }}
        />
      </div>
    </>
  )
}

export default Mermaid
