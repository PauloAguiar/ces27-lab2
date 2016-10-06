# CES-27 - LAB 2 - Dynamo
## Preparando o Repositório local

Os comandos a seguir requerem uma instalação da ferramenta [Git](https://git-scm.com/). Para verificar a sua instalação basta executar o seguinte comando:

```shell
$ git --version
```
> git version 2.9.2.windows.1

Caso a ferramenta não esteja instalada, seguir instruções do site oficial [Git](https://git-scm.com/)

Além disso, as ferramentas de Go devem estar devidamente instaladas e configuradas, incluindo o GOPATH (ver [Go - Getting Started](https://golang.org/doc/install) em caso de dúvidas). Para verificar a sua instalação, basta digitar no prompt o seguinte comando:

```shell
$ go env
```
>set GOARCH=amd64  
>set GOBIN=  
>set GOEXE=.exe  
>set GOHOSTARCH=amd64  
>set GOHOSTOS=windows  
>set GOOS=windows  
>set GOPATH=C:\gows  
>set GORACE=  
>set GOROOT=C:\tools\Go  
>set GOTOOLDIR=C:\tools\Go\pkg\tool\windows_amd64  
>set GO15VENDOREXPERIMENT=1  
>set CC=gcc  
>set GOGCCFLAGS=-m64 -mthreads -fmessage-length=0  
>set CXX=g++  
>set CGO_ENABLED=1  

Com todas as ferramentas devidamente configuradas, obtenha os arquivos do laboratório usando o seguinte comando:
```shell
$ go get -d github.com/pauloaguiar/ces27-lab2
```

Esse comando vai colocar no seu *workspace* os arquivos do laboratório, no diretório *src/github.com/pauloaguiar/ces27-lab2*.

Você não deve alterar a estrutura do diretório, pois isso quebrará as referências internas do código.