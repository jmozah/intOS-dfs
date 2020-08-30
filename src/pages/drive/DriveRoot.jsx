import React, {useEffect, useState} from "react";
import {useDispatch, useSelector} from "react-redux";
import {useHistory, useParams} from "react-router-dom";

// Sub-pages
import FolderView from "./pages/FolderView";

// Ids
const folderViewId = "folderViewId";

function getAccount(state) {
  return state.account;
}

function getContents(state) {
  return state.drive.fairdrive;
}

export function DriveRoot() {
  const params = useParams();
  const id = params.id;
  const account = useSelector(state => getAccount(state));
  const driveContent = useSelector(state => getContents(state)) || {
    directories: []
  };
  const dispatch = useDispatch();

  useEffect(() => {
    dispatch({type: "GET_DRIVE"});
    //console.log("account:", account);
  }, []);

  const [stage, setStage] = useState(folderViewId);

  const history = useHistory();

  const handleFileUpload = files => {
    console.log(files.length);
    dispatch({type: "UPLOAD_FILES", data: files});
  };

  // Router
  switch (stage) {
    case folderViewId:
      return (<FolderView id={id} account={account} handleFileUpload={handleFileUpload} contents={driveContent} nextStage={() => setStage()} exitStage={() => setStage()}></FolderView>);
    default:
      return <h1>Oops...</h1>;
  }
}

export default DriveRoot;
