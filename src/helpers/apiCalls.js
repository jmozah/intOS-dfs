import axios from "axios";
import qs from "querystring";

const axi = axios.create({baseURL: "http://localhost:9090/v0/", timeout: 120000});

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
      data: qs.stringify({password: "1234", pod: "Fairdrive"}),
      config: config,
      withCredentials: true
    });

    return response;
  } catch (error) {
    throw error;
  }
}

export async function isLoggedIn(username) {
  return true;
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
