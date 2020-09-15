import React, {useEffect, useState} from "react";
import {useDispatch, useSelector} from "react-redux";
import {useHistory, useParams} from "react-router-dom";
//import {receiveFile} from "helpers/apiCalls";
//import styles from "./drive.module.css";
import styles from "../filereceive.module.css";
import main from "styles.module.css";

export function FileReceiveAccept({
  shareId = "",
  account,
  filename = "The Matrix 720p.mp4",
  fileicon = "data:base64",
  nextStage
}) {
  return (<div className={styles.container}>
    <div className={styles.title}>Fairdrive Receive</div>
    <div className={styles.flexer}></div>
    {
      fileicon
        ? (<div>
          <img src={fileicon} className={styles.appicon}></img>
        </div>)
        : ("")
    }
    {
      filename
        ? (<div className={styles.subtitle}>
          <b>{filename}</b>&nbsp; will be saved to your Inbox folder.
        </div>)
        : ("")
    }{" "}
    {
      filename && fileicon
        ? (<div>
          {
            account.locked
              ? (<div tabIndex="2" className={main.button} onClick={nextStage}>
                <div>
                  <div className={main.buttontext}>Save</div>
                </div>
              </div>)
              : (<div tabIndex="2" className={main.button}>
                <div>
                  <div className={main.buttontext}>Unlock to allow</div>
                </div>
              </div>)
          }
        </div>)
        : ("")
    }
    <div className={styles.flexer}></div>
    {
      filename && fileicon
        ? (<div>
          {
            account.locked
              ? ("")
              : (<div className={main.link}>Get me out of here</div>)
          }
        </div>)
        : ("")
    }
  </div>);
}
export default FileReceiveAccept;
