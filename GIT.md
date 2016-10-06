# CES-27 - LAB 2 - Dynamo
## Gerenciando Alterações

No passo [Setup](SETUP.md), nós configuramos e fizemos um clone do repositorio do github na máquina local.  

Existem diversos tutoriais sobre a utilização de Gitpela internet, abaixo estão listados alguns:  

[Try Git](https://try.github.io)  
[Git - The Simple Guide](http://rogerdudler.github.io/git-guide/)  
[Git Official Documentation - Getting Started](https://git-scm.com/book/en/v2/Getting-Started-About-Version-Control)  

É interessante que quem nunca teve contato com Git realize pelo menos um dos tutoriais citado para um melhor entendimento dos passos a seguir.

## Commit - Passo-a-Passo

Para ver o estado das nossas alterações, executamos o seguinte comando:

```shell
wordcount$ git status
```
> On branch master  
> Your branch is up-to-date with 'origin/master'.  
> Changes not staged for commit:  
>   (use "git add <file>..." to update what will be committed)  
>   (use "git checkout -- <file>..." to discard changes in working directory)> 
>
>         modified:   wordcount.go   
>
> no changes added to commit (use "git add" and/or "git commit -a")  

Precisamos selecionar os arquivos para a etapa de Commit.

```shell
wordcount$ git add wordcount.go
```
>  

```shell
$ git status
```
> On branch master  
> Your branch is up-to-date with 'origin/master'.  
> Changes to be committed:  
>   (use "git reset HEAD <file>..." to unstage)>   
>   
>         modified:   wordcount.go  

Agora que as nossas alterações foram selecionadas, vamos realizar o commit:

```shell
$ git commit -m "Mensagem do Commit"
```

> [master a68a570] Imprimindo tamanho do input  
> 1 file changed, 5 insertions(+), 1 deletion(-)

Agora as nossas alterações foram adicionadas ao nosso repositório local:

```shell
$ git log -n 1
```
> commit a68a570562a41bd26b485dcd80ec2592b8e4c4a9  
> Author: Paulo Araujo <phaguiardm@gmail.com>  
> Date:   Wed Aug 24 21:04:56 2016 -0300>   
>  
>     Mensagem do Commit  


Para enviar as suas alterações para o seu repositório remoto (Fork no GitHub), seguir instruções em: [Entrega](ENTREGA.md)

## Atualizando Repositório Local

Caso o repositório *origin* tenha sido atualizado, você pode fazer um *pull* das novas alterações e um *rebase* com as alterações locais utilizando os seguintes comandos:

```shell
$ git fetch origin
$ git rebase origin/master
```
