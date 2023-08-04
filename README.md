# go-text-to-wav
Super simple text to speech (TTS) converter.

> This demo application was created to demonstrate how to use the [VoiceDock TTS API](https://voicedock.app/specs/#tts-api).

# Installation
Create directories for model data and configuration:
```bash
mkdir dataset
mkdir config
```
Download an example configuration (more information [here](https://github.com/voicedock/ttspiper)):
```bash
curl -o config/ttspiper.json https://raw.githubusercontent.com/voicedock/ttspiper/main/config/ttspiper.json
```
Launching a remote "ttspiper" server that converts text to speech and implements the
[VoiceDock TTS API](https://github.com/voicedock/voicedock-specs/blob/main/proto/voicedock/core/tts/v1/tts_api.proto).
```bash
docker run --rm \
  -v "$(pwd)/config:/data/config" \
  -v "$(pwd)/dataset:/data/dataset" \
  -p 9999:9999 \
  ghcr.io/voicedock/ttspiper:latest ttspiper
```

Clone this repo and run build:
```bash
go build
```

# Usage
Show supported language packs on the remote server:
```bash
./go-text-to-wav -a 127.0.0.1:9999 -c list
```
```
Lang    Speaker Downloaded
ru      Irina   true
en-us   Ryan    false
```
Start downloading a language pack on a remote server:
```bash
./go-text-to-wav -a 127.0.0.1:9999 -c download -l en-us -s Ryan
```
```
Staring download (lang: en-us, speaker: Ryan)
Download complete
```
Run text-to-speech and save the result to a wav file.
```bash
echo "Hi, friend! I speak the text in English." | ./go-text-to-wav -a 127.0.0.1:9999 -l en-us -s Ryan > out.wav
```
or
```bash
cat ./text.txt | ./go-text-to-wav -a 127.0.0.1:9999 -l en-us -s Ryan > out.wav
```

## CONTRIBUTING
Lint proto files:
```bash
docker run --rm -w "/work" -v "$(pwd):/work" bufbuild/buf:latest lint internal/api/grpc/proto
```
Generate grpc interface:
```bash
docker run --rm -w "/work" -v "$(pwd):/work" ghcr.io/voicedock/protobuilder:1.0.0 generate internal/api/grpc/proto --template internal/api/grpc/proto/buf.gen.yaml
```