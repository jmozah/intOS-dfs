
import { takeEvery } from "redux-saga/effects";
import * as t from "./actionTypes";
// Sagas
import unlockSystemSaga from "./sagas/unlockSystemSaga"


/******************************* Watchers *************************************/

export default function* systemRootSaga() {
    //yield systemSaga()
    yield takeEvery(t.UNLOCK_SYSTEM, unlockSystemSaga);
}
