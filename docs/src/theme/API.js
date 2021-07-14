import React from 'react'
import Redoc from '@theme/Redoc'
import './API.module.css'

function API({ spec }) {
  return <Redoc spec={spec} />
}

export default API
