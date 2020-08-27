import React, { useState } from "react";
import styles from "styles.module.css";
import createAccount from "../account.module.css";
import QRCode from "react-qr-code";

const MyIdentity = ({
    nextStage,
    exitStage,
    shortCode,
    resolveShort,
    avatarStage,
    avatar,
    username,
    setUsername,
    publicKey = "",
    address = "",
    QR = true
}) => (
        <div className={createAccount.formcontainer}>
            <div className={createAccount.closeButton} onClick={exitStage}>
                <div className={styles.exitgrayicon} />
            </div>
            {/* <div onClick={shortCode}>CREATE SHORT</div>
            <div onClick={resolveShort}>FIND SHORT</div> */}
            <div className={createAccount.placeholder}>
                <img className={createAccount.avatarImage} src={avatar} alt="avatar" />
            </div>
            <div className={createAccount.subtitle}>{username}</div>

            <div>
                <QRCode size="128" value={address}></QRCode>
            </div>
            <div className={createAccount.flexer} />
            <div tabIndex="2" className={styles.button} onClick={nextStage}>
                <div>
                    <div className={styles.buttontext}>edit account</div>
                </div>
            </div>
        </div>
    );

export default MyIdentity;
