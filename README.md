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

### ReplicaSet

Porem caso este pod trave por algum motivo este pode seja removido, nosso pod não volta a executar e isso não é algo que queremos que ocorra em produção, para solucionar este problema vamos criar o objeto ReplicaSet (k8s/replicaset.yaml)<br/>
`# kubectl apply -f k8s/replicaset.yaml` 

Desta maneira caso algum pod morra o ReplicaSet cria um novo pod automaticamente, mantendo sempre o mínimo de pods configurados no yaml

### Deployment

Ao gerar uma versão mais nova da aplicação o replicaset não atualiza os pods que estão rodando, para isso acrescentamos um novo objeto Deployment ( Deployment > ReplicaSet > Pod), a configuração de ambos é idêntica alterando apenas o kind o de ReplicaSet para Deployment, como podemos ver no arquivo deployment.yaml.<br/>
`# kubectl apply -f k8s/deployment.yaml`

Caso altere a versão da imagem e aplique as alterações os pods antigos serão removidos e gerado novos pods com as alterações realizadas.

### Rollout e Revisões

Caso a versão mais nova esteja com um bug e seja necessário voltar para versão anterior não precisamos alterar o yaml novamente, podemos fazer um rollout para uma versão anterior, para isso basta rodar o comando `# kubectl rollout history deployment batatinhago`  para listar as revisões e `# kubectl rollout undo deployment batatinhago` para voltar para ultima versão, caso queria voltar para uma revisão especifica basta utilizar `--to-revisio=[NUMERO REVISÃO]`
