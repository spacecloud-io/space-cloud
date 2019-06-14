import React, { Component } from 'react';
import './rules.css';
import { Row, Col, Icon } from 'antd';
import { Controlled as CodeMirror } from 'react-codemirror2';
import 'codemirror/theme/material.css';
import 'codemirror/lib/codemirror.css';

class Rules extends Component {
	constructor(props) {
		super(props);
		const rules = Object.keys(this.props.rules);
		var selectedRule = '';
		if (rules.length > 0) {
			selectedRule = rules[0];
		}
		this.state = { selectedRule: selectedRule };
		this.handleChange = this.handleChange.bind(this);
		this.handleClick = this.handleClick.bind(this);
	}
	handleChange(editor, data, value) {
		console.log(value);
		this.props.handleRuleChange(value);
	}
	handleClick(rule) {
		this.setState({ selectedRule: rule });
	}
	render() {
		var rules = Object.keys(this.props.rules);
		const values = Object.values(this.props.rules);
		const index = rules.indexOf(this.state.selectedRule);

		return (
			<div className="rules-main-wrapper">
				<Row>
					<Col span={6}>
						<div className="addaRule" onClick={this.props.handleAddTextClick}>
							<Icon className="addIcon" type="plus" /> {this.props.addText}
						</div>
						<div className="rulesTable">
							{rules.map((rule) => {
								return (
									<div
										className={`rule ${this.state.selectedRule === rule ? 'selected' : ''}`}
										id="rule"
										value={rule}
										key={rule}
										onClick={() => this.handleClick(rule)}
									>
										{rule}
									</div>
								);
							})}
						</div>
					</Col>
					<Col span={18}>
						<div className="code">
							<div className="code-hint">
								Hint : To indent press ctrl + A in the editor and then shift + tab
							</div>
							<div className="code-mirror">
								<CodeMirror
									value={values[index]}
									options={{
										mode: 'xml',
										// theme: 'material',
										autofocus: true,
										lineNumbers: true
									}}
									onBeforeChange={(editor, data, value) => {
										this.handleChange(editor, data, value);
									}}
								/>
							</div>
						</div>
					</Col>
				</Row>
			</div>
		);
	}
}

export default Rules;
