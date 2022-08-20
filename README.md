# ProofOfWork gateway to simple tcp service.

<!-- vscode-markdown-toc -->
* 1. [Problem statement](#Problemstatement)
* 2. [How to run](#Howtorun)
	* 2.1. [Requirements](#Requirements)
	* 2.2. [Building base image](#Buildingbaseimage)
	* 2.3. [Linter check](#Lintercheck)
	* 2.4. [Passing tests](#Passingtests)
	* 2.5. [Build images](#Buildimages)
	* 2.6. [Monitoring environment](#Monitoringenvironment)
	* 2.7. [Start server](#Startserver)
	* 2.8. [Start client](#Startclient)
	* 2.9. [Start load-test](#Startload-test)
* 3. [Solution](#Solution)
* 4. [Server environment variables](#Serverenvironmentvariables)
* 5. [Client command lines ??](#Clientcommandlines)

<!-- vscode-markdown-toc-config
	numbering=true
	autoSave=true
	/vscode-markdown-toc-config -->
<!-- /vscode-markdown-toc -->


##  1. <a name='Problemstatement'></a>Problem statement

Test task for Server Engineer

Design and implement “Word of Wisdom” tcp server.  
 • TCP server should be protected from DDOS attacks with the Prof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.  
 • The choice of the POW algorithm should be explained.  
 • After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.  
 • Docker file should be provided both for the server and for the client that solves the POW challenge

##  2. <a name='Howtorun'></a>How to run

###  2.1. <a name='Requirements'></a>Requirements

Possibly can work everywhere on Linux/MacOS X, but tested on this config:

1. OS: Mac OS X 12.5.1
3. Docker Desktop 4.4.2+
4. GNU Make 3.81
5. curl 7.79.1

###  2.2. <a name='Buildingbaseimage'></a>Building base image

```bash
make base-image
```

###  2.3. <a name='Lintercheck'></a>Linter check

```bash
make lint
```

###  2.4. <a name='Passingtests'></a>Passing tests

```bash
make test
```

###  2.5. <a name='Buildimages'></a>Build images

```bash
make docker-build
```

###  2.6. <a name='Monitoringenvironment'></a>Monitoring environment


```bash
make start-monitoring
```

###  2.7. <a name='Startserver'></a>Start server

```bash
make start-server
```

###  2.8. <a name='Startclient'></a>Start client

```bash
make start-client
```

###  2.9. <a name='Startload-test'></a>Start load-test

##  3. <a name='Solution'></a>Solution


##  4. <a name='Serverenvironmentvariables'></a>Server environment variables

##  5. <a name='Clientcommandlines'></a>Client command lines ??


