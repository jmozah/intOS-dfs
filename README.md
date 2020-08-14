# intOS-dfs
The decentralised file system (dfs) is a file system built for the internetOS (intOS).
This works as a think layer over Swarm (https://ethswarm.org/).

dfs can be used as follows
1) As a standalone, decentralised, personal data drive over the internet.
2) In conjunction with intOS-compute layer to work as a data provider for
   large scale parallel data processing engine over the internet. 

### User

every user is associated with a 24 word mnemonic based hd wallet. This wallet is 
passwod protected and stored in the datadir. whenever a user created a pod for
himself, a new key pair is created using this mnemonic. A user can use this
mnemonic and import his account in any device and instantly see all his pods.

### What is a pod?

A pod is a personal drive created by a user in intOS-dfs. It is used to store files and 
related metadata in a decentralised fashion. A pod is always under the control of the user
who created it. A user can create store any number of files or directories in a pod. 
The user can share files in his pod with any other user just like in other centralised 
drives like dropbox. Not only users, a pod can be used by decentralised applications (DApp's) 
to store data related to that user.

The basic storage unit in dfs is a pod. A user can create multiple pods and use it to organise 
his data. for ex: Personal-Pod, Applications-Pod etc.

### How to run dfs?

For now dfs is a command line program. Later there will be a UI that will be hosted in Swarm.
- git clone https://github.com/jmozah/intOS-dfs.git
- cd intOS-dfs
- make build
- cd bin
- ./dfs prompt  (this will start dfs in REPL mode with a "dfs >>>" prompt)

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
