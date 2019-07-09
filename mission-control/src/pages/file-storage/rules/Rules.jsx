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
import EditItemModal from "../../../components/edit-item-modal/EditItemModal";
import projectId from '../../../assets/projectId.svg'
import { get, set } from "automate-redux";

class Rules extends React.Component {
	constructor(props) {
		super(props)
		this.state = { modalVisible: false }
		this.handleModalVisiblity = this.handleModalVisiblity.bind(this);
	}

	handleModalVisiblity(visible) {
		this.setState({ modalVisible: visible });
	}
	render() {
		const noOfRules = Object.keys(this.props.rules).length
		return (
			<div>
				<Topbar showProjectSelector />
				<div className="flex-box">
					<Sidenav selectedItem="file-storage" />
					<div className="page-content">
						<div className="header-flex">
							<Header name="Rules" color="#000" fontSize="22px" />
							<Documentation url="https://spaceuptech.com/docs/file-storage" />
						</div>
						{noOfRules > 0 && <RulesComponent
							rules={this.props.rules}
							handleRuleChange={this.props.handleRuleChange}
							addText={'Add a rule'}
							handleAddRuleClick={() => this.handleModalVisiblity(true)}
						/>}
						{!noOfRules && <EmptyState
							graphics={rulesImg} desc="Guard your data with rules that define who has access to it and how it is structured."
							buttonText="Add a rule"
							handleClick={() => this.handleModalVisiblity(true)} />}
						<EditItemModal graphics={projectId} heading="Rule name" name="Give a rule name" desc="This name is just to identify a unique rule" placeholder="Enter a rule name" visible={this.state.modalVisible} handleCancel={() => this.handleModalVisiblity(false)} handleSubmit={this.props.handleCreateRule} />
					</div>
				</div>
			</div>
		);
	}
}

const mapStateToProps = (state, ownProps) => {
	return {
		rules: get(state, `config.modules.fileStore.rules`, {}),
	}
};

const mapDispatchToProps = (dispatch) => {
	return {
		handleRuleChange: (ruleName, value) => {
			dispatch(set(`config.modules.fileStore.rules.${ruleName}`, value))
		},
		handleCreateRule: (ruleName) => {
			const defaultRule = {
				prefix: "/",
				rule: {
					create: {
						rule: "allow"
					},
					read: {
						rule: "allow"
					},
					delete: {
						rule: "allow"
					}
				}
			}
			dispatch(set(`config.modules.fileStore.rules.${ruleName}`, JSON.stringify(defaultRule, null, 2)))
		},
	};
};

export default connect(mapStateToProps, mapDispatchToProps)(Rules);
