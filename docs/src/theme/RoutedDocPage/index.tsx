import React from 'react';
import type {Props} from '@theme/DocPage';
import DocPage from '@theme/DocPage'

export default function RoutedDocPage(props: Props): JSX.Element {
  return <div id={"route-identifier"} data-route={props.location.pathname}><DocPage {...props}/></div>
}
