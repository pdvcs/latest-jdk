#!/usr/bin/env bash
set -e

export OP="$1"
if [[ "$OP" == "lts" ]]; then
    export JAVA_UPGRADE_VERSION_FILE=installed_java_lts_version
    export CMDLINE_LATEST_JAVA_VER="latest-jdk -lts -jv"
    export CMDLINE_JAVA_UPGRADE_URL="latest-jdk -lts"
    export JDK_DIR_LINK_NAME=lts-latest
elif [[ "$OP" == "latest" ]]; then
    export JAVA_UPGRADE_VERSION_FILE=installed_java_version
    export CMDLINE_LATEST_JAVA_VER="latest-jdk -jv"
    export CMDLINE_JAVA_UPGRADE_URL="latest-jdk"
    export JDK_DIR_LINK_NAME=latest
else
    echo "Usage: upgrade-java.sh <lts | latest>"
    exit 0
fi

echo "Determining $OP version..."
export LATEST_JAVA_VER=$($CMDLINE_LATEST_JAVA_VER)
export JAVA_UPGRADE_URL=$($CMDLINE_JAVA_UPGRADE_URL)

mkdir -p ~/.local/etc
touch ~/.local/etc/$JAVA_UPGRADE_VERSION_FILE

TMP_WORK_DIR=$(mktemp -d)
pushd $TMP_WORK_DIR > /dev/null 2>&1

INSTALLED_JAVA_VER=$(cat ~/.local/etc/$JAVA_UPGRADE_VERSION_FILE)
if [[ "${LATEST_JAVA_VER}" > "${INSTALLED_JAVA_VER}" ]]; then
    echo -n "Available: ${LATEST_JAVA_VER}, Installed: ${INSTALLED_JAVA_VER}, upgrade? (Y/n) "
    read proceed
    if [[ ${proceed,,} == "n" ]]; then
        exit 0
    fi
    rm OpenJDK*.tar.gz > /dev/null 2>&1 || true
    echo "Downloading Java $LATEST_JAVA_VER from Adoptium ..."
    curl -#fLO "${JAVA_UPGRADE_URL}"
    echo "Decompressing to /usr/local/java ..."
    sudo tar -C /usr/local/java -xzf OpenJDK*.tar.gz
    LATEST_JDK_FOLDER=$(ls /usr/local/java | egrep "jdk-${LATEST_JAVA_VER}.*" | sort -r | head -1)
    sudo ln -sf /usr/local/java/$LATEST_JDK_FOLDER /usr/local/java/${JDK_DIR_LINK_NAME}
    echo -n "$LATEST_JAVA_VER" > ~/.local/etc/$JAVA_UPGRADE_VERSION_FILE
    echo "Now installed: $(java -version 2>&1 | head -1)"
else
    echo "Java version ${INSTALLED_JAVA_VER} is already installed, no upgrade necessary."
fi
popd > /dev/null 2>&1
rm -rf $TMP_WORK_DIR
