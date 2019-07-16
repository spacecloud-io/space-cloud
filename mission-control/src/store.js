import { createStore } from "redux";
import { generateReducers } from "automate-redux";

// Initial state of redux
const initialState = {
  uiState: {
    login: {
      formState: {
        userName: { value: "" },
        password: { value: "" }
      },
      isLoading: false
    }
  },
  config: {},
  savedConfig: {}
};

// Generate reducers with the initial state and pass it to the redux store
export default createStore(generateReducers(initialState), window.__REDUX_DEVTOOLS_EXTENSION__ && window.__REDUX_DEVTOOLS_EXTENSION__());