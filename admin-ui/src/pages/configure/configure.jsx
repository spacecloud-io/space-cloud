import React from 'react';
import { connect } from 'react-redux';
// import '../../../index.css'
import Sidenav from '../../components/sidenav/Sidenav';
import Topbar from '../../components/topbar/Topbar';
import Header from '../../components/header/Header';
import SecretConfigure from '../../components/configure/SecretConfigure';
import { Divider } from 'antd';
import SslConfigureForm from '../../components/configure/SslConfigure';
import RealtimeConfigure from '../../components/configure/RealtimeConfigure';
import FunctionConfigure from '../../components/configure/FunctionConfigure';
import FileStorage from '../../components/configure/FileStorageConfigure';
function Rules(props) {
	return (
		<div>
			<Topbar title="Configure" />
			<div className="flex-box">
				<Sidenav selectedItem="configure" />
				<div className="page-content">
					<div className="header-flex">
						<Header name="Project Configurations" color="#000" fontSize="22px" />
					</div>
					<SecretConfigure formState={props.secret} handleChange={props.handleSecretChange} />
					<Divider />
					<SslConfigureForm formState={props.sslFormState} handleChange={props.handleSslChange} />
					<Divider />
					<RealtimeConfigure formState={props.realtimeFormState} handleChange={props.handleRealtimeChange} />
					<Divider />
					<FunctionConfigure formState={props.functionFormState} handleChange={props.handleFunctionChange} />
					<Divider />
					<FileStorage formState={props.fileStorageFormState} handleChange={props.handleFileStorageChange} />
				</div>
			</div>
		</div>
	);
}

const mapStateToProps = (state, ownProps) => {
	return {
		secret: '',
		sslFormState: { cert: '', enabled: false, key: '' },
		realtimeFormState: { broker: 'nats', enabled: true, conn: 'a' },
		functionFormState: { broker: 'nats', enabled: true, conn: 'a' },
		fileStorageFormState: { storage: '', enabled: true, conn: 'a' }
	};
};

const mapDispatchToProps = (dispatch) => {
	return {
		handleSecretChange: (value) => {
			console.log('secret changed', value);
		},
		handleSslChange: (value) => {
			console.log('ssl changed', value);
		},
		handleRealtimeChange: (value) => {
			console.log('realtime changed', value);
		},
		handleFunctionChange: (value) => {
			console.log('realtime changed', value);
		},
		handleFileStorageChange: (value) => {
			console.log('realtime changed', value);
		}
	};
};

export default connect(mapStateToProps, mapDispatchToProps)(Rules);
