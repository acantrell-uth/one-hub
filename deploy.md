# Deploying One-Hub Service
=====================================

## Step 1: Translate Texts

Before deploying the service, we need to translate all texts to English. This can be done using the following command:

```bash
python i18n/translate.py --repository_path . --json_file_path i18n/en.json
```

This script will use the `en.json` file to update the translated text.

## Step 2: Build Docker Image

Next, we need to build a new Docker image for our service. This can be done using the following command:

```bash
docker build --no-cache -t one-hub .
```

This command will create a new image named `one-hub` from the current directory.

## Step 3: Start Service

Finally, we need to start the service using Docker Compose. This can be done using the following command:

```bash
docker compose up -d
```

This command will start the service in detached mode, meaning it will run in the background.

**Important:** If you need to change default passwords or the port used by the service, please modify the `docker-compose.yml` file accordingly. Changes made to this file will be reflected when running the `docker compose up -d` command.

By following these steps, you should now have your One-Hub service deployed and running!
