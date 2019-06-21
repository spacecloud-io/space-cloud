import React from 'react';
import { connect } from 'react-redux';
import '../../../index.css';
import Sidenav from '../../../components/sidenav/Sidenav';
import Topbar from '../../../components/topbar/Topbar';
import Header from '../../../components/header/Header';
import Documentation from '../../../components/documentation/Documentation';
import DbConfigure from '../../../components/database-rules/DbConfigure';
import EmptyState from '../../../components/rules/EmptyState';
import rulesImg from '../../../assets/rules.svg';
import RulesComponent from '../../../components/rules/Rules';

function Rules(props) {
	const noOfRules = Object.keys(props.rules).length
	return (
		<div>
			<Topbar title="Database" />
			<div className="flex-box">
				<Sidenav selectedItem="database" />
				<div className="page-content">
					<div className="header-flex">
						<Header name="Rules" color="#000" fontSize="22px" />
						<Documentation url="https://spaceuptech.com/docs/database" />
					</div>
					<DbConfigure updateFormState={props.updateFormState} formState={props.formState} />
					{noOfRules && <RulesComponent
						rules={props.rules}
						handleRuleChange={props.handleRuleChange}
						addText={'Add a table rule'}
						handleAddRuleClick={props.handleAddRuleClick}
					/>}
					{!noOfRules && <EmptyState
						graphics={rulesImg} desc="Guard your data with rules that define who has access to it and how it is structured."
						buttonText="Add a table"
						handleClick={props.handleAddRuleClick} />}


				</div>
			</div>
		</div>
	);
}

const mapStateToProps = (state, ownProps) => {
	return {
		formState: {
			enabled: true,
			conn: 'http://localhost/8080'
		},
		rules: {
			Books: 'a',
			Todos: 'b',
			Users: 'c',
			xyz: 'd',
			sas: 'e',
			qwqw: 'f',
			rewwr: 'g'
		},

	};
};

const mapDispatchToProps = (dispatch) => {
	return {
		handleRuleChange: (value) => {
			console.log('Rule Changed', value);
		},
		handleAddRuleClick: () => {
			console.log('Clicked Add rule');
		},
		updateFormState: (fields) => {
			console.log(fields);
		}
	};
};

export default connect(mapStateToProps, mapDispatchToProps)(Rules);
