import React from 'react';
import { connect } from 'react-redux';
import store from '../../../store'
import '../../../index.css';
import Sidenav from '../../../components/sidenav/Sidenav';
import Topbar from '../../../components/topbar/Topbar';
import Header from '../../../components/header/Header';
import Documentation from '../../../components/documentation/Documentation';
import EmptyState from '../../../components/rules/EmptyState';
import rulesImg from '../../../assets/rules.svg';
import RulesComponent from '../../../components/rules/Rules';
import { get, set, push } from "automate-redux";

const Rules = (props) => {
	const noOfRules = props.rules.length
	return (
		<div>
			<Topbar showProjectSelector />
			<div className="flex-box">
				<Sidenav selectedItem="static" />
				<div className="page-content">
					<div className="header-flex">
						<Header name="Rules" color="#000" fontSize="22px" />
						<Documentation url="https://spaceuptech.com/docs/static" />
					</div>
					{noOfRules > 0 && <RulesComponent
						array={true}
						rules={props.rules}
						handleRuleChange={props.handleRuleChange}
						addText={'Add a rule'}
						handleAddRuleClick={props.handleCreateRule}
					/>}
					{!noOfRules && <EmptyState
						graphics={rulesImg} desc="Guard your data with rules that define who has access to it and how it is structured."
						buttonText="Add a rule"
						handleClick={props.handleCreateRule} />}
				</div>
			</div>
		</div>
	)
}

const mapStateToProps = (state) => {
	return {
		rules: get(state, `config.modules.static.rules`, []),
	}
};

const mapDispatchToProps = (dispatch) => {
	return {
		handleRuleChange: (index, value) => {
			let rules = get(store.getState(), "config.modules.static.rules", []).slice()
			rules[index] = value
			dispatch(set(`config.modules.static.rules`, rules))
		},
		handleCreateRule: () => {
			const defaultRule = {
				prefix: "/",
				proxy: "http://localhost:3000"
			}
			dispatch(push(`config.modules.static.rules`, JSON.stringify(defaultRule, null, 2)))
		},
	};
};

export default connect(mapStateToProps, mapDispatchToProps)(Rules);
