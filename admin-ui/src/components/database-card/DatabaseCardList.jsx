import React from 'react'
import DatabaseCard from './DatabaseCard'

function DatabaseCardList(props){
  return(
      <div>
      {props.cards.map((card) => 
      <DatabaseCard key={card.key} name={card.name} desc={card.desc} graphics={card.graphics} handleEnable={() => props.handleEnable(card.key)}/>)}
    </div>
  )
}

export default DatabaseCardList