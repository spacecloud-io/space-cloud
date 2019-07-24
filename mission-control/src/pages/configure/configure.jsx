import React, { useState } from 'react';
import ReactGA from 'react-ga';
import { connect } from 'react-redux';
import '../../index.css'
import Sidenav from '../../components/sidenav/Sidenav';
import Topbar from '../../components/topbar/Topbar';
import Header from '../../components/header/Header';
import SecretConfigure from '../../components/configure/SecretConfigure';
import { Divider } from 'antd';
import RealtimeConfigure from '../../components/configure/RealtimeConfigure';
import FunctionConfigure from '../../components/configure/FunctionConfigure';
import FileStorage from '../../components/configure/FileStorageConfigure';
import StaticConfigure from '../../components/configure/StaticConfigure';
import { get, set } from 'automate-redux';
import store from ".././../store";
import './configure.css'

function Rules(props) {
	useState(() => {
    ReactGA.pageview("/projects/configure");
  }, [])
	return (
		<div className="configurations">
			<Topbar showProjectSelector />
			<div className="flex-box">
				<Sidenav selectedItem="configure" />
				<div className="page-content">
					<Header name="Project Configurations" color="#000" fontSize="22px" />
					<SecretConfigure formState={props.secret} handleChange={props.handleSecretChange} />
					<Divider />
					<RealtimeConfigure formState={props.realtime} handleChange={props.handleRealtimeChange} />
					<Divider />
					<FunctionConfigure formState={props.functions} handleChange={props.handleFunctionChange} />
					<Divider />
					<FileStorage formState={props.fileStorage} handleChange={props.handleFileStorageChange} />
					<Divider />
					<StaticConfigure formState={props.static} handleChange={props.handleStaticChange} />
				</div>
			</div>
		</div>
	);
}

const mapStateToProps = (state, ownProps) => {
	return {
		secret: get(state, "config.secret"),
		realtime: get(state, "config.modules.realtime", {}),
		functions: get(state, "config.modules.functions", {}),
		fileStorage: get(state, "config.modules.fileStore", {}),
		static: get(state, "config.modules.static", {})
	};
};

const mapDispatchToProps = (dispatch) => {
	return {
		handleSecretChange: (value) => {
			dispatch(set("config.secret", value))
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
		},
		handleStaticChange: (value) => {
			const config = get(store.getState(), "config.modules.static", {})
			dispatch(set("config.modules.static", Object.assign({}, config, value)))
		}
	};
};

export default connect(mapStateToProps, mapDispatchToProps)(Rules);
