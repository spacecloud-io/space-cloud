import React from 'react'
import { Select } from 'antd';
import mysqlSmall from '../../assets/mysqlSmall.svg'
import postgresqlSmall from '../../assets/postgresSmall.svg'
import mongodbSmall from '../../assets/mongoSmall.svg'
import './db-selector.css'

const { Option } = Select;

function DbSelector(props) {
  return (
    <div className="db-dropdown">
      <Select
        style={{ width: 200, }}
        placeholder="Select a database"
        onChange={props.handleSelect}
        value={props.selectedDb}
      >
        <Option value="sql-mysql"><img src={mysqlSmall} alt="mySQL" className="drop-icon"/>  MySQL</Option>
        <Option value="sql-postgres"><img src={postgresqlSmall} alt="postgresSQl" className="drop-icon"/> PostgreSQL</Option>
        <Option value="mongo"><img src={mongodbSmall} alt="mongoDB" className="drop-icon"/> MongoDB</Option>
      </Select>,
    </div>
  )
}

export default DbSelector