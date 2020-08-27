
import { takeEvery } from "redux-saga/effects";
import * as t from "./actionTypes";
// Sagas
import unlockSystemSaga from "./sagas/unlockSystemSaga"


/******************************* Watchers *************************************/

export default function* driveRootSaga() {
    //yield systemSaga()
    //yield takeEvery(t.SET_DRIVE, setDriveSaga);
}
