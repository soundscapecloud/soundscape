# Streamlist - open source self-hosted personal music server

## Features

* **Playlists**
  * Create unlimited playlists from your library
  * Shareable link for each playlist
* **HTTP Streaming to any device**
  * Transcodes to streamable MP4.
  * Converts any source to AAC-encoded MP4 (.m4a) audio files using ffmpeg
  * Plays in all modern web browsers, media players, and podcast apps
* **Private Podcast URLs**
  * Each Playlist has its own private podcast URL
  * Works with any podcast app on desktop or mobile
  * Enables downloading for offline access
* **Lightweight player**
  * Uses browser-native player with minimal additions
  * Keyboard shortcuts
  * `space` (play/pause) `p` (previous) `n` (next) `m` (mute) `-`/`+` (volume)
  * Remembers your volume setting
* **Import from YouTube**
  * Built in search with queueing downloader
  * Cancel downloading/transcode jobs
* **Simple self-hosting**
  * Public Docker image
  * Single static Go binary with assets bundled
  * Automatic TLS using Let's Encrypt
  * Redirects http to https

## Running


### 1. Install Docker

```bash
# Update apt
$ sudo apt-get update

# Remove old docker install.
$ sudo apt-get remove docker docker-engine docker.io

# Ensure we have basics for apt-get.
$ sudo apt-get install \
    apt-transport-https \
    ca-certificates \
    curl \
    software-properties-common

# Add Docker's public key.
$ curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -

# Add Docker's apt repo
$ sudo add-apt-repository \
    "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
    $(lsb_release -cs) \
    stable"

# Update apt
$ sudo apt-get update

# Install Docker
$ sudo apt-get install docker-ce

# Run the hello-world test image
$ sudo docker run hello-world

```

### 2. Run as a Docker container

The official image is `streamlist/streamlist`, which should run in any up-to-date Docker environment.

```bash

# Your download directory should be bind-mounted as `/data` inside the container using the `--volume` flag.
$ mkdir /home/<username>/Music

$ sudo docker create                            \
    --name streamlist --init --restart always   \
    --publish 80:80 --publish 443:443           \
    --volume /home/<username>/Music:/data   \
    streamlist/streamlist:latest --letsencrypt --http-host music.example.com

$ sudo docker start streamlist

$ sudo docker logs -f streamlist
1.503869865804371e+09    info    Streamlist URL: https://music.example.com/streamlist
1.503869865804527e+09    info    Login credentials:  streamlist  /  1134423142

```

### 3. Updating the container image

Pull the latest image, remove the container, and re-create the container as explained above.

```bash
# Pull the latest image
$ sudo docker pull streamlist/streamlist

# Stop the container
$ sudo docker stop streamlist

# Remove the container (data is stored on the mounted volume)
$ sudo docker rm streamlist

# Re-create and start the container
$ sudo docker create ... (see above)

```

