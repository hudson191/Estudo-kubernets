# Kubernets

## O que é?

"Kubernetes é uma plataforma de código aberto, portável e extensiva para o gerenciamento de cargas de trabalho e serviços distribuídos em contêineres, que facilita tanto a configuração declarativa quanto a automação. Ele possui um ecossistema grande, e de rápido crescimento"<Br/>
Fonte: https://kubernetes.io/pt-br/docs/concepts/overview/what-is-kubernetes/

## Da onde veio?

O nome Kubernetes tem origem no Grego, significando timoneiro ou piloto. K8s é a abreviação derivada pela troca das oito letras "ubernete" por "8", se tornado K"8"s.

## Pontos importantes

- Kubernets é disponibilizado através de um conjunto de APIs que normalmente acessamos usando a CLI: kubectl
- Tudo é baseado em estado. Você configura o estado de cada objeto
- Kubernets Master
  - Kube-apiserver
  - Kube-controller-manager
  - Kube-scheduler
- Outros Nodes (Não Master)
  - Kubelet
  - kubeproxy

**Cluster**: Conjunto de máquinas (Nodes)<Br/>
Cada máquina possui uma quantidade de vCPU e Memória

**Pods**: Menor objeto do kubernets que contém o container provisionado<Br/>

## **kubernets local com kind**

Kind é uma ferramenta para executar cluster Kubernets local usando nós de contêiner  do Docker.<Br/>
Projetado principalmente para testar o próprio kubernets, mas também pode usado para desenvolvimento local.<Br/>
URL: https://kind.sigs.k8s.io/

### Como instalar

Primeiramente precisamos instalar o kubectl, para isso basta seguir o passo a passo descrido em https://kubernetes.io/docs/tasks/tools/<Br/>
Após a instalado o kubectl basta seguir o passo a passo para instalar o kind em https://kind.sigs.k8s.io/docs/user/quick-start#installation


### Criando cluster utilizando o kind

Para criar um cluster basta utilizar o comando abaixo<Br/>
`# kind create cluster`

Podemos notar que utilizamos este processo temos apenas um node, para criar um cluster multi node precisamos criar um arquivo yaml informando os nodes que desejamos, para melhor entendimento recomento a leitura da documentação https://kind.sigs.k8s.io/docs/user/configuration/

Após criado o arquivo yaml  de configuração, podemos criar o cluster utilizado o parametro `--config` informando o local do arquivo.<Br/>
`# kind create cluster --config=k8s/kind.yaml`

Caso deseje escolher o nome do cluster podemos utilziar `--name`.<Br/>
`# kind create cluster --config=k8s/kind.yaml --name=batata`

# Primeiro passos na prática

## Criando aplicação exemplo e imagem

Utilizando nossa aplicação GO que exibe a mensagem "Batatinha Frita" iremos criar um pod, toda a configuração do pod são feitas e arquivos yaml (k8s/pod.yaml) para criar o pode basta executar o comando<br/>
`# kubectl apply -f k8s/pod.yaml`

## ReplicaSet

Porem caso este pod trave por algum motivo este pode seja removido, nosso pod não volta a executar e isso não é algo que queremos que ocorra em produção, para solucionar este problema vamos criar o objeto ReplicaSet (k8s/replicaset.yaml)<br/>
`# kubectl apply -f k8s/replicaset.yaml` 

Desta maneira caso algum pod morra o ReplicaSet cria um novo pod automaticamente, mantendo sempre o mínimo de pods configurados no yaml

## Deployment

Ao gerar uma versão mais nova da aplicação o replicaset não atualiza os pods que estão rodando, para isso acrescentamos um novo objeto Deployment ( Deployment > ReplicaSet > Pod), a configuração de ambos é idêntica alterando apenas o kind o de ReplicaSet para Deployment, como podemos ver no arquivo deployment.yaml.<br/>
`# kubectl apply -f k8s/deployment.yaml`

Caso altere a versão da imagem e aplique as alterações os pods antigos serão removidos e gerado novos pods com as alterações realizadas.

## Rollout e Revisões

