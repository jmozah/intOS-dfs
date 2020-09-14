import React, {useState, useEffect} from "react";
import styles from "styles.module.css";
import {useHistory} from "react-router-dom";

import {useDispatch, useSelector} from "react-redux";
import {Dialog} from "@material-ui/core";
import {Input} from "react-advanced-form-addons";

import {logIn, getAvatar} from "helpers/apiCalls";

function getAccount(state) {
  return state.account;
}

function getSystem(state) {
  return state.system;
}

export default function PasswordUnlock({open}) {
  const dispatch = useDispatch();
  const history = useHistory();
  const account = useSelector(state => getAccount(state));
  const system = useSelector(state => getSystem(state));

  const [password, setPassword] = useState();

  const handleSetPassword = e => {
    setPassword(e.target.value);
  };

  async function onLogin() {
    const isUserLoggedIn = await logIn(account.username, password);

    if (isUserLoggedIn.data.code == 200) {
      dispatch({
        type: "SET_SYSTEM",
        data: {
          passWord: password
        }
      });
      history.push("/drive/root");
    } else {
      console.log("user not logged in");
    }
  }

  function anotherAccount() {
    history.push("/login");
  }

  useEffect(() => {
    console.log(system);
    if (system.unlocked) {
      history.push("/drive/root");
    }
  });

  return (<div className={styles.dialogBox}>
    <div className={styles.flexer}></div>
    <div className={styles.flexer}></div>
    <div className={styles.flexer}></div>
    <div className={styles.flexer}></div>
    <div className={styles.title}>Unlock your account</div>
    <div className={styles.flexer}></div>
    <div className={styles.flexer}></div>
    <div className={styles.flexer}></div>

    <img src={account.avatar} className={styles.dialogAvatar}></img>

    <div className={styles.username}>{account.username}</div>
    <div className={styles.flexer}></div>
    <div className={styles.flexer}></div>

    <div className={styles.dialogPasswordBox}>
      <input className={styles.dialogPassword} type="password" placeholder="Password" onChange={e => handleSetPassword(e)}></input>
    </div>

    <div tabIndex="2" className={styles.button} onClick={onLogin}>
      <div>
        <div className={styles.buttontext}>continue</div>
      </div>
    </div>

    <div className={styles.flexer}></div>
    <div className={styles.link} onClick={anotherAccount}>
      Sign in with another account
    </div>
  </div>);
}
