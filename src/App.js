import React, {useEffect} from "react";
import {Route, useHistory} from "react-router-dom";
import pages from "pages";
import styles from "styles.module.css";
import {useDispatch, useSelector} from "react-redux";
import {createMuiTheme, ThemeProvider} from "@material-ui/core/styles";
import {green, orange} from "@material-ui/core/colors";
import {logIn, isLoggedIn} from "helpers/apiCalls";

const outerTheme = createMuiTheme({
  palette: {
    primary: {
      main: "#92e7fa",
      light: "#cecece",
      background: "#333333"
    },
    secondary: {
      main: green[500]
    }
  }
});

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

  const loginState = useEffect(() => {
    // define function async
    async function checkAccountStatus() {
      if (account.status === "noAccount") {
        history.push("/account-create");
      } else {
        // do the api all to see if the user is logged in
        const loggedIn = await isLoggedIn(account.username).catch(e => {
          return e;
        });
        console.log(loggedIn);
        //history.push("/drive/root");
        // when not logged in
        // history.push("/unlock");
        // const res = await api.checkIsLoggedIn()
      }
    }
    checkAccountStatus().catch(e => console.log(e));
  }, [account.status]);

  return (
  < ThemeProvider theme = {
    outerTheme
  } > {
    " "
  } < div className = {
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
  } < /div>{" "} <
        /ThemeProvider >);
}

export default App;
