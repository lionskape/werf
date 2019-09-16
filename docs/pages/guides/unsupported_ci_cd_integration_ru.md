---
title: Интеграция с другими CI/CD системами
sidebar: documentation
permalink: ru/documentation/guides/unsupported_ci_cd_integration.html
author: Timofey Kirillov <timofey.kirillov@flant.com>
ref: documentation_guides_unsupported_ci_cd_integration
lang: ru
---

В настоящий момент Werf поддерживает только работу с [Gitlab CI]({{ site.baseurl }}/documentation/reference/plugging_into_cicd/gitlab_ci.html). Мы планируем в [скором будущем](https://github.com/flant/werf/issues/1682) обеспечить поддержку top-10 популярных CI систем.

Чтобы использовать Werf с любой CI/CD системой которая пока не поддерживается, необходимо выполнить шаги описанные [здесь]({{ site.baseurl }}/ru/documentation/reference/plugging_into_cicd/overview.html#what-is-ci-env), с помощью собственного скрипта. Запуск такого скрипта нужно производить вместо вызова команды `werf ci-env`, но поведение скрипта должно быть похожим на результат выполнения команды `werf ci-env`. Запуск скрипта также должен осуществляться перед выполнением любых других команд werf, вначале задания CI/CD.

Необходимого результата можно добиться выполнив некоторые действия и определив ряд переменных окружения из [списка]({{ site.baseurl }}/documentation/reference/plugging_into_cicd/overview.html#complete-list-of-ci-env-params-and-customizing).

## Ci-env procedures

### Интеграция с Docker registry

Согласно процедуре [интеграции с Docker registry]({{ site.baseurl }}/documentation/reference/plugging_into_cicd/overview.html#docker-registry-integration), необходимо определить следующие переменные:
 * [`DOCKER_CONFIG`]({{ site.baseurl }}/documentation/reference/plugging_into_cicd/overview.html#docker_config);
 * [`WERF_IMAGES_REPO`]({{ site.baseurl }}/documentation/reference/plugging_into_cicd/overview.html#werf_images_repo).

Также необходимо создать папку для временного файла конфигурации docker в рамках выполнения задания. Пример:

```bash
mkdir .docker
export DOCKER_CONFIG=$(pwd)/.docker
export WERF_IMAGES_REPO=DOCKER_REGISTRY_REPO
```

### Git integration

According to [git integration]({{ site.baseurl }}/documentation/reference/plugging_into_cicd/overview.html#git-integration) procedure, variables to define:
 * [`WERF_TAG_GIT_TAG`]({{ site.baseurl }}/documentation/reference/plugging_into_cicd/overview.html#werf_tag_git_tag);
 * [`WERF_TAG_GIT_BRANCH`]({{ site.baseurl }}/documentation/reference/plugging_into_cicd/overview.html#werf_tag_git_branch).

### CI/CD pipelines integration

According to [CI/CD pipelines integration]({{ site.baseurl }}/documentation/reference/plugging_into_cicd/overview.html#ci-cd-pipelines-integration) procedure, variables to define:
 * [`WERF_ADD_ANNOTATION_GIT_REPOSITORY_URL`]({{ site.baseurl }}/documentation/reference/plugging_into_cicd/overview.html#werf_add_annotation_git_repository_url).

### CI/CD configuration integration

According to [CI/CD configuration integration]({{ site.baseurl }}/documentation/reference/plugging_into_cicd/overview.html#ci-cd-configuration-integration) procedure, variables to define:
 * [`WERF_ENV`]({{ site.baseurl }}/documentation/reference/plugging_into_cicd/overview.html#werf_env).

### Configure modes of operation in CI/CD systems

According to [configure modes of operation in CI/CD systems]({{ site.baseurl }}/documentation/reference/plugging_into_cicd/overview.html#configure-modes-of-operation-in-ci-cd-systems) procedure, variables to define:

Variables to define:
 * [`WERF_GIT_TAG_STRATEGY_LIMIT`]({{ site.baseurl }}/documentation/reference/plugging_into_cicd/overview.html#werf_git_tag_strategy_limit);
 * [`WERF_GIT_TAG_STRATEGY_EXPIRY_DAYS`]({{ site.baseurl }}/documentation/reference/plugging_into_cicd/overview.html#werf_git_tag_strategy_expiry_days);
 * [`WERF_LOG_COLOR_MODE`]({{ site.baseurl }}/documentation/reference/plugging_into_cicd/overview.html#werf_log_color_mode);
 * [`WERF_LOG_PROJECT_DIR`]({{ site.baseurl }}/documentation/reference/plugging_into_cicd/overview.html#werf_log_project_dir);
 * [`WERF_ENABLE_PROCESS_EXTERMINATOR`]({{ site.baseurl }}/documentation/reference/plugging_into_cicd/overview.html#werf_enable_process_exterminator);
 * [`WERF_LOG_TERMINAL_WIDTH`]({{ site.baseurl }}/documentation/reference/plugging_into_cicd/overview.html#werf_log_terminal_width).

## Ci-env script

Copy following script and place into `werf-ci-env.sh` in the root of the project:

```bash
mkdir .docker
export DOCKER_CONFIG=$(pwd)/.docker
export WERF_IMAGES_REPO=DOCKER_REGISTRY_REPO
docker login -u USER -p PASSWORD $WERF_IMAGES_REPO

export WERF_TAG_GIT_TAG=GIT_TAG
export WERF_TAG_GIT_BRANCH=GIT_BRANCH
export WERF_ADD_ANNOTATION_GIT_REPOSITORY_URL="project.werf.io/ci-url=https://cicd.domain.com/project/x"
export WERF_ENV=ENV
export WERF_GIT_TAG_STRATEGY_LIMIT=10
export WERF_GIT_TAG_STRATEGY_EXPIRY_DAYS=30
export WERF_LOG_COLOR_MODE=on
export WERF_LOG_PROJECT_DIR=1
export WERF_ENABLE_PROCESS_EXTERMINATOR=1
export WERF_LOG_TERMINAL_WIDTH=100
```

This script needs to be customized to your CI/CD system: change `WERF_*` environment variables values to the real ones. Consult with the following pages to get an idea and examples of how to retrieve real values for werf variables:
 * [Gitlab CI integration]({{ site.baseurl }}/documentation/reference/plugging_into_cicd/gitlab_ci.html)

Copy following script and place into `werf-ci-env-cleanup.sh`:

```bash
rm -rf .docker
```

`werf-ci-env.sh` should be called in the beginning of everr CI/CD job prior running any werf commands.
`werf-ci-env-cleanup.sh` should be called in the end of every CI/CD job.
