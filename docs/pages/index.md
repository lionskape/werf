---
title: Main page
permalink: /
layout: default
---

<div class="welcome">
    <div class="page__container">
        <div class="welcome__content">
            <h1 class="welcome__title">
                GitOps CLI tool
            </h1>
            <div class="welcome__subtitle">
                 Deliver your application fast & easy.<br/>Open Source. Written in Golang.
            </div>
            <form action="https://www.google.com/search" class="welcome__search" method="get" name="searchform" target="_blank">
                <input name="sitesearch" type="hidden" value="werf.io">
                <input autocomplete="on" class="page__input welcome__search-input" name="q" placeholder="Search the documentation" required="required"  type="text">
                <button type="submit" class="page__icon page__icon_search welcome__search-btn"></button>
            </form>
        </div>
    </div>
</div>

<div class="page__container">
    <div class="intro">
        <div class="intro__image"></div>
        <div class="intro__content">
            <div class="intro__title">
                CLI tool to construct CI/CD pipelines
            </div>
            <div class="intro__text">
                <ul class="intro__list">
                    <li>
                        Werf is a single CLI tool that integrates well known tools:<br/> <code>git</code>, <code>helm</code> and <code>docker</code>.
                    </li>
                    <li>
                        Werf can be embedded into any existing CI/CD system (like Gitlab CI) <br>to implement CI/CD pipelines using provided building blocks:
                        <ul class="intro__list_c2">
                            <li><code>werf build-and-publish</code>;</li>
                            <li><code>werf deploy</code>;</li>
                            <li><code>werf dismiss</code>;</li>
                            <li><code>werf cleanup</code>.</li>
                        </ul>
                    </li>
                    <li>
                        Open Source, written in Golang.
                    </li>
                    <li>
                        Werf is not a SAAS, we consider it a new generation<br/> of high-level CI/CD tools.
                    </li>
                </ul>
            </div>
        </div>
    </div>
</div>

<div class="page__container">
    <ul class="intro-extra">
        <li class="intro-extra__item">
            <div class="intro-extra__item-title">
                Comfortable deployment
            </div>
            <div class="intro-extra__item-text">
                <ul class="intro__list">
                    <li>Full compatibility with Helm.</li>
                    <li>Easy RBAC definition.</li>
                    <li>Applied deploy configuration in kubernetes does not guarantee a rolled out functioning application. With werf you receive that guarantee.</li>
                    <li>Werf fails fast when a problem has been detected in the CI/CD job allowing faster debug of new application versions without unnecessary kubectl invocations.</li>
                    <li>Configurable resource errors detector and resource readiness detector using resource annotations.</li>
                    <li>Rich logging and error reporting.</li>
                </ul>
            </div>
        </li>
        <li class="intro-extra__item">
            <div class="intro-extra__item-title">
                Image Lifecycle Management
            </div>
            <div class="intro-extra__item-text">
                <ul class="intro__list">
                    <li>Build images with Dockerfiles or advanced image builder with incremental rebuilds based on the git history and ansible.</li>
                    <li>Publish images to the registry using advanced image naming schemas.</li>
                    <li>Deploy application images into Kubernetes cluster.</li>
                    <li>Cleanup registry from unused images based on policies.</li>
                </ul>
            </div>
        </li>
    </ul>
</div>

<div class="stats">
    <div class="page__container">
        <div class="stats__content">
            <div class="stats__title">Active development & adoption</div>
            <ul class="stats__list">
                <li class="stats__list-item">
                    <div class="stats__list-item-num">4</div>
                    <div class="stats__list-item-title">releases per week</div>
                    <div class="stats__list-item-subtitle">on average for the last year</div>
                </li>
                <li class="stats__list-item">
                    <div class="stats__list-item-num">1200</div>
                    <div class="stats__list-item-title">installations</div>
                    <div class="stats__list-item-subtitle">of large and small projects</div>
                </li>
                <li class="stats__list-item">
                    <div class="stats__list-item-num gh_counter">1010</div>
                    <div class="stats__list-item-title">stars on GitHub</div>
                    <div class="stats__list-item-subtitle">let’s make it more ;)</div>
                </li>
            </ul>
        </div>
    </div>
