import {call, put, select, fork} from "redux-saga/effects";
import EthCrypto from "eth-crypto";
import {getAccount} from "../selectors";
import axios from "axios";
import qs from "querystring";

const axi = axios.create({baseURL: "http://localhost:9090/v0/", timeout: 5000});

export default function* loginAccountSaga(action) {
  console.log("login account saga started");
  try {
    const requestBody = {
      user: action.data.username,
      password: action.data.password
    };

    const config = {
      headers: {
        "Content-Type": "application/x-www-form-urlencoded"
      }
    };

    const response = yield axi({method: "POST", url: "user/login", config: config, data: qs.stringify(requestBody), withCredentials: true});
    yield put({
      type: "SET_SYSTEM",
      data: {
        unlocked: true
      }
    });
    console.log(response);
  } catch (e) {
    console.log("error", e.message);
  }
}