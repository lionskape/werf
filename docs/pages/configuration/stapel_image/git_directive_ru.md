---
title: Добавление исходного кода из git-репозиториев
sidebar: documentation
permalink: ru/documentation/configuration/stapel_image/git_directive.html
ref: documentation_configuration_stapel_image_git_directive
lang: ru
summary: |
  <a class="google-drawings" href="https://docs.google.com/drawings/d/e/2PACX-1vRUYmRNmeuP14OcChoeGzX_4soCdXx7ZPgNqm5ePcz9L_ItMUqyolRoJyPL7baMNoY7P6M0B08eMtsb/pub?w=2031&amp;h=144" data-featherlight="image">
      <img src="https://docs.google.com/drawings/d/e/2PACX-1vRUYmRNmeuP14OcChoeGzX_4soCdXx7ZPgNqm5ePcz9L_ItMUqyolRoJyPL7baMNoY7P6M0B08eMtsb/pub?w=1016&amp;h=72">
  </a>

  <div class="tabs">
    <a href="javascript:void(0)" class="tabs__btn active" onclick="openTab(event, 'tabs__btn', 'tabs__content', 'local')">Local</a>
    <a href="javascript:void(0)" class="tabs__btn" onclick="openTab(event, 'tabs__btn', 'tabs__content', 'remote')">Remote</a>
  </div>

  <div id="local" class="tabs__content active">
  <div class="language-yaml highlighter-rouge"><div class="highlight"><pre class="highlight"><code><span class="na">git</span><span class="pi">:</span>
  <span class="pi">-</span> <span class="na">add</span><span class="pi">:</span> <span class="s">&lt;absolute path in git repository&gt;</span>
    <span class="na">to</span><span class="pi">:</span> <span class="s">&lt;absolute path inside image&gt;</span>
    <span class="na">owner</span><span class="pi">:</span> <span class="s">&lt;owner&gt;</span>
    <span class="na">group</span><span class="pi">:</span> <span class="s">&lt;group&gt;</span>
    <span class="na">includePaths</span><span class="pi">:</span>
    <span class="pi">-</span> <span class="s">&lt;path or glob relative to path in add&gt;</span>
    <span class="na">excludePaths</span><span class="pi">:</span>
    <span class="pi">-</span> <span class="s">&lt;path or glob relative to path in add&gt;</span>
    <span class="na">stageDependencies</span><span class="pi">:</span>
      <span class="na">install</span><span class="pi">:</span>
      <span class="pi">-</span> <span class="s">&lt;path or glob relative to path in add&gt;</span>
      <span class="na">beforeSetup</span><span class="pi">:</span>
      <span class="pi">-</span> <span class="s">&lt;path or glob relative to path in add&gt;</span>
      <span class="na">setup</span><span class="pi">:</span>
      <span class="pi">-</span> <span class="s">&lt;path or glob relative to path in add&gt;</span></code></pre>
  </div></div>     
  </div>

  <div id="remote" class="tabs__content">
  <div class="language-yaml highlighter-rouge"><div class="highlight"><pre class="highlight"><code><span class="na">git</span><span class="pi">:</span>
  <span class="pi">-</span> <span class="na">url</span><span class="pi">:</span> <span class="s">&lt;git repo url&gt;</span>
    <span class="na">branch</span><span class="pi">:</span> <span class="s">&lt;branch name&gt;</span>
    <span class="na">commit</span><span class="pi">:</span> <span class="s">&lt;commit&gt;</span>
    <span class="na">tag</span><span class="pi">:</span> <span class="s">&lt;tag&gt;</span>
    <span class="na">add</span><span class="pi">:</span> <span class="s">&lt;absolute path in git repository&gt;</span>
    <span class="na">to</span><span class="pi">:</span> <span class="s">&lt;absolute path inside image&gt;</span>
    <span class="na">owner</span><span class="pi">:</span> <span class="s">&lt;owner&gt;</span>
    <span class="na">group</span><span class="pi">:</span> <span class="s">&lt;group&gt;</span>
    <span class="na">includePaths</span><span class="pi">:</span>
    <span class="pi">-</span> <span class="s">&lt;path or glob relative to path in add&gt;</span>
    <span class="na">excludePaths</span><span class="pi">:</span>
    <span class="pi">-</span> <span class="s">&lt;path or glob relative to path in add&gt;</span>
    <span class="na">stageDependencies</span><span class="pi">:</span>
      <span class="na">install</span><span class="pi">:</span>
      <span class="pi">-</span> <span class="s">&lt;path or glob relative to path in add&gt;</span>
      <span class="na">beforeSetup</span><span class="pi">:</span>
      <span class="pi">-</span> <span class="s">&lt;path or glob relative to path in add&gt;</span>
      <span class="na">setup</span><span class="pi">:</span>
      <span class="pi">-</span> <span class="s">&lt;path or glob relative to path in add&gt;</span>
  </code></pre>
  </div></div>
  </div>
