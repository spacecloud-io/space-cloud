import React from 'react';
import './FileStorageConfigure.css';
import { Form, Input, Select, Switch } from 'antd';
import { createFormField } from 'rc-form';
const { Option } = Select;

function FileStorageConfigure(props) {
	const { getFieldDecorator } = props.form;
	return (
		<div className="fileStorage-configuration">
			<div className="conn-string">FileStorage</div>

			<Form className="conn-form" layout="inline">
				<div className="conn-form-options">
					<Form.Item>
						{getFieldDecorator('storage', {
							rules: [ { required: true, message: '' } ]
						})(
							<Select placeholder="Broker" className="select">
								<Option value="localstore">LocalStore</Option>
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
						})(<Input style={{ width: 350 }} placeholder="Enter Connection String" />)}
					</Form.Item>
				</div>
			</Form>
		</div>
	);
}

const WrappedFileStorageConfigureForm = Form.create({
	mapPropsToFields(props) {
		return {
			storage: createFormField({ value: props.formState.storage }),
			enabled: createFormField({ value: props.formState.enabled }),
			conn: createFormField({ value: props.formState.conn })
		};
	},
	onValuesChange(props, changedValues) {
		props.handleChange(changedValues);
	}
})(FileStorageConfigure);

export default WrappedFileStorageConfigureForm;
