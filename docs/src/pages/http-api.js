import React from 'react'
import ApiDoc from '@theme/ApiDoc'
import useBaseUrl from '@docusaurus/useBaseUrl'
import { useActiveVersion } from '@theme/hooks/useDocs'
import { Redirect } from '@docusaurus/router'
import config from '../../config'

function CustomPage() {
  const { path } = useActiveVersion()
  if (!config.enableRedoc) {
    return <Redirect to={useBaseUrl(path)} />
  }
  return (
    <ApiDoc
      layoutProps={{
        title: 'HTTP API Docs',
        description: `Read the HTTP API reference documentation`
      }}
      spec={{
        type: 'url',
        content: useBaseUrl(`${path}.static/api.json`)
      }}
    />
  )
}

export default CustomPage
