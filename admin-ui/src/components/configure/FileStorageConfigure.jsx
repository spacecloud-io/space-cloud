import React from 'react';
import './configure.css';
import { Form, Input, Select, Switch } from 'antd';
import { createFormField } from 'rc-form';
const { Option } = Select;

function FileStorageConfigure(props) {
	const { getFieldDecorator } = props.form;
	return (
		<div className="configure">
			<div className="conn-string">FileStorage : </div>

			<Form className="conn-form" layout="inline">
				<div className="conn-form-options">
					<Form.Item>
						{getFieldDecorator('storeType', {
							rules: [ { required: true, message: '' } ]
						})(
							<Select placeholder="Store Type" className="select">
								<Option value="local">Local</Option>
								<Option value="s3">S3</Option>
							</Select>
						)}
					</Form.Item>
					<Form.Item label="Enabled" className="switch">
						{getFieldDecorator('enabled', { valuePropName: 'checked' })(<Switch size="small" />)}
					</Form.Item>
				</div>
				<div className="conn-form-cert">
					<Form.Item className="conn-form-cert-input">
						{getFieldDecorator('conn', {
							rules: [ { required: true, message: '' } ]
						})(<Input style={{ width: 600 }} placeholder="Enter Connection String" />)}
					</Form.Item>
				</div>
			</Form> <br />
		</div>
	);
}

const WrappedFileStorageConfigureForm = Form.create({
	mapPropsToFields(props) {
		return {
			storeType: createFormField({ value: props.formState.storeType }),
			enabled: createFormField({ value: props.formState.enabled }),
			conn: createFormField({ value: props.formState.conn })
		};
	},
	onValuesChange(props, changedValues) {
		props.handleChange(changedValues);
	}
})(FileStorageConfigure);

export default WrappedFileStorageConfigureForm;
