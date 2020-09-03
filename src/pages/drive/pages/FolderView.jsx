import React, {useRef, useState} from "react";
import styles from "../drive.module.css";
import {Route, NavLink} from "react-router-dom";
import {useDispatch, useSelector} from "react-redux";
import {useHistory} from "react-router-dom";

import {
  AddCircleOutline,
  Cloud,
  Folder,
  HighlightOff,
  LibraryMusic,
  Subject,
  FileCopySharp
} from "@material-ui/icons/";

import {CircularProgress, LinearProgress} from "@material-ui/core";
import defaultAvatar from "images/defaultAvatar.png";
import {fileUpload} from "helpers/apiCalls";

export function FolderView({
  nextStage,
  exitStage,
  path,
  contents,
  account,
  refresh
}) {
  const [uploadShown, setUploadShown] = useState(false);
  const [uploadStatus, setUploadStatus] = useState("ready");
  const [uploadProgress, setUploadProgress] = useState(0);

  const hiddenFileInput = React.useRef(null);
  const dispatch = useDispatch();
  const history = useHistory();

  const handleClick = event => {
    hiddenFileInput.current.click();
  };

  function handleChange(event) {
    handleFileUpload(event.target.files);
  }

  async function handleFileUpload(files) {
    setUploadStatus("uploading");
    await fileUpload(files, path, function (progress, total) {
      setUploadProgress(Math.round((progress / total) * 100));
      if (progress == total) {
        setUploadStatus("swarmupload");
      }
    }).then(() => {
      toggleUploadShown();
      setUploadStatus("ready");
      refresh(path);
    }).catch(() => {
      setUploadStatus("error");
    });

    //dispatch({type: "GET_DRIVE"});
  }

  function toggleUploadShown() {
    setUploadShown(!uploadShown);
  }

  function handleLocation(item) {
    console.log(item);
    history.push("/drive/" + item);
  }

  const selectedIcon = icon => {
    switch (icon) {
      case "Dir":
        return <Folder></Folder>;
        break;
      case "txt":
        return <Subject></Subject>;
        break;
      case "mp3":
        return <LibraryMusic></LibraryMusic>;
      default:
        return <img className={styles.fileIcon} src={defaultAvatar}></img>;
        break;
    }
  };

  const UploadStage = status => {
    switch (status) {
      case "ready":
        return (<div className={styles.uploadSpace} onClick={handleClick}>
          <Cloud fontSize="large"></Cloud>
          <div>Upload some files</div>
          <input multiple="multiple" type="file" ref={hiddenFileInput} onChange={handleChange} style={{
              display: "none"
            }}/>
        </div>);
        break;
      case "uploading":
        return (<div className={styles.uploadSpace} onClick={handleClick}>
          <div>Uploading...</div>
          <LinearProgress className={styles.progress} variant="determinate" value={uploadProgress}/>
        </div>);
        break;
      case "swarmupload":
        return (<div className={styles.uploadSpace} onClick={handleClick}>
          <div>Storing on Swarm...</div>
          <LinearProgress className={styles.progress}/>
        </div>);
        break;
      case "error":
        return (<div className={styles.uploadSpace} onClick={handleClick}>
          <div>Error</div>
        </div>);
        break;
    }
  };
  return (<div className={styles.container}>
    <div className={styles.topbar}>
      <div className={styles.topmenu}>
        <div className={styles.user}>
          <div className={styles.username}>{account.username}</div>
          <div className={styles.balance}>
            {account.balance}
            &nbsp; BZZ
          </div>
        </div>
        <div className={styles.addButton} onClick={() => toggleUploadShown()}>
          {
            uploadShown
              ? (<HighlightOff fontSize="large"></HighlightOff>)
              : (<AddCircleOutline fontSize="large"></AddCircleOutline>)
          }
        </div>
      </div>
      <div className={styles.flexer}></div>
      {
        uploadShown
          ? (<div>{UploadStage(uploadStatus)}</div>)
          : (<div>
            <div className={styles.title}>
              {
                path === "root"
                  ? "My Fairdrive"
                  : path
              }
            </div>
            <div className={styles.status}>~3211MB</div>
          </div>)
      }
    </div>
    <div className={styles.innercontainer}>
      {
        contents.Entries
          ? (contents.Entries.map(item => (<div className={styles.rowItem} onClick={() => handleLocation(item.name)}>
            <div>{selectedIcon(item.type)}</div>
            <div className={styles.folderText}>{item.name}</div>
          </div>)))
          : (<div className={styles.folderLoading}>
            <CircularProgress></CircularProgress>
          </div>)
      }
    </div>
  </div>);
}

export default FolderView;
