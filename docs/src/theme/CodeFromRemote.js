import React, {useEffect, useState} from 'react'
import fetch from 'node-fetch'
import CodeBlock from '@theme/CodeBlock'
import styles from './CodeFromRemote.module.css';

const CodeFromRemote = ({src, link, lang}) => {
  const [content, setContent] = useState('')

  useEffect(() => {
    fetch(src).then(body => body.text()).then(setContent).catch(console.err)
  }, [])

  return (
    <div className={styles.container}>
      <CodeBlock metastring={link && `title="${link}"`} className={lang && `language-${lang}`} children={content}/>
    </div>
  )
}

export default CodeFromRemote
