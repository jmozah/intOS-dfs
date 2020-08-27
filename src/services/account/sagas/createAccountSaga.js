import {call, put, select, fork} from "redux-saga/effects";
import EthCrypto from "eth-crypto";
import {getAccount} from "../selectors";
import axios from "axios";
import qs from "querystring";

const axi = axios.create({baseURL: "http://127.0.0.1:9090/v0/", timeout: 5000});

export default function* createAccountSaga(action) {
  console.log("create account saga started");
  try {
    const requestBody = {
      user: action.data.username,
      password: action.data.password
    };

    console.log("request: ", requestBody);
    const config = {
      headers: {
        "Content-Type": "application/x-www-form-urlencoded"
      }
    };

    const response = yield axi({method: "POST", url: "user/signup", config: config, data: qs.stringify(requestBody)});
    //const response = yield axi.post("http://127.0.0.1:9090/v0/user/signup", requestBody);

    console.log(response);

    // encrypt wallet0
    const userObject = {
      status: "accountSet",
      username: action.data.username,
      avatar: action.data.avatar,
      address: response.data.reference
    };

    yield put({type: "SET_ACCOUNT", data: userObject});
    // yield put({ type: 'SET_SYSTEM', data: { mnemonic: decryptedMnemonic } })
    // yield put({ type: 'SET_SYSTEM', data: { privatekey: decryptedPrivateKey.privateKey } })
    yield put({
      type: "SET_SYSTEM",
      data: {
        unlocked: true
      }
    });
    // const fairdrive = yield window.fairdrive.getFairdrive(decryptedMnemonic.toString())
    // if (fairdrive.res) {
    //     yield put({ type: 'SET_DRIVE', data: fairdrive.res })
    // } else {
    //     console.log(fairdrive.err)
    // }
  } catch (e) {
    console.log("error on timeout", e);
  }
}