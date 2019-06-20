import React from 'react';
import { connect } from 'react-redux';
import '../../index.css'
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
					<Header name="Project Configurations" color="#000" fontSize="22px" />

					<SecretConfigure formState={props.secret} handleChange={props.handleSecretChange} />
					<Divider />
					<SslConfigureForm formState={props.ssl} handleChange={props.handleSslChange} />
					<Divider />
					<RealtimeConfigure formState={props.realtime} handleChange={props.handleRealtimeChange} />
					<Divider />
					<FunctionConfigure formState={props.function} handleChange={props.handleFunctionChange} />
					<Divider />
					<FileStorage formState={props.fileStorage} handleChange={props.handleFileStorageChange} />
				</div>
			</div>
		</div>
	);
}

const mapStateToProps = (state, ownProps) => {
	return {
		secret: '',
		ssl: { cert: '', enabled: false, key: '' },
		realtime: { broker: 'nats', enabled: true, conn: 'a' },
		function: { broker: 'nats', enabled: true, conn: 'a' },
		fileStorage: { storage: '', enabled: true, conn: 'a' }
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
