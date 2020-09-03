import React, {useState} from "react";
import {useHistory} from "react-router-dom";
import {useDispatch} from "react-redux";
import defaultAvatar from "images/defaultAvatar.png";
import {createAccount, createDirectory, createPod} from "helpers/apiCalls";

// Sub-pages
import AccountCreateIntro from "./pages/AccountCreateIntro";
import MnemonicShow from "./pages/MnemonicShow";
import MnemonicCheck from "./pages/MnemonicCheck";
import ChooseUsername from "./pages/ChooseUsername";
import ChoosePassword from "./pages/ChoosePassword";
import ChooseAvatar from "./pages/ChooseAvatar";
import CreatingAccount from "./pages/CreatingAccount";
import {createNextState} from "@reduxjs/toolkit";

// Ids
const accountCreateIntroId = "accountCreateIntroId";
const mnemonicShowId = "mnemonicShowId";
const mnemonicCheckId = "mnemonicCheckId";
const chooseUsernameId = "chooseUsernameId";
const chooseAvatarId = "chooseAvatarId";
const choosePasswordId = "choosePasswordId";
const creatingAccountId = "creatingAccountId";

export function AccountCreateRoot() {
  const dispatch = useDispatch();

  const [stage, setStage] = useState(accountCreateIntroId);
  const history = useHistory();
  // Mnemonic for debugging
  //   const [mnemonic, setMnemonic] = useState([
  //     "scissors",
  //     "system",
  //     "judge",
  //     "reveal",
  //     "slogan",
  //     "rice",
  //     "option",
  //     "body",
  //     "bronze",
  //     "insane",
  //     "evolve",
  //     "matter"
  //   ]);
  const [mnemonic, setMnemonic] = useState([]);
  const [wallet, setWallet] = useState();
  const [collection, setCollection] = useState();
  const [avatar, setAvatar] = useState(defaultAvatar);
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState();

  const [accountCreateDone, setAccountCreateDone] = useState(false);
  const [item0, setItem0] = useState(false);
  const [item1, setItem1] = useState(false);
  const [item2, setItem2] = useState(false);
  const [item3, setItem3] = useState(false);

  // Create account function
  const createAccountProcess = async () => {
    setStage(creatingAccountId);
    const mnemonicJoined = mnemonic.join(" ");
    const newUser = await createAccount(username, password, mnemonicJoined);

    setItem0(true);
    await createPod(password, "Fairdrive");

    setItem1(true);
    await createDirectory(password, "Documents");
    await createDirectory(password, "Movies");
    await createDirectory(password, "Music");
    await createDirectory(password, "Pictures");

    setItem2(true);
    // store account in Redux
    const userObject = {
      status: "accountSet",
      username: username,
      avatar: avatar,
      address: newUser.reference,
      balance: 0.0
    };

    dispatch({type: "SET_ACCOUNT", data: userObject});
    dispatch({
      type: "SET_SYSTEM",
      data: {
        hasAcount: true
      }
    });

    setItem3(true);
    history.push("/drive/root");
  };

  // Router
  switch (stage) {
    case accountCreateIntroId:
      return (<AccountCreateIntro createStage={() => setStage(mnemonicShowId)} restoreStage={() => setStage()} exitStage={() => history.goBack()}/>);

    case mnemonicShowId:
      return (<MnemonicShow fairdrive={window.fairdrive} nextStage={() => setStage(mnemonicCheckId)} exitStage={() => setStage(accountCreateIntroId)} setMnemonic={setMnemonic} mnemonic={mnemonic} setCollection={setCollection}/>);
    case mnemonicCheckId:
      return (<MnemonicCheck nextStage={() => setStage(chooseUsernameId)} prevStage={() => setStage(mnemonicShowId)} exitStage={() => setStage(accountCreateIntroId)} mnemonic={mnemonic}/>);
    case chooseUsernameId:
      return (<ChooseUsername avatar={avatar} setUsername={setUsername} username={username} nextStage={() => setStage(choosePasswordId)} exitStage={() => setStage(accountCreateIntroId)} avatarStage={() => setStage(chooseAvatarId)}></ChooseUsername>);
    case chooseAvatarId:
      return (<ChooseAvatar avatar={defaultAvatar} exitStage={() => setStage(chooseUsernameId)} setAvatar={setAvatar}></ChooseAvatar>);
    case choosePasswordId:
      return (<ChoosePassword createAccount={createAccountProcess} exitStage={() => setStage(accountCreateIntroId)} setPassword={setPassword} password={password}/>);
    case creatingAccountId:
      return (<CreatingAccount accountCreateDone={accountCreateDone} item0={item0} item1={item1} item2={item2} item3={item3}/>);
    default:
      return <h1>Oops...</h1>;
  }
}

export default AccountCreateRoot;
