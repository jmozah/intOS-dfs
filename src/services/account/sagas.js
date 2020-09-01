import {takeEvery} from "redux-saga/effects";
import * as t from "./actionTypes";
// Sagas
import createAccountSaga from "./sagas/createAccountSaga";
import loginAccountSaga from "./sagas/loginAccountSaga";
import logoutAccountSaga from "./sagas/logoutAccountSaga";
/******************************* Watchers *************************************/

export default function* accountRootSaga() {
  yield takeEvery(t.CREATE_ACCOUNT, createAccountSaga);
  yield takeEvery(t.LOGIN_ACCOUNT, loginAccountSaga);
  yield takeEvery(t.LOGOUT_ACCOUNT, logoutAccountSaga);
}
