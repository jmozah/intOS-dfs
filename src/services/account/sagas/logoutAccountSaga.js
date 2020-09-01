import {call, put, select, fork} from "redux-saga/effects";
import EthCrypto from "eth-crypto";
import {getAccount} from "../selectors";
import axios from "axios";
import qs from "querystring";

const axi = axios.create({baseURL: "http://localhost:9090/v0/", timeout: 5000});

export default function* logoutAccountSaga(action) {
  console.log("logout account saga started");
  try {
    const config = {
      headers: {
        "Content-Type": "application/x-www-form-urlencoded"
      }
    };

    const response = yield axi({method: "POST", url: "user/logout", config: config, withCredentials: true});
    //const response = yield axi.post("http://127.0.0.1:9090/v0/user/signup", requestBody);

    console.log(response);
    //console.log(response.headers["set-cookie"]);
  } catch (e) {
    console.log("error on timeout", e);
  }
}