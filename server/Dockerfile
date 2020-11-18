FROM golang:1.14 AS build

RUN apt-get update && apt-get install -y \
    python3 python3-pip \
    fonts-liberation libappindicator3-1 libasound2 libatk-bridge2.0-0 \
    libnspr4 libnss3 lsb-release xdg-utils libxss1 libdbus-glib-1-2 \
    curl unzip wget \
    xvfb


# install geckodriver and firefox

RUN GECKODRIVER_VERSION=`curl https://github.com/mozilla/geckodriver/releases/latest | grep -Po 'v[0-9]+.[0-9]+.[0-9]+'` && \
    wget https://github.com/mozilla/geckodriver/releases/download/$GECKODRIVER_VERSION/geckodriver-$GECKODRIVER_VERSION-linux64.tar.gz && \
    tar -zxf geckodriver-$GECKODRIVER_VERSION-linux64.tar.gz -C /usr/local/bin && \
    chmod +x /usr/local/bin/geckodriver && \
    rm geckodriver-$GECKODRIVER_VERSION-linux64.tar.gz

RUN FIREFOX_SETUP=firefox-setup.tar.bz2 && \
    apt-get purge firefox && \
    wget -O $FIREFOX_SETUP "https://download.mozilla.org/?product=firefox-latest&os=linux64" && \
    tar xjf $FIREFOX_SETUP -C /opt/ && \
    ln -s /opt/firefox/firefox /usr/bin/firefox && \
    rm $FIREFOX_SETUP

# Install OpenJDK-11
RUN apt-get update && \
    apt-get install -y \
        default-jre-headless \ 
        # firefox \
        software-properties-common \
    && apt-get clean;

WORKDIR /usr/src/app/server
RUN mkdir -p /usr/src/app/server
RUN mkdir -p /usr/local/bin/
COPY . /usr/src/app/server

RUN chmod +x start.sh
RUN go mod download
RUN GOOS=linux go build -o backend main.go

# COPY /usr/src/app/server/backend /usr/local/bin/backend

ENV HOST=0.0.0.0

CMD ["./start.sh"]