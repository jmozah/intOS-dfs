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
  const [folderToEdit, setFolderToEdit] = useState("");
  const [openFolder, setFolderOpen] = useState(false);
  const [openNew, setNewOpen] = useState(false);

  const [newFolderName, setNewFolderName] = useState("");

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

  function handleFolderNameChange(e) {
    setNewFolderName(e.target.value);
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

  function stripLastPath(path) {
    return path.split("/").pop();
  }

  function pathToArray(path) {
    console.log(urlPath(path).split("/"));
    return urlPath(path).split("/");
  }

  async function handleDeleteFolder(folderName) {
    console.log(folderName);
    await deleteDirectory(folderName);
    refresh(path);
    handleFolderClose();
  }

  function handleGotoAccount() {
    history.push("/account");
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

  function getPathForItem(item, path) {
    let patharray = pathToArray(path);
    let index = patharray.indexOf(item);
    patharray = patharray.slice(0, index + 1);
    let urlPath = patharray.join("&");
    return urlPath;
  }

  function breadCrumb(path) {
    let patharray = pathToArray(path);
    console.log(patharray);
    patharray = patharray.slice(0, patharray.length - 1);
    return patharray.map(item => (<div className={styles.breadcrumbspace}>
      <NavLink className={styles.breadcrumbitem} to={getPathForItem(item, path)}>
        {item + "/"}
      </NavLink>
    </div>));
  }

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

      <div>
        <div className={styles.title}>
          {
            path === "root"
              ? "My Fairdrive"
              : stripLastPath(urlPath(path))
          }
        </div>
        {
          path != "root"
            ? (<div className={styles.breadcrumb}>
              <div>Back to: &nbsp;</div>
              <NavLink className={styles.breadcrumbitem} to="/drive/root">
                My Fairdrive/
              </NavLink>
              {breadCrumb(path)}
            </div>)
            : ("")
        }
      </div>
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
