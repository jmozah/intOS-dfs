import React, {useEffect, useRef, useState} from "react";
import styles from "../drive.module.css";
import {Route, NavLink} from "react-router-dom";
import {useDispatch, useSelector} from "react-redux";
import {useHistory} from "react-router-dom";
import sortByProp from "helpers/sortByProp";
import urlPath from "helpers/urlPath";
import NewDialog from "./components/NewDialog";

import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  List,
  ListItem,
  ListItemIcon,
  ListItemText
} from "@material-ui/core";

import {
  mdiFolder,
  mdiFolderEdit,
  mdiSettingsHelper,
  mdiShare,
  mdiTrashCan,
  mdiZipBox
} from "@mdi/js";
import Icon from "@mdi/react";

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
import {createDirectory, deleteDirectory, fileUpload} from "helpers/apiCalls";

export function FolderView({
  nextStage,
  exitStage,
  path,
  contents,
  account,
  refresh,
  setFolderShown
}) {
  const [uploadShown, setUploadShown] = useState(false);
  const [uploadStatus, setUploadStatus] = useState("ready");
  const [uploadProgress, setUploadProgress] = useState(0);

  const [folderToEdit, setFolderToEdit] = useState("");
  const [openFolder, setFolderOpen] = useState(false);
  const [openNew, setNewOpen] = useState(true);

  const [newFolderName, setNewFolderName] = useState("");
  const hiddenFileInput = useRef(null);

  const toSortProp = "name";
  const [toSort, setToSort] = useState(toSortProp);
  const orderProp = "asc";

  const dispatch = useDispatch();
  const history = useHistory();

  function handleFolderClickOpen() {
    setFolderOpen(true);
  }

  function handleFolderClose() {
    setFolderOpen(false);
  }

  function handleNewClickOpen() {
    setNewOpen(true);
  }

  function handleNewClickClose() {
    setNewOpen(false);
  }

  useEffect(() => {
    setFolderShown(true);
  }, []);

  function handleClick(event) {
    hiddenFileInput.current.click();
  }

  function handleChange(event) {
    handleFileUpload(event.target.files);
  }

  function handleFolderNameChange(e) {
    setNewFolderName(e.target.value);
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

  function toggleFolderMenuShown(item) {
    setFolderToEdit(item);
    handleFolderClickOpen();
  }

  function handleLocation(item) {
    console.log(item);
    let writePath = "";
    if (path == "root") {
      writePath = "";
    } else {
      writePath = path + "&";
    }
    history.push("/drive/" + writePath + item);
  }

  function handleGotoAccount() {
    history.push("/account");
  }

  async function handleDeleteFolder(folderName) {
    console.log(folderName);
    await deleteDirectory(folderName);
    refresh(path);
    handleFolderClose();
  }

  const selectedIcon = icon => {
    switch (icon) {
      case "inode/directory":
        return <Icon path={mdiFolder}></Icon>;
        break;
      case "application/zip":
        return <Icon path={mdiZipBox}></Icon>;
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

  const Entries = contents => {
    switch (contents.entries.length) {
      case 0:
        return <div className={styles.folderLoading}>Nothing here yet.</div>;
        break;
      default:
        return contents.entries.sort(sortByProp(toSort, orderProp)).map(item => (<div key={item.name} className={styles.rowItem}>
          <div onClick={() => handleLocation(item.name)}>
            {selectedIcon(item.content_type)}
          </div>
          <div onClick={() => handleLocation(item.name)} className={styles.folderText}>
            {item.name}
          </div>
          <div>
            <Icon path={mdiSettingsHelper} onClick={() => toggleFolderMenuShown(item.name)} className={styles.custom} rotate={90} size="36px"></Icon>
          </div>
        </div>));
        break;
    }
  };

  const FolderDialogFragment = () => {
    return (<Dialog open={openFolder} onClose={handleFolderClose} fullWidth="fullWidth">
      <DialogTitle>
        <span className={styles.folderMenuTitle}>{folderToEdit}</span>
      </DialogTitle>
      <DialogContent>
        <List>
          <ListItem button="button" divider="divider" role="listitem">
            <ListItemIcon>
              <Icon path={mdiShare} size="24px"></Icon>
            </ListItemIcon>
            <ListItemText primary="Share"/>
          </ListItem>
          <ListItem button="button" divider="divider" role="listitem">
            <ListItemIcon>
              <Icon path={mdiFolderEdit} size="24px"></Icon>
            </ListItemIcon>
            <ListItemText primary="Rename"/>
          </ListItem>
          <ListItem onClick={() => handleDeleteFolder(folderToEdit)} button="button" divider="divider" role="listitem">
            <ListItemIcon>
              <Icon path={mdiTrashCan} size="24px"></Icon>
            </ListItemIcon>
            <ListItemText primary="Delete"/>
          </ListItem>
        </List>
      </DialogContent>
      <DialogActions>
        <Button onClick={handleFolderClose}>Close</Button>
      </DialogActions>
    </Dialog>);
  };

  return (<div className={styles.container}>
    <div className={styles.topbar}>
      <div className={styles.topmenu}>
        <div className={styles.user}>
          <div onClick={() => handleGotoAccount()} className={styles.username}>
            {account.username}
          </div>
          <div className={styles.balance}>
            {account.balance}
            &nbsp; BZZ
          </div>
        </div>
        <div className={styles.addButton} onClick={() => handleNewClickOpen()}>
          <AddCircleOutline fontSize="large"></AddCircleOutline>
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
        contents
          ? (Entries(contents))
          : (<div className={styles.folderLoading}>
            <CircularProgress></CircularProgress>
          </div>)
      }
    </div>
    {FolderDialogFragment()}
    <NewDialog open={openNew} onClose={() => handleNewClickClose()} path={path} refresh={refresh}></NewDialog>
  </div>);
}

export default FolderView;
