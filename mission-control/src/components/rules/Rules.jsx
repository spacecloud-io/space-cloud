import React, { useState, useEffect } from 'react';
import './rules.css';
import { Row, Col, Icon } from 'antd';
import { Controlled as CodeMirror } from 'react-codemirror2';
import 'codemirror/theme/material.css';
import 'codemirror/lib/codemirror.css';
import 'codemirror/mode/javascript/javascript'
import 'codemirror/addon/selection/active-line.js'
import 'codemirror/addon/edit/matchbrackets.js'
import 'codemirror/addon/edit/closebrackets.js'

const Rules = (props) => {
	const [selected, setSelected] = useState(null)

	useEffect(() => {
		if (props.array && props.rules.length) {
			setSelected(0)
		} else if (!props.array && Object.keys(props.rules).length) {
			setSelected(Object.keys(props.rules)[0])
		}
	}, [])

	const handleDeleteClick = (e, rule) => {
		e.stopPropagation()
		props.handleDeleteRule(rule)
	}
	var rules = props.array ? props.rules.map((_, index) => (`Rule ${index + 1}`)) : Object.keys(props.rules);

	return (
		<div className="rules-main-wrapper">
			<Row>
				<Col span={6}>
					<div className="addaRule" onClick={props.handleAddRuleClick}>
						<Icon className="addIcon" type="plus" /> {props.addText}
					</div>
					<div className="rulesTable">

						{rules.map((rule, index) => {
							return (
								<div
									className={`rule ${selected === (props.array ? index : rule) ? 'selected' : ''}`}
									id="rule"
									value={rule}
									key={rule}
									onClick={() => setSelected(props.array ? index : rule)}
								>
									<div className="add-a-rule">
										{rule}
										<i class="material-icons delete-icon" onClick={(e) => handleDeleteClick(e, props.array ? index : rule)}>delete</i>
									</div>
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
								value={props.rules[selected]}
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
									props.handleRuleChange(selected, value);
								}}
							/>
						</div>
					</div>
				</Col>
			</Row>
		</div>
	);
}

export default Rules;
