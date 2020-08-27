import React, { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { useHistory, useParams } from "react-router-dom";

// Sub-pages
import FolderView from './pages/FolderView';

// Ids
const folderViewId = 'folderViewId';

function getAccount(state) {
    return state.account
}

export function DriveRoot() {

    const params = useParams()
    const id = params.id
    const account = useSelector(state => getAccount(state))

    useEffect(() => {
        console.log('account:', account)
    }, [account.privateKey])

    const [stage, setStage] = useState(folderViewId)

    const dispatch = useDispatch()
    const history = useHistory()

    // Router
    switch (stage) {
        case folderViewId:
            return (
                <FolderView
                    id={id}
                    nextStage={() => setStage()}
                    exitStage={() => setStage()}>
                </FolderView>
            );
        default:
            return <h1>Oops...</h1>;
    }
}

export default DriveRoot;
