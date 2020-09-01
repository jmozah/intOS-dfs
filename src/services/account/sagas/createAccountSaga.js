import {call, put, select, fork} from "redux-saga/effects";
import EthCrypto from "eth-crypto";
import {getAccount} from "../selectors";
import axios from "axios";
import qs from "querystring";

const axi = axios.create({baseURL: "http://localhost:9090/v0/", timeout: 5000});

export default function* createAccountSaga(action) {
  console.log("create account saga started");
  try {
    const requestBody = {
      user: action.data.username,
      password: action.data.password,
      mnemonic: action.data.mnemonic
    };

    console.log("request: ", requestBody);

    const config = {
      headers: {
        "Content-Type": "application/x-www-form-urlencoded"
      }
    };

    const response = yield axi({method: "POST", url: "user/signup", config: config, data: qs.stringify(requestBody), withCredentials: true});
    //const response = yield axi.post("http://127.0.0.1:9090/v0/user/signup", requestBody);

    console.log(response);
    //console.log(response.headers["set-cookie"]);

    // encrypt wallet0
    const userObject = {
      status: "accountSet",
      username: action.data.username,
      avatar: action.data.avatar,
      address: response.data.reference,
      balance: 0.0
    };

    const podName = new Date().toISOString();

    const podRequest = {
      password: action.data.password,
      pod: "Fairdrive"
    };

    const createPod = yield axi({method: "POST", url: "pod/new", config: config, data: qs.stringify(podRequest), withCredentials: true});

    console.log(createPod);

    const createDocumentsDirectory = yield axi({
      method: "POST",
      url: "dir/mkdir",
      config: config,
      data: qs.stringify({dir: "Documents"}),
      withCredentials: true
    });

    const createMoviesDirectory = yield axi({
      method: "POST",
      url: "dir/mkdir",
      config: config,
      data: qs.stringify({dir: "Movies"}),
      withCredentials: true
    });

    const createMusicDirectory = yield axi({
      method: "POST",
      url: "dir/mkdir",
      config: config,
      data: qs.stringify({dir: "Music"}),
      withCredentials: true
    });

    const createPictursDirectory = yield axi({
      method: "POST",
      url: "dir/mkdir",
      config: config,
      data: qs.stringify({dir: "Pictures"}),
      withCredentials: true
    });

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