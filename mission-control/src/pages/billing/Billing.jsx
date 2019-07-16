import React from 'react';
import { connect } from 'react-redux';
import { get, increment } from 'automate-redux';
import store from '../../store';
import service from '../../index';
import { generateId, notify, isUserSignedIn, triggerSignin } from "../../utils"
import { PAYU_MERCHANT_KEY } from "../../constants"

import { Tabs } from 'antd';
import './billing.css';
import Sidenav from '../../components/sidenav/Sidenav';
import Topbar from '../../components/topbar/Topbar';
import TotalCredit from '../../components/billing/total-credit/TotalCredit';
import BillingTable from '../../components/billing/billing-table/BillingTable';
import BillingHistory from '../../components/billing/billing-history-table/BillingHistoryTable';
import Upgrade from "../../components/upgrade/Upgrade";
import RechargeModal from "../../components/billing/recharge-modal/RechargeModal";
import EmptyState from "../../components/empty-state/EmptyState"
import someGraphics from '../../assets/projectId.svg'
const { TabPane } = Tabs;

class Billing extends React.Component {
	constructor(props) {
		super(props)
		this.state = { isRechargeModalVisible: false }
	}

	handleRechargeModalVisibility = visible => {
		this.setState({ isRechargeModalVisible: visible })
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
								<RechargeModal visible={this.state.isRechargeModalVisible}
									handleCancel={() => this.handleRechargeModalVisibility(false)}
									handleSubmit={(option) => this.props.handleRecharge(option)} />
								<div>
									<TotalCredit amount={this.props.totalCredit} handleClick={() => this.handleRechargeModalVisibility(true)} />
									<div className="interval-text"></div>
									<BillingTable data={this.props.billing} title={'Usage this month'} />
								</div>
							</React.Fragment>
							:
							<EmptyState
								graphics={someGraphics}
								handleClick={triggerSignin}
								desc="Lorem ipsum dolor sit amet consectetur adipisicing elit. Vitae qui id nulla ipsa maiores fugit ipsum inventore esse iste magnam. Porro blanditiis possimus animi voluptatum? Similique vel illo at asperiores."
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
		handleRecharge: option => {
			// Note: Amount is to be expressed in stringified float with double digit precision otherwise PAYU hash does not matches 
			let amount = "10.00", extraCredits = 0
			switch (option) {
				case 1:
					amount = "100.00"
					extraCredits = 10
					break
				case 2:
					amount = "1000.00"
					extraCredits = 150
					break
			}
			//Create a Data object that is to be passed to LAUNCH method of Bolt
			let pd = {
				key: PAYU_MERCHANT_KEY, txnid: generateId().substring(0, 29),
				amount: amount, udf1: extraCredits, udf2: '', udf3: '', udf4: '', udf5: '',
				firstname: get(store.getState(), "user.name").split(" ")[0], email: get(store.getState(), "user.email"),
				phone: '6111111111', productinfo: `Space Cloud Enterprise Credits`, hash: '',
				surl: 'https://www.payumoney.com/merchant-dashboard/#/transactions',
				furl: 'https://www.payumoney.com/merchant-dashboard/#/transactions'
			}

			// Data to be Sent to API to generate hash.
			let data = {
				txnid: pd.txnid, email: pd.email, amount: pd.amount, productinfo: pd.productinfo, firstname: pd.firstname, udf1: pd.udf1
			}

			service.startPayment(data).then(hash => {
				pd.hash = hash
				window.bolt.launch(pd, {
					responseHandler: function (response) {
						// your payment response Code goes here
						if (response.response.txnStatus !== "CANCEL") {
							service.completePayment(Object.assign({}, response.response, { udf1: pd.udf1 })).then(() => {
								dispatch(increment("credits", Number(amount) + extraCredits))
								notify("success", "Success", "Payment was successfull")
							}).catch(ex => {
								console.log("Error", ex)
								notify("error", "Error", "Payment failed")
							})
						}
					},
					catchException: function (response) {
						console.log("Exception", response)
						notify("error", "Error", "Something went wrong")
					}
				});

			}).catch(ex => {
				console.log("Error", ex)
				notify("error", "Error", "Payment could not be initiated")
			})
		},
		handleRangeChange: (dateStringprops) => {
			console.log('range changed', dateStringprops);
		}
	};
};

export default connect(mapStateToProps, mapDispatchToProps)(Billing);
