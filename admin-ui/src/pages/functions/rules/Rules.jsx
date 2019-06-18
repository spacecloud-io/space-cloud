import React from 'react'
import '../../../index.css'
import Sidenav from '../../../components/sidenav/Sidenav'
import Topbar from '../../../components/topbar/Topbar'
import Header from '../../../components/header/Header'
import Documentation from '../../../components/documentation/Documentation'
import rules from '../../../assets/rules.svg'
import EmptyState from '../../../components/rules/EmptyState'

function Rules() {
  return (
    <div>
      <Topbar title="Functions" />
      <div className="flex-box">
        <Sidenav selectedItem="functions" />
        <div className="page-content">
          <div className="header-flex">
            <Header name="Rules" color="#000" fontSize="22px" />
            <Documentation url="https://spaceuptech.com/docs/functions" />
          </div>
          <EmptyState graphics={rules} desc="Guard your data with rules that define who has access to it and how it is structured." buttonText="Add a function" />
        </div>
      </div>
    </div>
  )
}

export default Rules;
