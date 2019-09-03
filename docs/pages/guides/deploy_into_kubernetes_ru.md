---
title: Деплой в Kubernetes
sidebar: documentation
permalink: ru/documentation/guides/deploy_into_kubernetes.html
author: Timofey Kirillov <timofey.kirillov@flant.com>
ref: documentation_guides_deploy_into_kubernetes
lang: ru
---

## Обзор задачи

Будет рассмотрено, как деплоить приложение в Kubernetes с помощью Werf.

Werf использует (с некоторыми изменениями и дополнениями) [Helm](helm.sh) для деплоя приложений в Kubernetes, и в статье мы создадим простое web-приложение, соберем все необходимые для него образы, создадим helm-template'ы и запустим приложение в вашем кластере Kubernetes.

## Требования

 * Работающий кластер Kubernetes. Для выполнения примера вы можете использовать как обычный Kubernetes кластер, так и Minikube. Если вы решили использовать Minikube, прочитайте [статью о настройке Minikube]({{ site.baseurl }}/documentation/reference/development_and_debug/setup_minikube.html), чтобы запустить Minikube и Docker-registry.
 * Работающий Docker-registry.
   * Доступ от хостов Kubernetes с правами на push образов в registry.
   * Доступ от хостов Kubernetes с правами на pull образов в registry.
 * Установленные [зависимости Werf]({{ site.baseurl }}/documentation/guides/installation.html#install-dependencies).
 * Установленный [Multiwerf](https://github.com/flant/multiwerf).
 * Установленный `kubectl` и сконфигурированный для доступа в кластер Kubernetes (<https://kubernetes.io/docs/tasks/tools/install-kubectl/>).

**Внимание!** Далее, в качестве адреса репозитория мы будем использовать значение — `:minikube` . Если вы используете ваш существующий кластер Kubernetes и отдельный экземпляр docker-registry, указывайте его вместо аргумента `:minikube`.


### Выбор версии Werf

Перед началом работы с Werf, нужно выбрать версию Werf, которую вы будете использовать. Для выбора актуальной версии Werf в канале beta, релиза 1.0, выполните в вашей shell-сессии:

```shell
source <(multiwerf use 1.0 beta)
```

## Архитектура приложения

Пример представляет собой простейшее web-приложение, для запуска которого нам нужен только web-сервер.

Архитектура подобных приложений в Kubernetes выглядит как правило слещующе:

     .----------------------.
     | backend (Deployment) |
     '----------------------'
                |
                |
      .--------------------.
      | frontend (Ingress) |
      '--------------------'

Здесь `backend` — web-сервер с приложением, `frontend` — прокси-сервер, который выступает точкой входа и перенаправления внешнего трафика в приложение.

## Файлы приложения

Werf ожидает что все файлы, необходимые для сборки и развертывания приложения, находятся в папке с исходным кодом самого приложения (строго говоря, исходного кода прилоения может и не быть) — папке приложения (папке проекта).

Создадим пустую папку, на машине, где будет происходить сборка:

```shell
mkdir myapp
cd myapp
```

## Подготовка образа

Нам нужно подготовить образ приложения с web-сервером внутри. Для этого, создайте фаил `werf.yaml` в папке приложения следующего содержания:

```yaml
project: myapp
configVersion: 1
---

image: ~
from: python:alpine
ansible:
  install:
  - file:
      path: /app
      state: directory
      mode: 0755
  - name: Prepare main page
    copy:
      content:
        <!DOCTYPE html>
        <html>
          <body>
            <h2>Congratulations!</h2>
            <img src="https://flant.com/images/logo_en.png" style="max-height:100%;" height="76">
          </body>
        </html>
      dest: /app/index.html
```

Наше web-приложение состоит из единственной статической HTML-страницы, которая создается прямо на этапе сборки образа, в инструкциях сборки. Содержимое этой страницы будет отдавать Python HTTP-сервер.

Соберите образ приложения и загрузите его в Docker-registry:

```shell
werf build-and-publish --stages-storage :local --tag-custom myapp --images-repo :minikube
```

Название собранного образа приложения состоит из адреса Docker-registry (`REPO`) и тэга (`TAG`). При указании `:minikube` в качестве адреса Docker-registry, Werf использует в качестве адреса Docker-registry адрес `werf-registry.kube-system.svc.cluster.local:5000/myapp`. Так как мы указали в качестве тэга образа `myapp`, Werf загрузит в Docker-registry образ с именем `werf-registry.kube-system.svc.cluster.local:5000/myapp:myapp`.

## Подготовка конфигурации деплоя

Werf использует код [Helm](helm.sh) *для применения* конфигурации в Kubernetes. Для *описания* объектов Kubernetes, Werf в том числе использует конфигурационные файлы Helm — шаблоны, файлы параметров, такие как `values.yaml` и т.д. Также, Werf использует расширенные конфигурации, такие как — шифрованные файлы, файлы с секретами (`secret-values.yaml`), собственные helm-шаблоны для подстановки имен образов и переменных при деплое.

### Backend

Создайте фаил конфигурации backend `.helm/templates/010-backend.yaml`, и далее мы рассмотрим его подробнее:

{% raw %}
```yaml
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: {{ .Chart.Name }}-backend
spec:
  replicas: 4
  template:
    metadata:
      labels:
        service: {{ .Chart.Name }}-backend
    spec:
      containers:
      - name: backend
        workingDir: /app
        command: [ "python3", "-m", "http.server", "8080" ]
{{ include "werf_container_image" . | indent 8 }}
        livenessProbe:
          httpGet:
            path: /
            port: 8080
            scheme: HTTP
        readinessProbe:
          httpGet:
            path: /
            port: 8080
            scheme: HTTP
        ports:
        - containerPort: 8080
          name: http
          protocol: TCP
        env:
{{ include "werf_container_env" . | indent 8 }}
---
apiVersion: v1
kind: Service
metadata:
  name: {{ .Chart.Name }}-backend
spec:
  clusterIP: None
  selector:
    service: {{ .Chart.Name }}-backend
  ports:
  - name: http
    port: 8080
    protocol: TCP
```
{% endraw %}

В конфигурации описывается создание Deployment'а `myapp-backend` (конструкция {% raw %}`{{ .Chart.Name }}-backend`{% endraw %} будет преобразована в `myapp-backend`) с несколькими репликами.

Конструкция {% raw %}`{{ include "werf_container_image" . | indent 8 }}`{% endraw %} использует внутренний helm-шаблон Werf, который:
* всегда возвращает поле `image:` объекта Kubernetes с корректным именем образа, учитывая используемую схему тегирования (в примере это — `werf-registry.kube-system.svc.cluster.local:5000/myapp:latest`)
* дополнительно может возвращать другие поля объекта Kubernetes, такие как `imagePullPolicy`, на основании заложенной логики и некоторых внешних условий.

Go-шаблон `werf_container_image` предоставляет удобный способ указания имени образа в объекте Kubernetes **из описанной конфигурации**. Как использовать шаблон в случае если в конфигурации описано несколько образов, читай подробнее [в соответствующей статье]({{ site.baseurl }}/ru/documentation/reference/deploy_process/deploy_into_kubernetes.html#werf_container_image).

Конструкция {% raw %}`{{ include "werf_container_env" . | indent 8 }}`{% endraw %} использует другой внутренний helm-шаблон Werf, который *может* возвращать содержимое секции `env:` соответствующего контейнера объекта Kubernetes. It is needed for kubernetes to shut down and restart deployment pods only when docker image has been changed, [see the reference for details]({{ site.baseurl }}/documentation/reference/deploy_process/deploy_into_kubernetes.html#werf_container_env).

Finally, in this configuration Service `myapp-backend` specified to access Pods of Deployment `myapp-backend`.

### Frontend configuration

To describe `frontend` configuration place file `.helm/templates/090-frontend.yaml` with the following content:

```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: myapp-frontend
  annotations:
    kubernetes.io/ingress.class: "nginx"
    nginx.ingress.kubernetes.io/ssl-redirect: "false"
spec:
  rules:
  - host: myapp.local
    http:
      paths:
      - path: /
        backend:
          serviceName: myapp-backend
          servicePort: 8080
```

This Ingress configuration set up nginx proxy server for host `myapp.local` to our backend web server `myapp-backend`.

## Run deploy

If you use minikube then enable ingress addon before running deploy:

```shell
minikube addons enable ingress
```

Run deploy with werf:

```shell
werf deploy --stages-storage :local --images-repo :minikube --tag-custom myapp --env dev
```

With this command werf will create all kubernetes resources using helm and watch until `myapp-backend` Deployment is ready (when all replicas Pods are up and running).

[Environment]({{ site.baseurl }}/documentation/reference/deploy_process/deploy_into_kubernetes.html#environment) `--env` is a required param needed to generate helm release name and kubernetes namespace.

Helm release with name `myapp-dev` will be created. This name consists of [project name]({{ site.baseurl }}/documentation/configuration/introduction.html#meta-configuration-doc) `myapp` (which you've placed in the `werf.yaml`) and specified environment `dev`. Check docs for details about [helm release name generation]({{ site.baseurl }}/documentation/reference/deploy_process/deploy_into_kubernetes.html#release-name).

Kubernetes namespace `myapp-dev` will also be used. This name also consists of [project name]({{ site.baseurl }}/documentation/configuration/introduction.html#meta-configuration-doc) `myapp` and specified environment `dev`. Check docs for details about [kubernetes namespace generation]({{ site.baseurl }}/documentation/reference/deploy_process/deploy_into_kubernetes.html#kubernetes-namespace).

## Check your application

Now it is time to know the IP address of your kubernetes cluster. If you use minikube get it with (in the most cases the IP address will be `192.168.99.100`):

```shell
minikube ip
```

Make sure that host name `myapp.local` is resolving to this IP address on your machine. For example append this record to the `/etc/hosts` file:

```shell
192.168.99.100 myapp.local
```

Then you can check application by url: `http://myapp.local`.

## Delete application from kubernetes

To completely remove deployed application run this dismiss werf command:

```shell
werf dismiss --env dev --with-namespace
```

## See also

For all werf deploy features such as secrets [take a look at reference]({{ site.baseurl }}/documentation/reference/deploy_process/deploy_into_kubernetes.html).
