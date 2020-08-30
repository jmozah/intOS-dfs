import React, {useRef, useState} from "react";
import styles from "../drive.module.css";
import {Route, NavLink} from "react-router-dom";
import {
  AddCircleOutline,
  Cloud,
  Folder,
  HighlightOff,
  LibraryMusic,
  Subject
} from "@material-ui/icons/";
import defaultAvatar from "images/defaultAvatar.png";

export function FolderView({
  nextStage,
  exitStage,
  id,
  contents,
  account,
  handleFileUpload
}) {
  const directories = contents.directories;

  const [uploadShown, setUploadShown] = useState(true);

  const hiddenFileInput = React.useRef(null);

  const handleClick = event => {
    hiddenFileInput.current.click();
  };

  const handleChange = event => {
    handleFileUpload(event.target.files);
  };

  const toggleUploadShown = () => {
    setUploadShown(!uploadShown);
  };

  const selectedIcon = icon => {
    switch (icon) {
      case "folder":
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
  return (<div className={styles.container}>
    <div className={styles.topbar}>
      <div className={styles.topmenu}>
        <div className={styles.user}>
          <div className={styles.username}>{account.username}</div>
          <div className={styles.balance}>
            {account.balance}
            BZZ
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
          ? (<div className={styles.uploadSpace} onClick={handleClick}>
            <Cloud fontSize="large"></Cloud>
            <div>Upload some files</div>
            <input multiple="multiple" type="file" ref={hiddenFileInput} onChange={handleChange} style={{
                display: "none"
              }}/>
          </div>)
          : (<div>
            <div className={styles.title}>
              {
                id === "root"
                  ? "Your Fairdrive"
                  : id
              }
            </div>
            <div className={styles.status}>~3211MB</div>
          </div>)
      }
    </div>
    <div className={styles.innercontainer}>
      {
        directories.map(item => (<div className={styles.rowItem}>
          <div>{selectedIcon(item.icon)}</div>
          <div className={styles.folderText}>{item}</div>
        </div>))
      }
    </div>
  </div>);
}

export default FolderView;
