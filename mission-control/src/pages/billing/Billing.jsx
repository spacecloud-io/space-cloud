import React from 'react';
import { connect } from 'react-redux';
import './billing.css';
import Sidenav from '../../components/sidenav/Sidenav';
import Topbar from '../../components/topbar/Topbar';
import { Tabs } from 'antd';
import TotalCredit from '../../components/billing/total-credit/TotalCredit';
import BillingTable from '../../components/billing/billing-table/BillingTable';
import BillingHistory from '../../components/billing/billing-history-table/BillingHistoryTable';
const { TabPane } = Tabs;

function Billing(props) {
	const reducer = (accumulator, currentValue) => accumulator + currentValue.amount;

	var total = props.enterprise.reduce(reducer, 0);
	total = props.hosted.reduce(reducer, total);

	return (
		<div>
			<Topbar showProjectSelector />
			<div className="flex-box">
				<Sidenav selectedItem="billing" />
				<div className="page-content">
					<Tabs defaultActiveKey="billing-history">
						<TabPane tab="Expense" key="expense">
							<TotalCredit amount={props.totalCredit} />
							<div className="interval-text">{props.intervalText}</div>
							<BillingTable data={props.enterprise} title={'Enterprise'} />
							<BillingTable data={props.hosted} title={'Hosted'} />
							<div className="amount-wrapper">
								<div className="amount-text">Total Amount</div>
								<div className="amount">$ {total}</div>
							</div>
						</TabPane>
						<TabPane tab="Billing History" key="billing-history">
							<BillingHistory data={props.billingHistory} handleRangeChange={props.handleRangeChange} />
						</TabPane>
					</Tabs>
				</div>
			</div>
		</div>
	);
}

const mapStateToProps = (state, ownProps) => {
	return {
		intervalText: '1st - 31st May,2019',
		totalCredit: 10.53,
		enterprise: [
			{
				project: 'abcd',
				item: 'Number of Usage',
				usage: '20',
				amount: 10.67
			}
		],
		hosted: [
			{
				project: 'abcd',
				item: 'Number of Usage',
				usage: '20',
				amount: 10.67
			},
			{
				project: 'xyz',
				item: 'Number of Usage',
				usage: '20',
				amount: 10.49
			},
			{
				project: 'pqrs',
				item: 'Number of Usage',
				usage: '20',
				amount: 10.59
			}
		],
		billingHistory: [
			{
				interval: '21 Mar- 21 Apr',
				project: 'abcd',
				item: 'Number of Usage',
				usage: '20',
				amount: -10.67
			},
			{
				interval: '21 APr- 21 May',
				project: 'xyz',
				item: 'Number of Usage',
				usage: '20',
				amount: 10.49
			},
			{
				interval: '21 May- 21 Jun',
				project: 'pqrs',
				item: 'Number of Usage',
				usage: '20',
				amount: 10.59
			}
		]
	};
};

const mapDispatchToProps = (dispatch) => {
	return {
		handleRangeChange: (dateStringprops) => {
			console.log('range changed', dateStringprops);
		}
	};
};

export default connect(mapStateToProps, mapDispatchToProps)(Billing);
