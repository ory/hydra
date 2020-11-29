/**
 * Copyright (c) Facebook, Inc. and its affiliates.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

import React from 'react'
import useDocusaurusContext from '@docusaurus/useDocusaurusContext'
import Link from '@docusaurus/Link'
import Layout from '@theme/Layout'

import { useVersions, useLatestVersion } from '@theme/hooks/useDocs'

function capitalizeFirstLetter(string) {
  return string.charAt(0).toUpperCase() + string.slice(1)
}

function Version() {
  const { siteConfig } = useDocusaurusContext()
  const versions = useVersions()
  const latestVersion = useLatestVersion()
  const currentVersion = versions.find((version) => version.name === 'current')
  const pastVersions = versions.filter(
    (version) => version !== latestVersion && version.name !== 'current'
  )
  const repoUrl = `https://github.com/${siteConfig.organizationName}/${siteConfig.projectName}`
  const project = `ORY ${capitalizeFirstLetter(siteConfig.projectName)}`

  return (
    <Layout
      title="Versions"
      permalink="/versions"
      description={`Overview of all ${project} documentation versions.`}
    >
      <main className="container margin-vert--lg">
        <h1>{project} documentation versions</h1>

        <div className="margin-bottom--lg">
          <h3 id="next">Current version (Stable)</h3>
          <p>
            Here you can find the documentation for current released version.
          </p>
          <table>
            <tbody>
              <tr>
                <th>{latestVersion.name}</th>
                <td>
                  <Link to={latestVersion.path}>Documentation</Link>
                </td>
                <td>
                  <a href={`${repoUrl}/blob/master/CHANGELOG.md`}>Changelog</a>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        {currentVersion !== latestVersion && (
          <div className="margin-bottom--lg">
            <h3 id="next">Next version (Unreleased)</h3>
            <p>Here you can find the documentation for unreleased version.</p>
            <table>
              <tbody>
                <tr>
                  <th>next</th>
                  <td>
                    <Link to={currentVersion.path}>Documentation</Link>
                  </td>
                  <td>
                    <a href={repoUrl}>Source Code</a>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        )}

        {pastVersions.length > 0 && (
          <div className="margin-bottom--lg">
            <h3 id="archive">Past versions (Not maintained anymore)</h3>
            <p>Here you can find documentation for previous versions.</p>
            <table>
              <tbody>
                {pastVersions.map((version) => (
                  <tr key={version.name}>
                    <th>{version.label}</th>
                    <td>
                      <Link to={version.path}>Documentation</Link>
                    </td>
                    <td>
                      <a href={`${repoUrl}/blob/master/CHANGELOG.md`}>
                        Changelog
                      </a>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </main>
    </Layout>
  )
}

export default Version
