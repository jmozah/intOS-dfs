import {takeEvery} from "redux-saga/effects";
import * as t from "./actionTypes";
// Sagas
import createAccountSaga from "./sagas/createAccountSaga";

/******************************* Watchers *************************************/

export default function* accountRootSaga() {
  yield takeEvery(t.CREATE_ACCOUNT, createAccountSaga);
}
