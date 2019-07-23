import React from 'react';
import { Divider, Row, Col } from 'antd';
import './billing-table.css';
function BillingTable(props) {
	const reducer = (accumulator, currentValue) => accumulator + currentValue.amount;

	const total = props.data.reduce(reducer, 0);

	return (
		<div className="billing-table-main-wrapper">
			<div className="title">{props.title}</div>
			<hr className="divider" />
			<div className="table">
				{props.data.map((value) => {
					return (
						<Row>
							<Col xs={{ span: 12 }}>
								<div className="table-col">{value.item}</div>
							</Col>
							<Col xs={{ span: 8 }}>
								<div className="table-col">{value.usage}</div>
							</Col>
							<Col xs={{ span: 4 }}>
								<div className="table-col">$ {value.amount}</div>
							</Col>
						</Row>
					);
				})}
				<Row>
					<Col span={20} />
					<Col xs={{ span: 4 }}>
						<div className="total">
							Total: $
							<div className="total-amount">&nbsp;{total}</div>
						</div>
					</Col>
				</Row>
			</div>
		</div>
	);
}
export default BillingTable;
