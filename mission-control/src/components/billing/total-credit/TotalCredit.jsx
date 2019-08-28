import React from 'react';
import { Button } from 'antd';
import './total-credit.css'
function TotalCredit(props) {
	return (
		<div className="total-credit-main-wrapper">
			<div className="credit-wrapper">
				<div className="heading">Total Account Credit</div>
				<div className="amount">${props.amount}</div>
			</div>
      <div className="credit-text-wrapper">
        <div className="credit-text">Your credit</div>
      </div>
			<div className="button-wrapper">
				<Button type="primary" onClick={props.handleClick}>Pay Early</Button>
			</div>
		</div>
	);
}
export default TotalCredit;