---

## Что такое git-маппинг?

***Git-маппинг*** определяет, какой файл или папка из git-репозитория должны быть добавлены в конкретное место образа. Git-репозиторий может быть как локальным репозиторием, в котором находится файл конфигурации сборки (`werf.yaml`), так и удаленным (внешним) репозиторием (в этом случае указывается адрес репозитория и версия кода — ветка, тэг или конкретный коммит).

Werf добавляет файлы из git-репозитория в образ копируя их с помощью [git archive](https://git-scm.com/docs/git-archive) (при первоначальном добавлении файлов) либо накладывая git patch. При повторных сборках и появлении изменений в git-репозитории, Werf добавляет patch к собранному ранее образу — чтобы в конечном образе отразить необходимые изменения файлов и папок. Более подробно, механизм переноса файлов в образ и накладывания патчей рассматривается в соответствующей секции [далее...](#more-details-gitarchive-gitcache-gitlatestpatch)

Конфигурация _git-маппинга_ поддерживает фильтры, что позволяет используя необходимое количество _git_маппингов_ сформировать практически любую файловую структуру в образе. Также, вы можете указать группу и владельца конечных файлов в образе, что освобождает от необходимости делать это (`chown`) отдельной командой.

В Werf реализована поддержка сабмодулей git (git submodules), и если Werf определяет, что какая-то часть git-маппинга является сабмодулем, то принимаются соответствующие меры, чтобы обрабатывать изменения в сабмодулях корректно.

> Все git-сабмодули проекта связаны с конкретным коммитом, поэтому все разработчики работающие с репозиторием использующим сабмодуль, получают одинаковое содержимое. Werf не инициализирует, не обновляет сабмодули, а использует соответствующие связанные коммиты.

Пример добавления файлов из папки `/src` локального git-репозитория в папку `/app` собираемого образа, и добавления кода PhantomJS из удаленного репозитория в папку `/src/phantomjs` собираемого образа:

```yaml
git:
- add: /src
  to: /app
- url: https://github.com/ariya/phantomjs
  add: /
  to: /src/phantomjs
```

## Зачем использовать git-маппинг?

Основная идея использования git-маппинга — добавление истории к сборочному процессу.

### Наложение патчей вместо копирования

Большинство коммитов в репозитории реального приложения относятся к обновлению кода самого приложения. В этом случае, если компиляция приложения не требуется, то для получения нового образа достаточно применить исправления к файлам в предыдущем образе.

### Удаленные репозитории

Сборка конечного образа может зависеть от файлов в других репозитория. Werf позволяет добавлять файлы из удаленных репозиториев, а также отслеживать их изменение.

## Синтаксис

Для добавления кода из локального git-репозитория используется следующий синтаксис:

- `add` — (не обязательный параметр) путь к директории или файлу, содержимое которого (которой) нужно добавить в образ. Указывается абсолютный путь *относительно корня* репозитория, — т.е. он должен начинаться с `/`. По умолчанию копируется все содержимое репозитория, т.е. отсутствие параметра `add` равносильно указанию `add: /`;
- `to` — путь внутри образа, куда будет скопировано соответствующее содержимое;
- `owner` — имя или id пользователя-владельца файлов в образе;
- `group` — имя или id группы-владельца файлов в образе;
- `excludePaths` — список исключений (маска) при рекурсивном копировании файлов и папок. Указывается относительно пути, указанного в `add`;
- `includePaths` — список масок файлов и папок для рекурсивного копирования. Указывается относительно пути, указанного в `add`;
- `stageDependencies` — список масок файлов и папок для указания зависимости пересборки стадии от их изменений. Позволяет указать, при изменении каких файлов и папок необходимо принудительно пересобирать конкретную пользовательскую стадию. Более подробно рассматривается [здесь]({{ site.baseurl }}/ru/documentation/configuration/stapel_image/assembly_instructions.html).

При использовании удаленных репозиториев дополнительно используются следующие параметры:
- `url` — адрес удаленного репозитория;
- `branch`, `tag`, `commit` — имя ветки, тэга или коммита соответственно. По умолчанию — ветка master.

## Использование git-маппинга

### Копирование директорий

Параметр `add` определяет источник, путь в git-репозитории, откуда файлы рекурсивно копируются в образ и помещаются по адресу, указанному в параметре `to`. Если параметр не определен, то по умолчанию используется значение `/`, т.е. копируется весь репозиторий.
Пример простейшей конфигурации, добавляющей содержимое всего локального git-репозитория в образ в папку `/app`.

```yaml
git:
- add: /
  to: /app
```

Если в репозитории была следующая структура файлов и папок:

![git repository files tree]({{ site.baseurl }}/images/build/git_mapping_01.png)

То в образе будет следующая структура файлов и папок:

![image files tree]({{ site.baseurl }}/images/build/git_mapping_02.png)

Вы также можете указывать несколько _git-маппингов_. Пример:

```yaml
git:
- add: /src
  to: /app/src
- add: /assets
  to: /static
```

Если в репозитории была следующая структура файлов и папок:

![git repository files tree]({{ site.baseurl }}/images/build/git_mapping_03.png)

То в образе будет следующая структура файлов и папок:

![image files tree]({{ site.baseurl }}/images/build/git_mapping_04.png)

Следует отметить, что конфигурация git-маппинга не похожа например на копирование типа `cp -r / src / app`. Параметр `add` указывает *содержимое* каталога, которое будет рекурсивно копироваться из репозитория. Поэтому, если папка `/assets` со всем содержимым из репозитория должна быть скопирована в папку `/app/assets` образа, то имя *assets* вы должны указать два раза. Либо, как вариант, вы можете использовать [фильтр](#using-filters), — например параметр `includePaths`.

Примеры обоих вариантов, которые вы можете использовать для достижения одинакового результата:
```yaml
git:
- add: /assets
  to: /app/assets
```

либо

```yaml
git:
- add: /
  to: /app
  includePaths: assets
```

> В Werf нет какого-либо ограничения или соглашения на счет использования `/` в конце, как например в rsync. Т.о. `add: /src` и `add: /src/` — одно и тоже.

### Копирование файла

В случае с копированием файла (когда вы указываете в параметре `add` конкретный файл) действует тот же принцип — вы указываете в параметре `add` содержимое какого файла нужно скопировать, и в параметре `to` — название файла в образе, который будет содержать это содержимое (т.е. также два раза). Это дает вам возможнось изменять имя файла при добавлении его из git-репозитория в образ.

```yaml
git:
- add: /config/prod.yaml
  to: /app/conf/production.yaml
```

### Изменение владельца

При добавлении файла из git-репозитория вы можете указать имя и/или группу владельца файлов в образе. Добавляемым файлам и папкам в образе после копирования будут установлены соответствующие права. Пользователь и группа могут быть указаны как именем так и чиловым id (userid, groupid).

Пример использования:

```yaml
git:
- add: /src/index.php
  to: /app/index.php
  owner: www-data
```

Если указан только параметр `owner`, как в приведенном примере, то группой-владельцем устанавливается основная группа указанного пользователя в системе.

В результате, в папку `/app` образа будет добавлен файл `index.php` и ему будут установлены следующие права:

![index.php owned by www-data user and group]({{ site.baseurl }}/images/build/git_mapping_05.png)

Если значения параметра `owner` или `group` не числовые id, а текстовые (т.е. названия соответствнно пользователя и группы), то соответствующие пользователь и группа должны существовать в системе. Их нужно добавить заранее при необходимости, иначе при сборке возникнет ошибка.

```yaml
git:
- add: /src/index.php
  to: /app/index.php
  owner: wwwdata
```

### Использование фильтров

Парамеры фильтров — `includePaths` и `excludePaths` используются при составлении списка файлов для добавления. Эти параметры содержат набор путей или масок, применяемых соответственно для включения и исключения списка файлов и папок при добавлении в образ.

Фильтр `excludePaths` работает следующим образом: каждая маска списка применяется к каждому файлу, найденному по пути `add`. Если файл удовлетворяет хотябы одной маске, — файл исключается из списка файлов на добавление. Если файл не удовлетворяет ни одной маске — файл добавляется в образ.

Фильтр `includePaths` работает наоборот — если файл удовлетворяет хотябы одной маске, — файл добавляется в образ.

Конфигурация _Git-маппинга_ может содержать оба Фильтра. В этом случае файл добавляется в образ если его путь удовлетворяет хотябы одной маске `includePaths` и не удовлетворяет ни одной маске `excludePaths`.

Пример:

```yaml
git:
- add: /src
  to: /app
  includePaths:
  - '**/*.php'
  - '**/*.js'
  excludePaths:
  - '**/*-dev.*'
  - '**/*-test.*'
```

В приведенном примере добавляются `.php` и `.js` файлы из папки  `/src` исключая файлы с суффиксом `-dev.` или `-test.` в имени файла.

При определении соответствия файла маске, применяется следующий алгоритм:
 - определяется абсолютный путь к очередному файлу в репозитории;
 - путь сравнивается с масками, определенными в `includePaths` и `excludePaths`, либо с конкретным указанным путем:
   - путь в параметре `add` объединяется с маской или указанным путем из параметров `includePaths` и `excludePaths`;
   - оба варианта проверяются с учетом правил применения глобальных шаблонов: если файл удовлетворяет маске — он включается (в случае `includePaths`), либо исключается (в случае `excludePaths`).
 - путь сравнивается с масками, определенными в `includePaths` и `excludePaths`, либо с конкретным указанным путем с учетом дополнительных условий:
   - путь в параметре `add` объединяется с маской или указанным путем из параметров `includePaths` и `excludePaths` и объединяется с суффиксом `**/*` к шаблону;
   - оба варианта проверяются с учетом правил применения глобальных шаблонов: если файл удовлетворяет маске — он включается (в случае `includePaths`), либо исключается (в случае `excludePaths`).

> Последний шаг в алгоритме, с добалвнием  суффикса`**/*` сделан для удобства — вам достаточно указать название папки в параметрах *git-маппинга*, чтобы все ее содержимое удовлетворяло шаблону параметра.

Маска может содержать следующие шаблоны:

- `*` — удовлетворяет любому файлу. Шаблон включает `.` и исключает `/`.
- `**` — удовлетворяет директории со всем ее содержимым, рекурсивно.
- `?` — удовлетворяет любому однму символу в имени файла (аналогично regexp-шаблону `.{1}`)
- `[set]` — удовлетворяет любому символу из указанного набора символов. Аналогично использованию в regexp-шаблонах, включая указание диапазонов типа `[^a-z]`.
- `\` — экранирует следующий символ

Маска, которая начинается с шаблона `*` или `**`, должна быть взята в одинарныеили двойные кавычки в `werf.yaml`:
 - `"*.rb"` — двойные кавычки
- `'**/*'` — одинарные кавычки

Примеры фильтров:

```yaml
add: /src
to: /app
includePaths:
# удовлетворяет всем php файлам, расположенным конкретно в папке /src
- '*.php'

# удовлетворяет всем phph файлам рекурсивно, начиная с папки /src
# (также удовлетворяет файлам *.php, т.к. '.' включается шаблон **)
- '**/*.php'

# удовлетворяет всем файлам в папке /src/module1 рекурсивно
- module1
```

Фильтр `includePaths` может применяться для копирования одного файла без изменения имени. Пример:
```yaml
git:
- add: /src
  to: /app
  includePaths: index.php
```

### Наложение путей копирования

Если вы определяете несколько *git-маппингов*, вы дожлны учитывать, что при наложении путей в образе в парамерре `to` вы можете столкнуться с невозможностью добавления файлов. Пример:

```yaml
git:
- add: /src
  to: /app
- add: /assets
  to: /app/assets
```

Чтобы избежать ошибок сборки, Werf определяет возможные наложения касающиеся фильтров `includePaths` и `excludePaths`, и если такое наложение присутствует, то Werf пытается разрешить самые простые конфликты, неявно добавляя соответсвущий параметр `excludePaths` в git-маппинг. Однако, такое поведение может все-таки привести к ножиданным результатам, поэтому лучше всего избегать наложения путей приопределении git-маппингов.

В примере выше, Werf в итоге неявно добавит параметр  `excludePaths` и итоговая конфигурация будет равнозначна следующей:

```yaml
git:
- add: /src
  to: /app
  excludePaths:  # Werf добавил этот фильтр, чтобы исключить конфлакт наложения результирующих путей
  - assets       # между /src/assets и /assets
- add: /assets
  to: /app/assets
```

## Работа с удаленными репозиториями

Werf может использовать  удаленные (внешние) репозитории в качестве источника файлов. For this purpose, the _git mapping_ configuration contains an `url` parameter where you should specify the repository address. Werf supports `https` and `git+ssh` protocols.

### https

The syntax for https protocol is:

{% raw %}
```yaml
git:
- url: https://[USERNAME[:PASSWORD]@]repo_host/repo_path[.git/]
```
{% endraw %}

`https` access may require login and password.

For example, login and password from GitLab CI variables:

{% raw %}
```yaml
git:
- url: https://{{ env "CI_REGISTRY_USER" }}:{{ env "CI_JOB_TOKEN" }}@registry.gitlab.company.name/common/helper-utils.git
```
{% endraw %}

In this example, the [env](http://masterminds.github.io/sprig/os.html) method from the sprig library is used to access the environment variables.

### git, ssh

Werf supports access to the repository via the git protocol. Access via this protocol is typically protected using ssh tools: this feature is used by GitHub, Bitbucket, GitLab, Gogs, Gitolite, etc. Most often the repository address looks as follows:

```yaml
git:
- url: git@gitlab.company.name:project_group/project.git
```

To successfully work with remote repositories via ssh, you should understand how werf searches for access keys.


#### Working with ssh keys

Keys for ssh connects are provided by ssh-agent. The ssh-agent is a daemon that operates via file socket, the path to which is stored in the environment variable `SSH_AUTH_SOCK`. Werf mounts this file socket to all _assembly containers_ and sets the environment variable `SSH_AUTH_SOCK`, i.e., connection to remote git repositories is established with the use of keys that are registered in the running ssh-agent.

The ssh-agent is determined as follows:

- If werf is started with `--ssh-key` flags (there may be multiple flags):
  - A temporary ssh-agent runs with defined keys, and it is used for all git operations with remote repositories.
  - The already running ssh-agent is ignored in this case.
- No `--ssh-key` flags specified and ssh-agent is running:
  - `SSH_AUTH_SOCK` environment variable is used, and the keys added to this agent is used for git operations.
- No `--ssh-key` flags specified and ssh-agent is not running:
  - If `~/.ssh/id_rsa` file exists, then werf will run the temporary ssh-agent with the  key from `~/.ssh/id_rsa` file.
- If none of the previous options is applicable, then the ssh-agent is not started, and no keys for git operation are available. Build images with remote _git mappings_ ends with an error.

## More details: gitArchive, gitCache, gitLatestPatch

Let us review adding files to the resulting image in more detail. As stated earlier, the docker image contains multiple layers. To understand what layers werf create, let's consider the building actions based on three sample commits: `1`, `2` and `3`:

- Build of commit No. 1. All files are added to a single layer based on the configuration of the _git mappings_. This is done with the help of the git archive. This is the layer of the _gitArchive_ stage.
- Build of commit No. 2. Another layer is added where the files are changed by applying a patch. This is the layer of the _gitLatestPatch_ stage.
- Build of commit No. 3. Files have already added, so werf apply patches in the _gitLatestPatch_ stage layer.

Build sequence for these commits may be represented as follows:

| | gitArchive | --- | gitLatestPatch |
|---|:---:|:---:|:---:|
| Commit No. 1 is made, build at 10:00 |  files as in commit No. 1 | --- | - |
| Commit No. 2 is made, build at 10:05 |  files as in commit No. 1 | --- | files as in commit No. 2 |
| Commit No. 3 is made, build at 10:15 |  files as in commit No. 1 | --- | files as in commit No. 3 |

A space between the layers in this table is not accidental. After a while, the number of commits grows, and the patch between commit No. 1 and the current commit may become quite large, which will further increase the size of the last layer and the total _stages_ size. To prevent the growth of the last layer werf provides another intermediary stage — _gitCache_.
How does werf work with these three stages? Now we are going to need more commits to illustrate this, let it be `1`, `2`, `3`, `4`, `5`, `6` and `7`.

- Build of commit No. 1. As before, files are added to a single layer based on the configuration of the _git mappings_. This is done with the help of the git archive. This is the layer of the _gitArchive_ stage.
- Build of commit No. 2. The size of the patch between `1` and `2` does not exceed 1 MiB, so only the layer of the _gitLatestPatch_ stage is modified by applying the patch between `1` and `2`.
- Build of commit No. 3. The size of the patch between `1` and `3` does not exceed 1 MiB, so only the layer of the _gitLatestPatch_ stage is modified by applying the patch between `1` and `3`.
- Build of commit No. 4. The size of the patch between `1` and `4` exceeds 1 MiB. Now _gitCache_ stage layer is added by applying the patch between `1` and `4`.
- Build of commit No. 5. The size of the patch between `4` and `5` does not exceed 1 MiB, so only the layer of the _gitLatestPatch_ stage is modified by applying the patch between `4` and `5`.

This means that as commits are added starting from the moment the first build is done, big patches are gradually accumulated into the layer for the _gitCache_ stage, and only patches with moderate size are applied in the layer for the last _gitLatestPatch_ stage. This algorithm reduces the size of _stages_.

| | gitArchive | gitCache | gitLatestPatch |
|---|:---:|:---:|:---:|
| Commit No. 1 is made, build at 12:00 |  1 |  - | - |
| Commit No. 2 is made, build at 12:19 |  1 |  - | 2 |
| Commit No. 3 is made, build at 12:25 |  1 |  - | 3 |
| Commit No. 4 is made, build at 12:45 |  1 | *4 | - |
| Commit No. 5 is made, build at 12:57 |  1 |  4 | 5 |

\* — the size of the patch for commit `4` exceeded 1 MiB, so this patch is applied in the layer for the _gitCache_ stage.

### Rebuild of gitArchive stage

For various reasons, you may want to reset the _gitArchive_ stage, for example, to decrease the size of _stages_ and the image.

To illustrate the unnecessary growth of image size assume the rare case of 2GiB file in git repository. First build tranfers this file in the layer of the _gitArchive_ stage. Then some optimization occured and file is recompiled and it's size is decreased to 1.6GiB. The build of this new commit applies patch in the layer of the _gitCache_ stage. The image size become 3.6GiB of which 2GiB is a cached old version of the big file. Rebuilding from _gitArchive_ stage can reduce image size to 1.6GiB. This situation is quite rare but gives a good explanation of correlation between the layers of the _git stages_.

You can reset the _gitArchive_ stage specifying the **[werf reset]** or **[reset werf]** string in the commit message. Let us assume that, in the previous example commit `6` contains **[werf reset]** in its message, and then the builds would look as follows:

| | gitArchive | gitCache | gitLatestPatch |
|---|:---:|:---:|:---:|
| Commit No. 1 is made, build at 12:00 |  1 |  - | - |
| Commit No. 2 is made, build at 12:19 |  1 |  - | 2 |
| Commit No. 3 is made, build at 12:25 |  1 |  - | 3 |
| Commit No. 4 is made, build at 12:45 |  1 |  4 | - |
| Commit No. 5 is made, build at 12:57 |  1 |  4 | 5 |
| Commit No. 6 is made, build at 13:22 |  *6 |  - | - |

\* — commit `6` contains the **[werf reset]** string in its message, so the _gitArchive_ stage is rebuilt.

### _git stages_ and rebasing

Each _git stage_ stores service labels with commits SHA from which this _stage_ was built.
These commits are used for creating patches on the next _git stage_ (in a nutshell, `git diff COMMIT_FROM_PREVIOUS_GIT_STAGE LATEST_COMMIT` for each described _git mapping_).
So, if any saved commit is not in a git repository (e.g., after rebasing) then werf rebuilds that stage with latest commits at the next build.
