import React, {useEffect, useState} from 'react'
import fetch from 'node-fetch'
import CodeBlock from '@theme/CodeBlock'
import styles from './CodeFromRemote.module.css';

const CodeFromRemote = ({url, link, ...props}) => {
  const [content, setContent] = useState('')
  let {src, lang, title} = props

  if (url) {
    src = url
      .replace('github.com', 'raw.githubusercontent.com')
      .replace('/blob/', '/')
  }

  if (link) {
    title= link
  }

  useEffect(() => {
    fetch(src)
      .then(body => body.text())
      .then((content) => {
        // https://github.com/ory/kratos-selfservice-ui-react-native/blob/master/App.tsx#L37-L65
        const params = src.match(/^https:\/\/raw.githubusercontent.com\/.+#L([0-9]+)-L([0-9]+)$/) || []

        if (params.length === 3) {
          const lines = content.split('\n')
          lines.splice(0, params[1] - 1)
          lines.splice(params[2] - params[1] - 1)
          return lines.join('\n')
        }

        return content
      })
      .then(setContent)
      .catch(console.err)
  }, [])

  return (
    <div className={styles.container}>
      <CodeBlock metastring={title && `title="${title}"`} className={lang && `language-${lang}`} children={content}/>
    </div>
  )
}

export default CodeFromRemote
