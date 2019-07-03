import React, { Component } from 'react';
import './rules.css';
import { Row, Col, Icon } from 'antd';
import { Controlled as CodeMirror } from 'react-codemirror2';
import 'codemirror/theme/material.css';
import 'codemirror/lib/codemirror.css';
import 'codemirror/mode/javascript/javascript'
import 'codemirror/addon/selection/active-line.js'
import 'codemirror/addon/edit/matchbrackets.js'
import 'codemirror/addon/edit/closebrackets.js'

class Rules extends Component {
	constructor(props) {
		super(props);
		const rules = Object.keys(this.props.rules);
		var selectedRule = '';
		if (rules.length > 0) {
			selectedRule = rules[0];
		}
		this.state = { selectedRule: selectedRule };

		this.handleClick = this.handleClick.bind(this);
	}

	handleClick(rule) {
		this.setState({ selectedRule: rule });
	}
	render() {
		var rules = Object.keys(this.props.rules);

		return (
			<div className="rules-main-wrapper">
				<Row>
					<Col span={6}>
						<div className="addaRule" onClick={this.props.handleAddRuleClick}>
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
									value={this.props.rules[this.state.selectedRule]}
									options={{
										mode: { name: "javascript", json: true },
										lineNumbers: true,
										styleActiveLine: true,
										matchBrackets: true,
										autoCloseBrackets: true,
										tabSize: 2,
										autofocus: true
									}}
									onBeforeChange={(editor, data, value) => {
										this.props.handleRuleChange(this.state.selectedRule, value);
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
