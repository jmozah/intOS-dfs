import React, { useState } from "react";
import styles from "../drive.module.css";
import { Route, NavLink } from "react-router-dom";
import { Folder, LibraryMusic, Subject } from '@material-ui/icons/';
import defaultAvatar from "images/defaultAvatar.png"

export function FolderView({
    nextStage,
    exitStage,
    id,
    contents = [
        { title: "Documents", icon: "folder" },
        { title: "Pictures", icon: "folder" },
        { title: "Movies", icon: "folder" },
        { title: "DappConnect", icon: "folder" },
        { title: "Shared", icon: "folder" },
        { title: "Keith Jarret - Koln Concert Extra long title.mp3", icon: "mp3" },
        { title: "HelloWorld.txt", icon: "txt" },
        { title: "MyAvatar.png", icon: "sad" },
        { title: "BBC.Documentary.mp4", icon: "asdasd" },
    ]
}) {

    const selectedIcon = (icon) => {
        switch (icon) {
            case "folder":
                return <Folder></Folder>
                break;
            case "txt":
                return <Subject></Subject>
                break;
            case "mp3":
                return <LibraryMusic></LibraryMusic>
            default:
                return <img className={styles.fileIcon} src={defaultAvatar}></img>
                break;
        }
    }
    console.log("someting")
    return (
        <div className={styles.container}>
            <div className={styles.topbar}>
                <div className={styles.username}>Michelle</div>
                <div className={styles.balance}>102.32 BZZ</div>
                <div className={styles.flexer}></div>
                <div className={styles.title}>{id === "root" ? "Your Fairdrive" : id}</div>
                <div className={styles.status}>~3211MB</div>
            </div>
            <div className={styles.innercontainer}>
                {contents.map(item => (
                    <div className={styles.rowItem}>
                        <div>{selectedIcon(item.icon)}</div>
                        <div className={styles.folderText}>{item.title}</div>
                    </div>
                ))}
            </div>
        </div>
    )
}

export default FolderView;
