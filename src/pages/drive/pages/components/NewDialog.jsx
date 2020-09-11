import React, {useEffect, useRef, useState} from "react";
import styles from "../../drive.module.css";
import urlPath from "helpers/urlPath";

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
import {createDirectory, deleteDirectory, fileUpload} from "helpers/apiCalls";

export default function NewDialog({open, path, refresh, onClose}) {
  console.log("from newdialog: ", open, path);

  const homeId = "homeId";
  const newFolderId = "newFolderId";

  //const [openNew, setNewOpen] = useState(open);
  const [newDialogContentState, setNewDialogContentState] = useState(homeId);

  const [newFolderName, setNewFolderName] = useState();

  function handleNewClose() {
    setNewDialogContentState(homeId);
    onClose();
  }

  function handleFolderNameChange(e) {
    setNewFolderName(e.target.value);
  }

  async function handleNewFolder() {
    console.log(newFolderName);
    let writePath = "";
    if (path == "root") {
      writePath = "/";
    } else {
      writePath = "/" + urlPath(path);
    }
    await createDirectory(writePath + newFolderName);
    refresh(path);
    handleNewClose();
  }

  const NewDialogContent = () => {
    switch (newDialogContentState) {
      case homeId:
        return (<div>
          <DialogContent>
            <List>
              <ListItem onClick={() => setNewDialogContentState(newFolderId)} button="button" divider="divider" role="listitem">
                <ListItemIcon>
                  <Icon path={mdiShare} size="24px"></Icon>
                </ListItemIcon>
                <ListItemText primary="New Folder"/>
              </ListItem>
              <ListItem button="button" divider="divider" role="listitem">
                <ListItemIcon>
                  <Icon path={mdiFolderEdit} size="24px"></Icon>
                </ListItemIcon>
                <ListItemText primary="Upload Items"/>
              </ListItem>
            </List>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleNewClose}>Close</Button>
          </DialogActions>
        </div>);
        break;
      case newFolderId:
        return (<div>
          <DialogContent className={styles.center}>
            <input className={styles.nameInput} placeholder="Folder name" type="text" onChange={e => handleFolderNameChange(e)}></input>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleNewFolder}>Save</Button>

            <Button onClick={handleNewClose}>Close</Button>
          </DialogActions>
        </div>);
      default:
        break;
    }
  };

  return (<Dialog open={open} fullWidth="fullWidth">
    <div className={styles.dialogContainer}>{NewDialogContent()}</div>
  </Dialog>);
}