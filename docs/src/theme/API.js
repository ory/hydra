import React from 'react'
import useBaseUrl from '@docusaurus/useBaseUrl'
import { useActiveVersion } from '@theme/hooks/useDocs'
import Redoc from '@theme/Redoc'
import styles from './API.module.css'

function join(...args) {
  return args
    .map((part, i) => {
      if (i === 0) {
        return part.trim().replace(/[\/]*$/g, '')
      } else {
        return part.trim().replace(/(^[\/]*|[\/]*$)/g, '')
      }
    })
    .filter((x) => x.length)
    .join('/')
}

function API({ spec }) {
  const { path } = useActiveVersion()
  return <Redoc spec={spec} />
}

export default API