Caso a versão mais nova esteja com um bug e seja necessário voltar para versão anterior não precisamos alterar o yaml novamente, podemos fazer um rollout para uma versão anterior, para isso basta rodar o comando `# kubectl rollout history deployment batatinhago`  para listar as revisões e `# kubectl rollout undo deployment batatinhago` para voltar para ultima versão, caso queria voltar para uma revisão especifica basta utilizar `--to-revisio=[NUMERO REVISÃO]`

# Services

Ter nossas aplicações rodando não significa que a aplicação pode ser acessada, para isso utilizamos Services, ele será responsável de escolher para qual pod a requisição será direcionada.

## ClusterIP

Como em toda configuração o Service também é feito por um arquivo yaml (service-clusterip.yaml), nele informamos qual o app em selector que iremos fazer o balanceamento, para isso utilizamos o nome que colocamos ao configurar o pod em matchLabels, para com isso o Service poder filtrar os containers que ira utilizar para balanceamento.

Podemos utilizar o comando `# kubectl get svc` para listar os serviços, com isso podemos notar que foi atribuído um IP interno do kubernets para o serviço, porem a porta não esta publica, com isso precisamos realizar um port-foward para poder acessar o serviço de fora do container.<br/>
`# kubectl port-forward svc/batatinhago-service 8000:80`  

## NodePort

O serviço NodePort é um tipo de serviço que permite acessar o kuster de fora do kubernets, lembrando que precisamos atribuir uma porta maior que **30000** e menor que **32767** e que a porta informada é liberada em todos os nodes

## Environment

Para utilizarmos variáveis  de ambientes, podemos colocar as informação **env** dentro do arquivo de deployment (deployment-env.yaml). 

### ConfigMap

Para não utilizar as variáveis  de ambiente em hard code dentro do deployment podemos utilizar o ConfigMap e criar um arquivo com as variáveis  conforme o k8s/configmap-env.yaml, para o deployment ler o config map basta trocar a opção **value** do **env** por **valueFrom** e informar o nome do config map e qual a chave deseja buscar. 

Outra maneira que podemos trabalhar é buscar todos os parâmetros do ConfigMap, para isso basta trocar a opção env por envFrom (k8s/deployment-env-configmap.yaml).

#### Injetando ConfigMap na aplicação

Existe casos em que precisamos criar arquivos para aplicação, como configuração do NGINX, podemos fazer isso criando um volume e montando ele dentro do container, confirme feito no arquivo k8s/deployment-file-configmap.yaml.

Para verificar se o arquivo foi criado, você pode acessar o pod com o seguinte comando:<br/>
`# kubectl exec -it [NOME_DO_POD] -- bash`

### Sercrets

Secrets é uma forma de trabalhar com os dados "ofuscados", não é algo muito seguro uma vez que as informações são armazenadas em Base64 e facilmente decodificada, para isso criamos um arquivo de configuração k8s/secret.yaml, e configurar dentro do **envFrom** similar ao config map, porem substitui o  **configMapRef** por **secretRef**

## HealthCheck

### LivenessProbes

Muitas vezes nossa aplicação para de funcionar, neste caso precisamos ficar verificando se o sistema está funcionando para isso utilizamos o livenessProbe, ele possui 3 tipos HTTP, command, tcp. Configuramos em k8s/deployment-livenessprobe.yaml.

- failureThreshold: Quantidade de falhas para reiniciar o POD
- successThreshold: Quantidade de teste OK para considerar a aplicação saudavel

https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/

### Readiness

Muitas vezes nossa aplicação apresenta algum erro e precisa parar de receber tráfego, para isso configuramos o **readiness** que valida se a aplicacação esta saudável para receber as requisições, diferente do liveness o readiness apenas bloqueia as requisições enquanto o liveness reinicia o pod. A configuração é similar ao **livenessProbes** como podemos ver no arquivo k8s/deployment-readiness.yaml

### StartupProb

Muito similar ao **readiness** porem realiza a verificção somente quando o pod esta subindo e o **readines** e **liviness** roda somente após o startup falar que a aplicação esta pronta.