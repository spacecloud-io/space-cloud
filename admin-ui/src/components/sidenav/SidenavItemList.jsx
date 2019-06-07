import React from 'react'
import SidenavItem from './SidenavItem'

function SidenavItemList(props) {
  
  return (
    <div className="list">
      {props.items.map((item) =>  <SidenavItem name={item.name} icon={item.icon} active={item.key == props.selectedItem} key={item.key}/>)}
    </div>
  )
}

export default SidenavItemList