</div>

<div class="features">
    <div class="page__container">
        <ul class="features__list">
            <li class="features__list-item">
                <div class="features__list-item-icon features__list-item-icon_easy"></div>
                <div class="features__list-item-title">Easy to start</div>
                <div class="features__list-item-text">Remain regular build process with Dockerfile. Take your project and put into werf easily right now.</div>
            </li>
            <li class="features__list-item">
                <div class="features__list-item-icon features__list-item-icon_config"></div>
                <div class="features__list-item-title">Compact configuration file</div>
                <div class="features__list-item-text">Build multiple images with a single configuration file, share common configuration parts using go-templates.</div>
            </li>
            <li class="features__list-item">
                <div class="features__list-item-icon features__list-item-icon_lifecycle"></div>
                <div class="features__list-item-title">Complete application lifecycle management</div>
                <div class="features__list-item-text">Manage image building process, deploy applications into Kubernetes and remove unused images easily.</div>
            </li>
            <li class="features__list-item">
                <div class="features__list-item-icon features__list-item-icon_size"></div>
                <div class="features__list-item-title">Reduce image size</div>
                <div class="features__list-item-text">Detach source data and build tools using artifacts, mounts and stapel.</div>
            </li>
            <li class="features__list-item">
                <div class="features__list-item-icon features__list-item-icon_ansible"></div>
                <div class="features__list-item-title">Build images with <span>Ansible</span></div>
                <div class="features__list-item-text">Use the powerful and popular infrastructure-as-a-code tool.</div>
            </li>
            <li class="features__list-item">
                <div class="features__list-item-icon features__list-item-icon_debug"></div>
                <div class="features__list-item-title">Advanced tools for debugging the build process</div>
                <div class="features__list-item-text">In the process of assembling, you can access a certain stage using introspection options.</div>
            </li>
            <li class="features__list-item"></li>
            <li class="features__list-item">
                <div class="features__list-item-icon features__list-item-icon_kubernetes"></div>
                <div class="features__list-item-title">Comfortable deployment to <span>Kubernetes</span></div>
                <div class="features__list-item-text">Deploy to Kubernetes using standard Kubernetes package manager with interactive tracking of the deployment process and real-time logs browsing.</div>
            </li>
            <li class="features__list-item"></li>
        </ul>
        <a href="https://github.com/flant/werf#complete-features-list" target="_blank" class="page__btn page__btn_o features__btn">
            Check out a complete features list
        </a>
    </div>
</div>

<div class="community">
    <div class="page__container">
        <div class="community__content">
            <div class="community__title">Friendly growing community</div>
            <div class="community__subtitle">Werf’s developers are always in contact with community<br/> at Slack and Telegram.</div>
            <div class="community__btns">
                <a href="https://t.me/werf_ru" target="_blank" class="page__btn page__btn_w community__btn">
                    <span class="page__icon page__icon_telegram"></span>
                    Join via Telegram
                </a>
                <a href="https://cloud-native.slack.com/messages/CHY2THYUU" target="_blank" class="page__btn page__btn_w community__btn">
                    <span class="page__icon page__icon_slack"></span>
                    Join via Slack
                </a>
            </div>
        </div>
    </div>
</div>

