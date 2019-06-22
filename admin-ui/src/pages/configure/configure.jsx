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
import { get, set } from 'automate-redux';
import store from ".././../store";

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
					<SslConfigureForm formState={props.ssl} handleChange={props.handleSslChange} />
					<Divider />
					<RealtimeConfigure formState={props.realtime} handleChange={props.handleRealtimeChange} />
					<Divider />
					<FunctionConfigure formState={props.functions} handleChange={props.handleFunctionChange} />
					<Divider />
					<FileStorage formState={props.fileStorage} handleChange={props.handleFileStorageChange} />
				</div>
			</div>
		</div>
	);
}

const mapStateToProps = (state, ownProps) => {
	return {
		secret: get(state, "config.secret"),
		ssl: get(state, "config.ssl", {}),
		realtime: get(state, "config.modules.realtime", {}),
		functions: get(state, "config.modules.functions", {}),
		fileStorage: get(state, "config.modules.fileStore", {})
	};
};

const mapDispatchToProps = (dispatch) => {
	return {
		handleSecretChange: (value) => {
			dispatch(set("config.secret", value))
		},
		handleSslChange: (value) => {
			const config = get(store.getState(), "config.ssl", {})
			dispatch(set("config.ssl", Object.assign({}, config, value)))
		},
		handleRealtimeChange: (value) => {
			const config = get(store.getState(), "config.modules.realtime", {})
			dispatch(set("config.modules.realtime", Object.assign({}, config, value)))
		},
		handleFunctionChange: (value) => {
			const config = get(store.getState(), "config.modules.functions", {})
			dispatch(set("config.modules.functions", Object.assign({}, config, value)))
		},
		handleFileStorageChange: (value) => {
			const config = get(store.getState(), "config.modules.fileStore", {})
			dispatch(set("config.modules.fileStore", Object.assign({}, config, value)))
		}
	};
};

export default connect(mapStateToProps, mapDispatchToProps)(Rules);
