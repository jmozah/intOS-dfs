import axios from "axios";
import qs from "querystring";
import {Avatar} from "@material-ui/core";

const axi = axios.create({baseURL: "http://fairdrive.org/v0/", timeout: 120000});

export async function logIn(username, password) {
  try {
    const requestBody = {
      user: username,
      password: password
    };

    const config = {
      headers: {
        "Content-Type": "application/x-www-form-urlencoded"
      }
    };

    const response = await axi({method: "POST", url: "user/login", config: config, data: qs.stringify(requestBody), withCredentials: true});

    const openPod = await axi({
      method: "POST",
      url: "pod/open",
      data: qs.stringify({password: password, pod: "Fairdrive"}),
      config: config,
      withCredentials: true
    });

    return response;
  } catch (error) {
    throw error;
  }
}

export async function logOut(username) {
  try {
    const requestBody = {
      user: username
    };

    const config = {
      headers: {
        "Content-Type": "application/x-www-form-urlencoded"
      }
    };

    const response = await axi({method: "POST", url: "user/logout", config: config, data: qs.stringify(requestBody), withCredentials: true});

    return response;
  } catch (error) {
    throw error;
  }
}

export async function isLoggedIn(username) {
  try {
    const requestBody = {
      user: username
    };

    const config = {
      headers: {
        "Content-Type": "application/x-www-form-urlencoded"
      }
    };

    const response = await axi({method: "POST", url: "user/isloggedin", config: config, data: qs.stringify(requestBody), withCredentials: true});

    return response;
  } catch (error) {
    throw error;
  }
}

export async function isUsernamePresent(username) {
  try {
    const requestBody = {
      user: username
    };

    const config = {
      headers: {
        "Content-Type": "application/x-www-form-urlencoded"
      }
    };

    const response = await axi({method: "POST", url: "user/present", config: config, data: qs.stringify(requestBody), withCredentials: true});

    return response;
  } catch (error) {
    throw error;
  }
}

export async function fileUpload(files, directory, onUploadProgress) {
  const config = {
    headers: {
      "Content-Type": "multipart/form-data"
    }
  };

  const formData = new FormData();
  for (const file of files) {
    formData.append("files", file);
  }
  formData.append("pod_dir", "/" + directory);
  formData.append("block_size", "64Mb");

  const uploadFiles = await axi({
    method: "POST",
    url: "file/upload",
    data: formData,
    config: config,
    withCredentials: true,
    onUploadProgress: function (event) {
      onUploadProgress(event.loaded, event.total);
    }
  });

  console.log(uploadFiles);
  return true;
}

export async function getDirectory(directory) {
  try {
    const config = {
      headers: {
        "Content-Type": "application/x-www-form-urlencoded"
      }
    };

    const openPod = await axi({
      method: "POST",
      url: "pod/open",
      data: qs.stringify({password: "1234", pod: "Fairdrive"}),
      config: config,
      withCredentials: true
    });

    let data = "/";

    if (directory == "root") {
      data = qs.stringify({dir: "/"});
    } else {
      data = qs.stringify({
        dir: "/" + directory
      });
    }

    const response = await axi({method: "POST", url: "dir/ls", data: data, config: config, withCredentials: true});

    return response.data;
  } catch (error) {
    throw error;
  }
}

export async function createAccount(username, password, mnemonic) {
  try {
    const requestBody = {
      user: username,
      password: password,
      mnemonic: mnemonic
    };

    const config = {
      headers: {
        "Content-Type": "application/x-www-form-urlencoded"
      }
    };

    const response = await axi({method: "POST", url: "user/signup", config: config, data: qs.stringify(requestBody), withCredentials: true});
    return response.data;
  } catch (e) {
    console.log("error on timeout", e);
  }
}

function dataURLtoFile(dataurl, filename) {
  var arr = dataurl.split(","),
    mime = arr[0].match(/:(.*?);/)[1],
    bstr = atob(arr[1]),
    n = bstr.length,
    u8arr = new Uint8Array(n);

  while (n--) {
    u8arr[n] = bstr.charCodeAt(n);
  }

  return new File([u8arr], filename, {type: mime});
}

export async function storeAvatar(avatar) {
  try {
    //Usage example:
    var file = dataURLtoFile(avatar, "avatar.jpg");

    const formData = new FormData();
    formData.append("avatar", file);

    const config = {
      headers: {
        "Content-Type": "multipart/form-data"
      }
    };

    const response = await axi({method: "POST", url: "user/avatar", config: config, data: formData, withCredentials: true});
    return response.data;
  } catch (e) {
    console.log("error on timeout", e);
  }
}

export async function getAvatar(username) {
  try {
    const config = {
      responseType: "arraybuffer",

      headers: {
        "Content-Type": "application/x-www-form-urlencoded"
      }
    };

    const data = {
      username: username
    };

    const response = await axi({method: "GET", url: "user/avatar", config: config, data: qs.stringify(data), withCredentials: true});
    console.log(response);

    const blob = new Blob([response.data]);

    var reader = new FileReader();
    reader.readAsDataURL(blob);
    reader.onloadend = function () {
      var base64data = reader.result;
      console.log(base64data);
    };

    return response.data;
  } catch (e) {
    console.log("error on timeout", e);
  }
}

export async function createPod(passWord, podName) {
  try {
    const config = {
      headers: {
        "Content-Type": "application/x-www-form-urlencoded"
      }
    };
    const podRequest = {
      password: passWord,
      pod: podName
    };

    const createPod = await axi({method: "POST", url: "pod/new", config: config, data: qs.stringify(podRequest), withCredentials: true});
  } catch (error) {}
}

export async function createDirectory(passWord, directoryName) {
  try {
    const config = {
      headers: {
        "Content-Type": "application/x-www-form-urlencoded"
      }
    };

    const createPictursDirectory = await axi({
      method: "POST",
      url: "dir/mkdir",
      config: config,
      data: qs.stringify({dir: directoryName}),
      withCredentials: true
    });
    return true;
  } catch (error) {}
}