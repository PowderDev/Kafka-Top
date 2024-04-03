FROM node:20-bullseye

RUN apt-get update \
    && apt-get install -y --no-install-recommends --no-install-suggests build-essential

WORKDIR /app
COPY package.json ./package.json
COPY package-lock.json ./package-lock.json
RUN npm install
COPY . .

RUN npm run build

EXPOSE 4173

CMD [ "npm", "run", "preview"]
