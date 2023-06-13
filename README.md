# assignment_demo_2023

![Tests](https://github.com/TikTokTechImmersion/assignment_demo_2023/actions/workflows/test.yml/badge.svg)

This is a demo and template for backend assignment of 2023 TikTok Tech Immersion.

## Installation

Requirement:

- golang 1.18+
- docker

To install dependency tools:

```bash
make pre
```

## Run

```bash
docker-compose up -d
```

Check if it's running:

```bash
curl localhost:8080/ping
```

## To send Message

localhost:8080/api/send

```json
{
  "text": "Hello World",
  "chat":"Nishita-Mini",
  "Sender": "Nishita"
}
```

## To pull all Messages

GET localhost:8080/api/pull

```json
{
  "Chat": "Nishita-Mini",
  "Cursor": 0,
  "Limit": 10,
  "Reverse": false
}
```
