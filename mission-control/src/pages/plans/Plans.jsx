import React, { useState } from 'react';
import ReactGA from 'react-ga';
import { connect } from 'react-redux';
import { get, set, increment, decrement } from 'automate-redux';
import service from "../../index"
import store from "../../store";
import { notify, isUserSignedIn } from "../../utils"

import { Row, Col } from "antd"
import Sidenav from '../../components/sidenav/Sidenav';
import Topbar from '../../components/topbar/Topbar';
import Header from '../../components/header/Header';
import Plan from "../../components/plan/Plan";
import '../../index.css'

const NotifyContent = (props) => (
	<span>
		Follow <a href="https://spaceuptech.com/blog/firebase-and-heroku-on-kubernetes">this guide</a> to deploy apps on Kubernetes via Space Cloud
	</span>
)

function Plans({ mode, handleModeChange }) {
	useState(() => {
		ReactGA.pageview("/plans");
	}, [])
	return (
		<div>
			<Topbar showProjectSelector />
			<div className="flex-box">
				<Sidenav selectedItem="plans" />
				<div className="page-content">
					<Header name="Plans" color="#000" fontSize="22px" />
					<p>We strongly believe in open source. That's why all development related features in Space Cloud will be open source forever. Our paid plans help you move faster by automating devops and providing support.</p>
					<Row gutter={{ xs: 8, sm: 16, md: 24, lg: 32 }} type="flex" style={{ marginTop: '48px' }}>
						<Col span={8}>
							<Plan
								name="OpenSource"
								desc="Bootstrap your hobby projects with SC Open Source. Automate your backend"
								points={["All development features", "1 Project"]}
								pricing="Free forever" active={mode === 0} handleClick={() => handleModeChange(0)} />
						</Col>
						<Col span={8}>
							<Plan
								name="Standard"
								desc="Experience Firebase + Heroku on your Kubernetes cluster!"
								points={["Deploy to Kubernetes / Docker Swarm", "3 Projects"]}
								pricing="10$/hour/month" active={mode === 1} handleClick={() => handleModeChange(1)} />
						</Col>
						<Col span={8}>
							<Plan
								name="Premium"
								desc="Unlock all powers with the  Premium SC. Make enterprise ready apps"
								points={["Metrics + Reporting", "10 Projects"]}
								pricing="50$/hour/month" active={mode === 2} handleClick={() => handleModeChange(2)} />
						</Col>
					</Row>
				</div>
			</div>
		</div>
	);
}

const mapStateToProps = (state) => {
	return {
		mode: get(state, "operationConfig.mode", 0)
	};
};

const mapDispatchToProps = (dispatch) => {
	return {
		handleModeChange: (mode) => {
			if (!isUserSignedIn() && mode > 0) {
				dispatch(set("uiState.isSigninModalVisible", true))
				return
			}
			dispatch(increment("pendingRequests"))
			const newOperationConfig = Object.assign({}, get(store.getState(), "operationConfig", {}), { mode: mode })
			service.saveOperationConfig(newOperationConfig).then(() => {
				dispatch(set("operationConfig", newOperationConfig))
				notify("success", "Success", "Plan changed successfully")
				if (mode > 0) {
					notify("info", "Next steps", <NotifyContent />)
				}
			}).catch(error => {
				console.log("Error", error)
				notify("error", "Error", 'Could not change mode')
			}).finally(() => dispatch(decrement("pendingRequests")))
		}
	};
};

export default connect(mapStateToProps, mapDispatchToProps)(Plans);
