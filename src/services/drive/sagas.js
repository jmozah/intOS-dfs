import {takeEvery} from "redux-saga/effects";
import * as t from "./actionTypes";
// Sagas
import getDriveContentSaga from "./sagas/getDriveContentSaga";

/******************************* Watchers *************************************/

export default function* driveRootSaga() {
  //yield systemSaga()
  yield takeEvery(t.GET_DRIVE, getDriveContentSaga);
}
