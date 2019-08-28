import React from 'react';
import ReactGA from 'react-ga';
import { connect } from 'react-redux';
import { get, increment } from 'automate-redux';
import store from '../../store';
import service from '../../index';
import { notify, isUserSignedIn, triggerSignin } from "../../utils"

import { Tabs } from 'antd';
import './billing.css';
import Sidenav from '../../components/sidenav/Sidenav';
import Topbar from '../../components/topbar/Topbar';
import TotalCredit from '../../components/billing/total-credit/TotalCredit';
import BillingTable from '../../components/billing/billing-table/BillingTable';
import BillingHistory from '../../components/billing/billing-history-table/BillingHistoryTable';
import RechargeModal from "../../components/billing/recharge-modal/RechargeModal";
import EmptyState from "../../components/empty-state/EmptyState"
import someGraphics from '../../assets/projectId.svg'
import Upgrade from '../../components/billing/upgrade/Upgrade'
const { TabPane } = Tabs;

class Billing extends React.Component {
	constructor(props) {
		super(props)
		this.state = { isRechargeModalVisible: false }
	}

	handleRechargeModalVisibility = visible => {
		this.setState({ isRechargeModalVisible: visible })
	}

	componentDidMount() {
		ReactGA.pageview("/billing");
	}
	render() {
		const reducer = (accumulator, currentValue) => accumulator + currentValue.amount;
		const total = this.props.billing.reduce(reducer, 0)
		return (
			<div>
				<Topbar showProjectSelector />
				<div className="flex-box">
					<Sidenav selectedItem="billing" />
					<div className="page-content">
						{this.props.signedIn ?
							<React.Fragment>
								{/* <RechargeModal visible={this.state.isRechargeModalVisible}
									handleCancel={() => this.handleRechargeModalVisibility(false)}
									handleSubmit={(option) => this.props.handleRecharge(option)} /> */}
								<div>
									<TotalCredit amount={this.props.totalCredit} handleClick={() => this.props.handleRechargeClick()} />
									<div className="interval-text"></div>
									{this.props.mode === 0 && <Upgrade />}
									<BillingTable data={this.props.billing} title={'Usage this month'} />
								</div>
							</React.Fragment>
							:
							<EmptyState
								graphics={someGraphics}
								handleClick={triggerSignin}
								desc="With Space Cloud Enterprise you get the power to experience Firebase + Heroku on your Kubernetes cluster and much more. Signin to access the Space Cloud Enterprise and unlock new powers!"
								actionText="Register / Login" />
						}
					</div>
				</div>
			</div>
		)
	}
}


const mapStateToProps = (state, ownProps) => {
	return {
		mode: get(state, "operationConfig.mode", 0),
		signedIn: isUserSignedIn(),
		totalCredit: get(state, "credits", 0),
		billing: get(state, "billing", []).map(obj => ({
			amount: obj.price,
			item: `Space Cloud (${obj.mode === 1 ? 'Standard' : 'Premium'} plan)`,
			usage: obj.hours + " hours"
		}))
	};
};

const mapDispatchToProps = (dispatch) => {
	return {
		handleRechargeClick: () => {
			const { email, name } = get(store.getState(), "user", {})
			service.requestPayment(email, name).then(() => {
				notify("success", "Hey Buddy", "We are excited that you want to pay! You will receive an email within a day from our team to guide you through the next steps for payment", 20)
			}).catch(ex => {
				console.log("Error", ex)
				notify("error", "Error", "Error requesting payment")
			})
		},
		handleRangeChange: (dateStringprops) => {
			console.log('range changed', dateStringprops);
		}
	};
};

export default connect(mapStateToProps, mapDispatchToProps)(Billing);
