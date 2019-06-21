import React from 'react';
import { connect } from 'react-redux';
import '../../../index.css';
import Sidenav from '../../../components/sidenav/Sidenav';
import Topbar from '../../../components/topbar/Topbar';
import Header from '../../../components/header/Header';
import Documentation from '../../../components/documentation/Documentation';
import EmptyState from '../../../components/rules/EmptyState';
import rulesImg from '../../../assets/rules.svg';
import RulesComponent from '../../../components/rules/Rules';

function Rules(props) {
	const noOfRules = Object.keys(props.rules).length
	return (
		<div>
			<Topbar title="File Storage" />
			<div className="flex-box">
				<Sidenav selectedItem="file-storage" />
				<div className="page-content">
					<div className="header-flex">
						<Header name="Rules" color="#000" fontSize="22px" />
						<Documentation url="https://spaceuptech.com/docs/file-storage" />
					</div>
					{noOfRules && <RulesComponent
						rules={props.rules}
						handleRuleChange={props.handleRuleChange}
						addText={'Add a rule'}
						handleAddRuleClick={props.handleAddRuleClick}
					/>}
					{!noOfRules && <EmptyState
						graphics={rulesImg} desc="Guard your data with rules that define who has access to it and how it is structured."
						buttonText="Add a rule"
						handleClick={props.handleAddRuleClick} />}

				</div>
			</div>
		</div>
	);
}

const mapStateToProps = (state, ownProps) => {
	return {
		rules: {
			Books: 'a',
			Todos: 'b',
			Users: 'c',
			xyz: 'd',
			sas: 'e',
			qwqw: 'f',
			rewwr: 'g'
		}
	};
};

const mapDispatchToProps = (dispatch) => {
	return {
		handleRuleChange: (value) => {
			console.log('Rule Changed', value);
		},
		handleAddRuleClick: () => {
			console.log('Clicked Add rule');
		}
	};
};

export default connect(mapStateToProps, mapDispatchToProps)(Rules);
