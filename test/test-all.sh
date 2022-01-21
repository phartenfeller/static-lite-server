trap "trap - SIGTERM && kill -- -$$" SIGINT SIGTERM EXIT

cd ./01-basic && ../../static-lite-server ./config.json &
npx wait-on http://localhost:3002 &&
  npx cypress run --spec ./cypress/integration/01.spec.js &&
  killall static-lite-server