<div class="roadmap">
    <div class="page__container">
        <div class="roadmap__title">
            Roadmap
        </div>
        <div class="roadmap__content">
            <div class="roadmap__goals">
                <div class="roadmap__goals-content">
                    <div class="roadmap__goals-title">Goals</div>
                    <ul class="roadmap__goals-list">
                        <li class="roadmap__goals-list-item">
                            Feature complete version of werf which works well in an environment with a single dedicated persistent host to run all werf operations (build, deploy and cleanup).
                        </li>
                        <li class="roadmap__goals-list-item">
                            Proven approaches and recipes <br/>
                            for most popular CI systems.
                        </li>
                        <li class="roadmap__goals-list-item">
                            Building images completely in userspace, <br/>
                            a container or Kubernetes cluster.
                        </li>
                    </ul>
                </div>
            </div>
            <div class="roadmap__steps">
                <div class="roadmap__steps-content">
                    <div class="roadmap__steps-title">Steps</div>
                    <ul class="roadmap__steps-list">
                        <li class="roadmap__steps-list-item" data-roadmap-step="1616">
                            <a href="https://github.com/flant/werf/issues/1616" class="roadmap__steps-list-item-issue" target="_blank">#1616</a>
                            <span class="roadmap__steps-list-item-text">
                                Use <a href="https://kubernetes.io/docs/tasks/manage-kubernetes-objects/declarative-config/#merge-patch-calculation" target="_blank">3-way-merge</a> during helm release upgrade.
                            </span>
                        </li>
                        <li class="roadmap__steps-list-item" data-roadmap-step="1184">
                            <a href="https://github.com/flant/werf/issues/1184" class="roadmap__steps-list-item-issue" target="_blank">#1184</a>
                            <span class="roadmap__steps-list-item-text">
                                Content addressable tagging scheme.
                            </span>
                        </li>
                        <li class="roadmap__steps-list-item" data-roadmap-step="1617">
                            <a href="https://github.com/flant/werf/issues/1617" class="roadmap__steps-list-item-issue" target="_blank">#1617</a>
                            <span class="roadmap__steps-list-item-text">
                                Proven approaches and recipes<br/>
                                for most popular CI systems.
                            </span>
                        </li>
                        <li class="roadmap__steps-list-item" data-roadmap-step="1614">
                            <a href="https://github.com/flant/werf/issues/1614" class="roadmap__steps-list-item-issue" target="_blank">#1614</a>
                            <span class="roadmap__steps-list-item-text">
                                Distributed builds.
                            </span>
                        </li>
                        <li class="roadmap__steps-list-item" data-roadmap-step="1606">
                            <a href="https://github.com/flant/werf/issues/1606" class="roadmap__steps-list-item-issue" target="_blank">#1606</a>
                            <span class="roadmap__steps-list-item-text">
                                Helm 3 support.
                            </span>
                        </li>
                        <li class="roadmap__steps-list-item" data-roadmap-step="1618">
                            <a href="https://github.com/flant/werf/issues/1618" class="roadmap__steps-list-item-issue" target="_blank">#1618</a>
                            <span class="roadmap__steps-list-item-text">
                                Userspace builds without Docker daemon<br/>
                                (as in <a href="https://github.com/GoogleContainerTools/kaniko" target="_blank">kaniko</a>).
                            </span>
                        </li>
                    </ul>
                </div>
            </div>
        </div>
    </div>
</div>

<div class="page__container">
    <div class="documentation">
        <div class="documentation__image">
        </div>
        <div class="documentation__info">
            <div class="documentation__info-title">
                Complete documentation
            </div>
            <div class="documentation__info-text">
                Documentation of werf comprises ~100 articles which include common use cases (getting started, deploy to Kubernetes, CI/CD integration and more), comprehensive description of its functions & architecture, as well as CLI, commands.
            </div>
        </div>
        <div class="documentation__btns">
            <a href="https://github.com/flant/werf" target="_blank" class="page__btn page__btn_b documentation__btn">
                Get Werf
            </a>
            <a href="{{ site.baseurl }}/documentation/guides/getting_started.html" class="page__btn page__btn_o documentation__btn">
                Starters guide
            </a>
            <a href="{{ site.baseurl }}/documentation/cli/main/build.html" class="page__btn page__btn_o documentation__btn">
                Explore CLI
            </a>
        </div>
    </div>
</div>
