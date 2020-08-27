import React, {useEffect} from "react";
import {Route, useHistory} from "react-router-dom";
import pages from "pages";
import styles from "styles.module.css";
import {useDispatch, useSelector} from "react-redux";

function getSystem(state) {
  return state.system;
}

function getAccount(state) {
  return state.account;
}

function App() {
  const system = useSelector(state => getSystem(state));
  const account = useSelector(state => getAccount(state));

  const dispatch = useDispatch();
  const history = useHistory();

  // useEffect(() => {
  //   if (!system.unlocked) {
  //     history.push("/unlock");
  //   }
  // }, [system.locked, account.status]);

  return (
  < div className = {
    styles.swarmcity
  } > {
    " "
  } {
    pages.map(({path, exact, component}) => (
    < Route key = {
      path
    } {
      ...{
        path,
        exact,
        component
      }
    } />))
  } {
    " "
  } < /div>
    );
}

export default App;
