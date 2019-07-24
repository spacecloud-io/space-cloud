import React from 'react';
import ReactGA from 'react-ga';
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
import EditItemModal from "../../../components/edit-item-modal/EditItemModal";
import projectId from '../../../assets/projectId.svg'
import { get, set } from "automate-redux";
import store from "../../../store"

class Rules extends React.Component {
	constructor(props) {
		super(props)
		this.state = { modalVisible: false }
		this.handleModalVisiblity = this.handleModalVisiblity.bind(this);
	}

	componentDidMount() {
		ReactGA.pageview("/projects/database/rules");
	}

	handleModalVisiblity(visible) {
		this.setState({ modalVisible: visible });
	}
	render() {
		const noOfRules = Object.keys(this.props.rules).length
		return (
			<div>
				<Topbar showProjectSelector showDbSelector selectedDb={this.props.selectedDb}/>
				<div className="flex-box">
					<Sidenav selectedItem="database" />
					<div className="page-content">
						<div className="header-flex">
							<Header name="Rules" color="#000" fontSize="22px" />
							<Documentation url="https://spaceuptech.com/docs/database" />
						</div>
						<DbConfigure updateFormState={this.props.updateFormState} formState={this.props.formState} />
						{noOfRules > 0 && <RulesComponent
							rules={this.props.rules}
							handleRuleChange={this.props.handleRuleChange}
							addText={'Add a table rule'}
							handleAddRuleClick={() => this.handleModalVisiblity(true)}
							handleDeleteRule={this.props.handleDeleteRule}
						/>}
						{!noOfRules && <EmptyState
							graphics={rulesImg} desc="Guard your data with rules that define who has access to it and how it is structured."
							buttonText="Add a table"
							handleClick={() => this.handleModalVisiblity(true)} />}
						<EditItemModal graphics={projectId} heading="Table name" name="Give a table name" desc="Note: This doesn't actually creates a table. It's for writing rules for a table" placeholder="Enter a table name" visible={this.state.modalVisible} handleCancel={() => this.handleModalVisiblity(false)} handleSubmit={this.props.handleCreateRule} />
					</div>
				</div>
			</div >
		);
	}
}

const mapStateToProps = (state, ownProps) => {
	return {
		selectedDb: ownProps.match.params.database,
		formState: {
			enabled: get(state, `config.modules.crud.${ownProps.match.params.database}.enabled`, false),
			conn: get(state, `config.modules.crud.${ownProps.match.params.database}.conn`),
		},
		rules: get(state, `config.modules.crud.${ownProps.match.params.database}.collections`, {}),
	};
};

const mapDispatchToProps = (dispatch, ownProps) => {
	const selectedDb = ownProps.match.params.database
	return {
		handleRuleChange: (ruleName, value) => {
			dispatch(set(`config.modules.crud.${selectedDb}.collections.${ruleName}`, value))
		},
		handleDeleteRule: (ruleName) => {
			const rules = get(store.getState(), `config.modules.crud.${selectedDb}.collections`)
			delete rules[ruleName]
			dispatch(set(`config.modules.crud.${selectedDb}.collections`, rules))
		},
		handleCreateRule: (ruleName) => {
			const defaultRule = {
				isRealtimeEnabled: true,
				rules: {
					create: {
						rule: "allow"
					},
					read: {
						rule: "allow"
					},
					update: {
						rule: "allow"
					},
					delete: {
						rule: "allow"
					}
				}
			}
			dispatch(set(`config.modules.crud.${selectedDb}.collections.${ruleName}`, JSON.stringify(defaultRule, null, 2)))
		},
		updateFormState: (fields) => {
			const dbConfig = get(store.getState(), `config.modules.crud.${selectedDb}`, {})
			dispatch(set(`config.modules.crud.${selectedDb}`, Object.assign({}, dbConfig, fields)))
		}
	};
};

export default connect(mapStateToProps, mapDispatchToProps)(Rules);
