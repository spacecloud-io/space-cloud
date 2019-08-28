import React, { Component } from 'react';
import { DatePicker, Row, Col } from 'antd';
import moment from 'moment';
import './billing-history-table.css';
const { RangePicker } = DatePicker;
const date = new Date().getDate(); //Current Date
const month = new Date().getMonth() + 1; //Current Month
const year = new Date().getFullYear(); //Current Year
class BillingHistoryTable extends Component {
	constructor(props) {
		super(props);
		this.state = {
			rangeDateStart: `${year}/${month - 3}/${date}`,
			rangeDateEnd: `${year}/${month}/${date}`
		};
		this.handleChange = this.handleChange.bind(this);
	}
	handleChange(date, dateStringprops) {
		console.log(date);
		console.log(dateStringprops[0]);
		this.setState({ rangeDateStart: dateStringprops[0], rangeDateEnd: dateStringprops[1] });
		this.props.handleRangeChange(dateStringprops);
	}
	render() {
		const monthWordStart = new Date(this.state.rangeDateStart).toLocaleString('en-us', { month: 'short' });
		const monthWordEnd = new Date(this.state.rangeDateEnd).toLocaleString('en-us', { month: 'short' });
		var dateStart = this.state.rangeDateStart.split('/');
		var dateEnd = this.state.rangeDateEnd.split('/');
		const reducer = (accumulator, currentValue) => accumulator + currentValue.amount;
		const total = this.props.data.reduce(reducer, 0);
		return (
			<div className="billing-history">
				<div className="table-header">
					<div className="selected-range">
						{monthWordStart} {dateStart[2]} {dateStart[0]} - {monthWordEnd} {dateEnd[2]},{dateEnd[0]}
					</div>
					<RangePicker
						className="range-picker"
						onChange={this.handleChange}
						defaultValue={[
							moment([ this.state.rangeDateStart ], 'YYYY/MM/DD'),
							moment([ this.state.rangeDateEnd ], 'YYYY/MM/DD')
						]}
						format={'YYYY/MM/DD'}
					/>
				</div>
				<div className="table">
					<div className="heading-row">
						<Row>
							<Col xs={{ span: 6, offset: 1 }}>
								<div className="table-heading">Interval</div>
							</Col>
							<Col xs={{ span: 4, offset: 0 }}>
								<div className="table-heading">Project</div>
							</Col>
							<Col xs={{ span: 4, offset: 0 }}>
								<div className="table-heading">Item</div>
							</Col>
							<Col xs={{ span: 4, offset: 1 }}>
								<div className="table-heading">Usage</div>
							</Col>
							<Col xs={{ span: 4, offset: 0 }}>
								<div className="table-heading">Amount</div>
							</Col>
						</Row>
					</div>
					{this.props.data.map((value) => {
						return (
							<div className="row">
								<Row>
									<Col xs={{ span: 6, offset: 1 }}>
										<div className="table-col">{value.interval}</div>
									</Col>
									<Col xs={{ span: 4, offset: 0 }}>
										<div className="table-col">{value.project}</div>
									</Col>
									<Col xs={{ span: 4, offset: 0 }}>
										<div className="table-col">{value.item}</div>
									</Col>
									<Col xs={{ span: 4, offset: 1 }}>
										<div className="table-col">{value.usage}</div>
									</Col>
									<Col xs={{ span: 4, offset: 0 }}>
										<div
											className={`table-col ${value.amount >= 0
												? 'amountPositive'
												: 'amountNegative'}`}
										>
											$ {value.amount}
										</div>
									</Col>
								</Row>
							</div>
						);
					})}
					<Row>
						<Col span={18} />
						<Col xs={{ span: 4, offset: 2 }}>
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
}

export default BillingHistoryTable;
