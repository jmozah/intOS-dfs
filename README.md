# intOS-dfs
The decentralised file system (dfs) is a file system built for the internetOS (intOS).
This works as a think layer over Swarm (https://ethswarm.org/).

dfs can be used as follows
1) As a standalone, decentralised, personal data drive over the internet.
2) In conjunction with intOS-compute layer to work as a data provider for
   large scale parallel data processing engine over the internet. 

### User
![The first screen invites the user to create or restore an account](https://i.imgur.com/VsxTAVW.png)

The first step in dfs is to create a user. Every user is associated with a 12 
word mnemonic based hd wallet. This wallet is passwod protected and stored in 
the datadir. whenever a user created a pod for himself, a new key pair is created 
using this mnemonic. A user can use this mnemonic and import their account in any 
device and instantly see all their pods.

### What is a pod?
A pod is a personal drive created by a user in intOS-dfs. It is used to store files and related metadata in a decentralised fashion. A pod is always under the control of the user who created it. A user can create store any number of files or directories in a pod. 
The user can share files in his pod with any other user just like in other centralised drives like dropbox. Not only users, a pod can be used by decentralised applications (DApp's) to store data related to that user.

The basic storage unit in dfs is a pod. A user can create multiple pods and use it to organise their data. for ex: Personal-Pod, Applications-Pod etc.

### How to run dfs?
- git clone https://github.com/jmozah/intOS-dfs.git
- cd intOS-dfs
- make binary
- yarn 
- yarn build
- ./dist/dfs start (starts server and serves the front end too)

### How to build dfs?
- yarn build (will build the frontend and copy it over to the go project)
- make binary

### HTTP APIs
##### user related APIs
- POST -F 'user=\<username\>' -F 'password=\<password\>' -F 'mnemonic=<12 words from bip39 list>' http://localhost:9090/v0/user/signup
- POST -F 'user=\<username\>' -F 'password=\<password\>' http://localhost:9090/v0/user/signup
- POST -F 'user=\<username\>' -F 'password=\<password\>' http://localhost:9090/v0/user/login 
- POST http://localhost:9090/v0/user/logout
- POST http://localhost:9090/v0/user/avatar
- POST -F 'first_name=\<firstName\>' -F 'middle_name=\<middleName\>' -F 'last_name=\<lastName\>' -F 'surname=\<surName\>' http://localhost:9090/v0/user/name
- POST -F 'phone=\<phone\>' -F 'mobile=\<mobile\>' -F 'address_line_1=\<address1\>' -F 'address_line_2=\<address2\>' http://localhost:9090/v0/user/contact
- DELETE -F 'password=\<password\>' http://localhost:9090/v0/user/delete
- GET  -F 'user=\<username\>' http://localhost:9090/v0/user/present
- GET  -F 'user=\<username\>' http://localhost:9090/v0/user/isloggedin
- GET  http://localhost:9090/v0/user/stat
- GET  http://localhost:9090/v0/user/avatar
- GET  http://localhost:9090/v0/user/name
- GET  http://localhost:9090/v0/user/contact
- GET  http://localhost:9090/v0/user/share/inbox
- GET  http://localhost:9090/v0/user/share/outbox



##### pod related APIs   
- POST -F 'password=\<password\>' -F 'pod=\<podname\>'  http://localhost:9090/v0/pod/new
- POST -F 'password=\<password\>' -F 'pod=\<podname\>'  http://localhost:9090/v0/pod/open
- POST http://localhost:9090/v0/pod/sync
- POST http://localhost:9090/v0/pod/close
- DELETE http://localhost:9090/v0/pod/delete
- GET http://localhost:9090/v0/pod/ls
- GET -F 'user=\<username\>' -F 'pod=\<podname\>'  http://localhost:9090/v0/pod/stat


##### dir related APIs   
- POST -F 'dir=\<dir_with_path\>'  http://localhost:9090/v0/dir/mkdir
- DELETE -F 'dir=\<dir_with_path\>'  http://localhost:9090/v0/dir/rmdir
- GET  -F 'dir=\<dir_with_path\>'  http://localhost:9090/v0/dir/ls
- GET  -F 'dir=\<dir_with_path\>'  http://localhost:9090/v0/dir/stat

##### file related APIs   
- POST -F -H "intOS-dfs-Compression: snappy/gzip" 'pod_dir=\<dir_with_path\>' -F 'block_size=\<in_Mb\>' -F 'files=@\<filename1\>' -F 'files=@\<filename2\>' http://localhost:9090/v0/file/upload  (compression header optional)
- POST -F 'file=\<file_path\>'  http://localhost:9090/v0/file/download
- POST -F 'file=\<file_path\>' -F 'to=\<destination_user\>' http://localhost:9090/v0/file/share
- POST http://localhost:9090/v0/file/share -d &InboxEntry(json struct)
- GET  -F 'file=\<file_path\>'  http://localhost:9090/v0/file/stat
- DELETE -F 'file=\<file_path\>'  http://localhost:9090/v0/file/delete


### Comands in dfs
**dfs >>>** \<command\> where, \<command\> is listed below
##### user related commands
- user \<new\> (user-name) - create a new user and login as that user
- user \<del\> (user-name) - deletes a already created user
- user \<login\> (user-name) - login as a given user
- user \<logout\> (user-name) - logout as user
- user \<ls\> - lists all the user present in this instance
##### pod related commands
- pod \<new\> (pod-name) - create a new pod and login to that pod
- pod \<del\> (pod-name) - deletes a already created pod
- pod \<login\> (pod-name) - login to a already created pod
- pod \<stat\> (pod-name) - display meta information about a pod
- pod \<sync\> (pod-name) - sync the contents of a logged in pod from Swarm
- pod \<logout\>  - logout of a logged in pod
- pod \<ls\> - lists all the pods created for this account
##### directory & file related commands
- cd \<directory name\>
- ls 
- copyToLocal \<source file in pod, destination directory in local fs\>
- copyFromLocal \<source file in local fs, destination directory in pod, block size in MB\>
- mkdir \<directory name\>
- rmdir \<directory name\>
- rm \<file name\>
= pwd - show present working directory
- head \<no of lines\>
- cat  - stream the file to stdout
- stat \<file name or directory name\> - shows the information about a file or directory
##### REPL related commands
- help - display this help
- exit - exits from the prompt
