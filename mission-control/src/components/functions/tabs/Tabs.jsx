import React from 'react';
import { Tabs } from "antd"
import history from "../../../history"
import "./tabs.css"
const { TabPane } = Tabs;

export default ({ activeKey, projectId }) => {
  const onTabClick = (key) => {
    history.push(`/mission-control/projects/${projectId}/functions/${key}`)
  }
  return (
    <div className="functions-tabs">
      <Tabs activeKey={activeKey} onTabClick={onTabClick}>
        <TabPane tab="Rules" key="rules">
        </TabPane>
        <TabPane tab="Testing" key="testing">
        </TabPane>
      </Tabs>
    </div>
  )
}