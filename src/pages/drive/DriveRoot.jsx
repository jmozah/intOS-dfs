import React, {useEffect, useState} from "react";
import {useDispatch, useSelector} from "react-redux";
import {useHistory, useParams} from "react-router-dom";
import {getDirectory} from "helpers/apiCalls";

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
  const path = params.path;
  const account = useSelector(state => getAccount(state));
  const driveContent = useSelector(state => getContents(state)) || {
    directories: [],
    files: []
  };

  const dispatch = useDispatch();

  const [contents, setContents] = useState(driveContent);

  async function getDirectoryContent(path) {
    const content = await getDirectory(path);
    setContents(content);
    //return content;
  }

  useEffect(() => {
    //dispatch({type: "GET_DRIVE"});

    getDirectoryContent(path).catch(e => console.log(e));

    //console.log("account:", account);
  }, [path]);

  const [stage, setStage] = useState(folderViewId);

  const history = useHistory();

  // Router
  switch (stage) {
    case folderViewId:
      return (<FolderView path={path} refresh={getDirectoryContent} account={account} contents={contents} nextStage={() => setStage()} exitStage={() => setStage()}></FolderView>);
    default:
      return <h1>Oops...</h1>;
  }
}

export default DriveRoot;
