import React from "react"
import main from "styles.module.css"
import accountstyles from "../account-create.module.css"

export function AccountCreateIntro({ createStage, restoreStage, exitStage }) {

    return (
        <div className={accountstyles.container}>
            <div className={accountstyles.title}>Welcome to Fairdrive</div>
            <div className={accountstyles.subtitle}>
                In the next steps you will be creating a Fairdrive Wallet.
            </div>

            <div tabIndex="2" className={main.button} onClick={createStage}>
                <div>
                    <div className={main.buttontext}>create account</div>
                </div>
            </div>

            <div tabIndex="2" className={main.button}>
                <div>
                    <div className={main.buttontext} onClick={restoreStage}>restore account</div>
                </div>
            </div>
        </div>
    )
}

export default AccountCreateIntro