import React, {useEffect, useRef, useState} from "react";
import styles from "../../drive.module.css";
import rootStyles from "styles.module.css";
import urlPath from "helpers/urlPath";

import {CircularProgress, LinearProgress} from "@material-ui/core";

import {Dialog} from "@material-ui/core";
import {
  mdiFolder,
  mdiFolderPlus,
  mdiFolderEdit,
  mdiSettingsHelper,
  mdiShare,
  mdiTrashCan,
  mdiUpload,
  mdiZipBox
} from "@mdi/js";
import Icon from "@mdi/react";
import {createDirectory, deleteDirectory, deleteFile, fileUpload} from "helpers/apiCalls";

export default function FolderDialog({open, path, refresh, onClose, item}) {
  console.log("from newdialog: ", open, path);

  const homeId = "homeId";
  const shareId = "shareId";
  const renameId = "renameId";
  const deleteId = "deleteId";
  const errorId = "errorId";

  const [folderContentState, setFolderContentState] = useState(homeId);
  const [newFolderName, setNewFolderName] = useState();
  const [shareName, setShareName] = useState();

  const [newName, setNewName] = useState();

  useEffect(() => {
    setNewFolderName(item.name);
    console.log(item);
  }, [item]);

  function handleFolderClose() {
    setFolderContentState(homeId);
    onClose();
  }

  function handleFolderNameChange(e) {
    setNewFolderName(e.target.value);
  }

  function handleShareChange(e) {
    setShareName(e.target.value);
  }

  async function handleRename() {
    refresh(path);
    handleFolderClose();
  }

  async function handleDelete() {
    // what is it?
    let writePath = "";
    if (path == "root") {
      writePath = "/";
    } else {
      writePath = "/" + urlPath(path) + "/";
    }

    if (item.content_type == "inode/directory") {
      await deleteDirectory(writePath + item.name);
    } else {
      await deleteFile(writePath + item.name);
    }
    //await deleteFile
    refresh(path);
    handleFolderClose();
  }

  async function handleShare() {
    refresh(path);
    handleFolderClose();
  }

  const FolderDialogContent = () => {
    switch (folderContentState) {
      case homeId:
        return (<div className={styles.foldermenu}>
          <div className={styles.menutitle}>
            <div>Action: {item.name}</div>
            <div className={styles.close} onClick={handleFolderClose}>
              <div className={styles.closeicon}/>
            </div>
          </div>
          <div className={styles.menuitem} onClick={() => setFolderContentState(shareId)}>
            <Icon path={mdiFolderPlus} className={styles.blue} size="24px"></Icon>
            Share
          </div>
          <div className={styles.menuitem} onClick={() => setFolderContentState(renameId)}>
            <Icon path={mdiUpload} className={styles.blue} size="24px"></Icon>
            Rename
          </div>
          <div className={styles.menuitem} onClick={() => setFolderContentState(deleteId)}>
            <Icon path={mdiFolderPlus} className={styles.blue} size="24px"></Icon>
            Delete
          </div>
        </div>);
        break;
      case shareId:
        return (<div className={styles.foldermenu}>
          <div className={styles.menutitle}>
            <div>Share</div>
            <div className={styles.close} onClick={handleFolderClose}>
              <div className={styles.closeicon}/>
            </div>
          </div>

          <input className={styles.nameInput} placeholder="Friend's name" type="text" onChange={e => handleShareChange(e)}></input>
          <div className={styles.flexer}>
            <div onClick={handleShare} className={styles.buttonPlace}>
              <div className={rootStyles.buttontext}>> share</div>
            </div>
          </div>
        </div>);
        break;
      case renameId:
        return (<div className={styles.foldermenu}>
          <div className={styles.menutitle}>
            <div>Rename folder</div>
            <div className={styles.close} onClick={handleFolderClose}>
              <div className={styles.closeicon}/>
            </div>
          </div>

          <input className={styles.nameInput} value={newFolderName} type="text" onChange={e => handleFolderNameChange(e)}></input>
          <div className={styles.flexer}></div>

          <div onClick={handleRename} className={styles.buttonPlace}>
            <div className={rootStyles.buttontext}>> rename folder</div>
          </div>
        </div>);
      case deleteId:
        return (<div className={styles.foldermenu}>
          <div className={styles.menutitle}>
            <div>Delete {item.name}?</div>
            <div className={styles.close} onClick={handleFolderClose}>
              <div className={styles.closeicon}/>
            </div>
          </div>

          <div onClick={handleDelete} className={styles.buttonPlace}>
            <div className={rootStyles.buttontext}>> remove permanently</div>
          </div>
        </div>);

      case errorId:
        return (<div className={styles.foldermenu}>
          <div className={styles.close} onClick={handleFolderClose}>
            <div className={styles.closeicon}/>
          </div>

          <div className={styles.uploadSpace}>
            <div className={styles.statusText}>Error...</div>
          </div>
        </div>);
        break;
      default:
        break;
    }
  };

  return (<Dialog open={open} fullWidth="fullWidth">
    <div className={styles.dialogContainer}>{FolderDialogContent()}</div>
  </Dialog>);
}